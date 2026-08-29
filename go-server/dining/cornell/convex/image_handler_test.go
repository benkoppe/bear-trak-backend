package convex

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConvexImageHandlerCachesImage(t *testing.T) {
	t.Parallel()

	var upstreamCalls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalls.Add(1)
		if r.URL.Path != "/api/storage/image-id" {
			t.Errorf("unexpected upstream path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "image/webp")
		_, _ = io.WriteString(w, "image contents")
	}))
	t.Cleanup(upstream.Close)

	handler := NewImageHandler(upstream.URL)
	for range 2 {
		response := serveImageRequest(handler, "/convexImages/image-id")
		if response.Code != http.StatusOK || response.Body.String() != "image contents" {
			t.Errorf("response = (%d, %q), want (200, %q)", response.Code, response.Body.String(), "image contents")
		}
		if response.Header().Get("Cache-Control") != convexImageCacheControl {
			t.Errorf("Cache-Control = %q, want %q", response.Header().Get("Cache-Control"), convexImageCacheControl)
		}
		if response.Header().Get("ETag") == "" {
			t.Error("ETag is empty")
		}
	}
	if got := upstreamCalls.Load(); got != 1 {
		t.Errorf("upstream calls = %d, want 1", got)
	}
}

func TestConvexImageHandlerCoalescesConcurrentMisses(t *testing.T) {
	t.Parallel()

	var upstreamCalls atomic.Int32
	var startedOnce sync.Once
	fetchStarted := make(chan struct{})
	releaseFetch := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		upstreamCalls.Add(1)
		startedOnce.Do(func() { close(fetchStarted) })
		<-releaseFetch
		w.Header().Set("Content-Type", "image/webp")
		_, _ = io.WriteString(w, "image")
	}))
	t.Cleanup(upstream.Close)

	handler := NewImageHandler(upstream.URL)
	start := make(chan struct{})
	statuses := make(chan int, 10)
	var requests sync.WaitGroup
	for range 10 {
		requests.Add(1)
		go func() {
			defer requests.Done()
			<-start
			statuses <- serveImageRequest(handler, "/convexImages/image-id").Code
		}()
	}
	close(start)
	<-fetchStarted
	deadline := time.After(time.Second)
	for len(handler.responseSlots) < 10 {
		select {
		case <-deadline:
			t.Fatal("concurrent requests did not reach the handler")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	close(releaseFetch)
	requests.Wait()
	close(statuses)

	for status := range statuses {
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d", status, http.StatusOK)
		}
	}
	if got := upstreamCalls.Load(); got != 1 {
		t.Errorf("upstream calls = %d, want 1", got)
	}
}

func TestConvexImageHandlerCachesNotFoundAndRejectsInvalidFiles(t *testing.T) {
	t.Parallel()

	var missingCalls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/storage/missing":
			missingCalls.Add(1)
			http.NotFound(w, r)
		case "/api/storage/large":
			_, _ = io.WriteString(w, "12345")
		case "/api/storage/not-an-image":
			w.Header().Set("Content-Type", "text/plain")
			_, _ = io.WriteString(w, "text")
		}
	}))
	t.Cleanup(upstream.Close)

	handler := NewImageHandler(upstream.URL)
	handler.maxImageBytes = 4

	for range 2 {
		if status := serveImageRequest(handler, "/convexImages/missing").Code; status != http.StatusNotFound {
			t.Errorf("missing image status = %d, want %d", status, http.StatusNotFound)
		}
	}
	if got := missingCalls.Load(); got != 1 {
		t.Errorf("missing image upstream calls = %d, want 1", got)
	}
	if status := serveImageRequest(handler, "/convexImages/large").Code; status != http.StatusBadGateway {
		t.Errorf("large image status = %d, want %d", status, http.StatusBadGateway)
	}
	if status := serveImageRequest(handler, "/convexImages/not-an-image").Code; status != http.StatusBadGateway {
		t.Errorf("non-image status = %d, want %d", status, http.StatusBadGateway)
	}
}

func TestConvexImageMemoryCacheEvictsLeastRecentlyUsed(t *testing.T) {
	t.Parallel()

	cache := newConvexImageMemoryCache(4, 2)
	cache.put("first", &cachedConvexImage{data: []byte("12")})
	cache.put("second", &cachedConvexImage{data: []byte("34")})
	_, _ = cache.get("first")
	cache.put("third", &cachedConvexImage{data: []byte("56")})

	if _, ok := cache.get("second"); ok {
		t.Error("least recently used entry was not evicted")
	}
	if _, ok := cache.get("first"); !ok {
		t.Error("recently used entry was evicted")
	}
}

func TestConvexImageHandlerBoundsPendingFetches(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	release := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.Header().Set("Content-Type", "image/webp")
		_, _ = io.WriteString(w, "image")
	}))
	t.Cleanup(upstream.Close)

	handler := NewImageHandler(upstream.URL)
	handler.maxPending = 1
	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		serveImageRequest(handler, "/convexImages/first")
	}()
	<-started

	if status := serveImageRequest(handler, "/convexImages/second").Code; status != http.StatusServiceUnavailable {
		t.Errorf("overloaded status = %d, want %d", status, http.StatusServiceUnavailable)
	}
	close(release)
	<-firstDone
}

func TestConvexImageHandlerDoesNotFollowRedirects(t *testing.T) {
	t.Parallel()

	var redirectedRequests atomic.Int32
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		redirectedRequests.Add(1)
		w.Header().Set("Content-Type", "image/webp")
		_, _ = io.WriteString(w, "image")
	}))
	t.Cleanup(redirectTarget.Close)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	t.Cleanup(upstream.Close)

	if status := serveImageRequest(NewImageHandler(upstream.URL), "/convexImages/image-id").Code; status != http.StatusBadGateway {
		t.Errorf("redirect status = %d, want %d", status, http.StatusBadGateway)
	}
	if got := redirectedRequests.Load(); got != 0 {
		t.Errorf("redirect target requests = %d, want 0", got)
	}
}

func serveImageRequest(handler http.Handler, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response
}
