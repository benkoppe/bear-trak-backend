package convex

import (
	"bytes"
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	convexImageCacheMaxBytes         = 32 << 20
	convexImageCacheMaxEntries       = 128
	convexImageNotFoundMaxEntries    = 256
	convexImageMaxBytes              = 4 << 20
	convexImageMaxPendingFetches     = 64
	convexImageMaxConcurrentFetches  = 4
	convexImageMaxActiveResponses    = 16
	convexImageFetchTimeout          = 15 * time.Second
	convexImageNotFoundCacheDuration = 5 * time.Minute
	convexImageCacheControl          = "public, max-age=31536000, immutable"
)

type cachedConvexImage struct {
	data        []byte
	contentType string
	etag        string
}

type convexImageNotFoundCache struct {
	mutex      sync.Mutex
	entries    map[string]time.Time
	maxEntries int
}

func newConvexImageNotFoundCache(maxEntries int) *convexImageNotFoundCache {
	return &convexImageNotFoundCache{
		entries:    make(map[string]time.Time),
		maxEntries: maxEntries,
	}
}

func (c *convexImageNotFoundCache) contains(storageID string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expires, ok := c.entries[storageID]
	if !ok {
		return false
	}
	if time.Now().After(expires) {
		delete(c.entries, storageID)
		return false
	}
	return true
}

func (c *convexImageNotFoundCache) put(storageID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.entries[storageID]; !ok && len(c.entries) >= c.maxEntries {
		for existingID, expires := range c.entries {
			if time.Now().After(expires) {
				delete(c.entries, existingID)
			}
		}
	}
	if _, ok := c.entries[storageID]; !ok && len(c.entries) >= c.maxEntries {
		for existingID := range c.entries {
			delete(c.entries, existingID)
			break
		}
	}
	c.entries[storageID] = time.Now().Add(convexImageNotFoundCacheDuration)
}

type convexImageCacheEntry struct {
	storageID string
	image     *cachedConvexImage
}

type convexImageMemoryCache struct {
	mutex      sync.Mutex
	entries    map[string]*list.Element
	recency    *list.List
	maxBytes   int64
	maxEntries int
	usedBytes  int64
}

func newConvexImageMemoryCache(maxBytes int64, maxEntries int) *convexImageMemoryCache {
	return &convexImageMemoryCache{
		entries:    make(map[string]*list.Element),
		recency:    list.New(),
		maxBytes:   maxBytes,
		maxEntries: maxEntries,
	}
}

func (c *convexImageMemoryCache) get(storageID string) (*cachedConvexImage, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, ok := c.entries[storageID]
	if !ok {
		return nil, false
	}

	c.recency.MoveToFront(element)
	return element.Value.(*convexImageCacheEntry).image, true
}

func (c *convexImageMemoryCache) put(storageID string, image *cachedConvexImage) {
	imageSize := int64(len(image.data))
	if imageSize > c.maxBytes || c.maxEntries <= 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if existing, ok := c.entries[storageID]; ok {
		entry := existing.Value.(*convexImageCacheEntry)
		c.usedBytes -= int64(len(entry.image.data))
		entry.image = image
		c.usedBytes += imageSize
		c.recency.MoveToFront(existing)
	} else {
		entry := &convexImageCacheEntry{storageID: storageID, image: image}
		c.entries[storageID] = c.recency.PushFront(entry)
		c.usedBytes += imageSize
	}

	for c.usedBytes > c.maxBytes || len(c.entries) > c.maxEntries {
		oldest := c.recency.Back()
		c.remove(oldest)
	}
}

func (c *convexImageMemoryCache) remove(element *list.Element) {
	entry := element.Value.(*convexImageCacheEntry)
	delete(c.entries, entry.storageID)
	c.usedBytes -= int64(len(entry.image.data))
	c.recency.Remove(element)
}

// ImageHandler proxies and caches immutable files from Convex image storage.
type ImageHandler struct {
	baseURL       *url.URL
	client        *http.Client
	cache         *convexImageMemoryCache
	notFoundCache *convexImageNotFoundCache
	maxImageBytes int64
	maxPending    int
	pendingMutex  sync.Mutex
	pending       map[string]*pendingConvexImage
	fetchSlots    chan struct{}
	responseSlots chan struct{}
}

type pendingConvexImage struct {
	done  chan struct{}
	image *cachedConvexImage
	err   error
}

// NewImageHandler creates a bounded, in-memory caching proxy for Convex images.
func NewImageHandler(convexCloudURL string) *ImageHandler {
	baseURL, err := url.Parse(strings.TrimRight(strings.TrimSpace(convexCloudURL), "/"))
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		baseURL = nil
	}

	return &ImageHandler{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: convexImageFetchTimeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		cache:         newConvexImageMemoryCache(convexImageCacheMaxBytes, convexImageCacheMaxEntries),
		notFoundCache: newConvexImageNotFoundCache(convexImageNotFoundMaxEntries),
		maxImageBytes: convexImageMaxBytes,
		maxPending:    convexImageMaxPendingFetches,
		pending:       make(map[string]*pendingConvexImage),
		fetchSlots:    make(chan struct{}, convexImageMaxConcurrentFetches),
		responseSlots: make(chan struct{}, convexImageMaxActiveResponses),
	}
}

func (h *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	storageID, ok := strings.CutPrefix(r.URL.Path, "/convexImages/")
	if !ok || !validConvexStorageID(storageID) || h.baseURL == nil {
		http.NotFound(w, r)
		return
	}
	select {
	case h.responseSlots <- struct{}{}:
		defer func() { <-h.responseSlots }()
	default:
		w.Header().Set("Retry-After", "1")
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	image, err := h.load(r.Context(), storageID)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, errConvexImageNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, errConvexImagePendingLimit) {
			status = http.StatusServiceUnavailable
		}
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Cache-Control", convexImageCacheControl)
	w.Header().Set("Content-Type", image.contentType)
	w.Header().Set("ETag", image.etag)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(image.data))
}

func (h *ImageHandler) load(ctx context.Context, storageID string) (*cachedConvexImage, error) {
	if image, ok := h.cache.get(storageID); ok {
		return image, nil
	}
	if h.notFoundCache.contains(storageID) {
		return nil, errConvexImageNotFound
	}

	h.pendingMutex.Lock()
	if image, ok := h.cache.get(storageID); ok {
		h.pendingMutex.Unlock()
		return image, nil
	}
	if h.notFoundCache.contains(storageID) {
		h.pendingMutex.Unlock()
		return nil, errConvexImageNotFound
	}

	pending, ok := h.pending[storageID]
	if !ok {
		if len(h.pending) >= h.maxPending {
			h.pendingMutex.Unlock()
			return nil, errConvexImagePendingLimit
		}
		pending = &pendingConvexImage{done: make(chan struct{})}
		h.pending[storageID] = pending
		go h.fetchPending(storageID, pending)
	}
	h.pendingMutex.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-pending.done:
		return pending.image, pending.err
	}
}

func (h *ImageHandler) fetchPending(storageID string, pending *pendingConvexImage) {
	ctx, cancel := context.WithTimeout(context.Background(), convexImageFetchTimeout)
	defer cancel()

	var image *cachedConvexImage
	var err error
	select {
	case h.fetchSlots <- struct{}{}:
		image, err = h.fetch(ctx, storageID)
		<-h.fetchSlots
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err == nil {
		h.cache.put(storageID, image)
	} else if errors.Is(err, errConvexImageNotFound) {
		h.notFoundCache.put(storageID)
	}

	h.pendingMutex.Lock()
	pending.image = image
	pending.err = err
	close(pending.done)
	delete(h.pending, storageID)
	h.pendingMutex.Unlock()
}

func (h *ImageHandler) fetch(ctx context.Context, storageID string) (*cachedConvexImage, error) {
	upstreamURL := *h.baseURL
	upstreamURL.Path = strings.TrimRight(upstreamURL.Path, "/") + "/api/storage/" + storageID
	upstreamURL.RawPath = ""
	upstreamURL.RawQuery = ""
	upstreamURL.Fragment = ""

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, upstreamURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create Convex image request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Convex image: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: upstream returned HTTP %d", errConvexImageNotFound, resp.StatusCode)
		}
		return nil, fmt.Errorf("convex image request returned HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > h.maxImageBytes {
		return nil, fmt.Errorf("convex image exceeds %d-byte limit", h.maxImageBytes)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, h.maxImageBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read Convex image: %w", err)
	}
	if int64(len(data)) > h.maxImageBytes {
		return nil, fmt.Errorf("convex image exceeds %d-byte limit", h.maxImageBytes)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || !allowedConvexImageMediaType(mediaType) {
		return nil, fmt.Errorf("convex file has non-image content type %q", contentType)
	}
	hash := sha256.Sum256(data)

	return &cachedConvexImage{
		data:        data,
		contentType: contentType,
		etag:        `"` + hex.EncodeToString(hash[:]) + `"`,
	}, nil
}

func allowedConvexImageMediaType(mediaType string) bool {
	switch mediaType {
	case "image/avif", "image/gif", "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func validConvexStorageID(storageID string) bool {
	if storageID == "" || len(storageID) > 128 {
		return false
	}

	for _, char := range storageID {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' {
			continue
		}
		return false
	}
	return true
}

var (
	errConvexImageNotFound     = errors.New("convex image not found")
	errConvexImagePendingLimit = errors.New("convex image pending fetch limit reached")
)
