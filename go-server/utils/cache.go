// Package utils includes all utility methods.
package utils

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// NoExpiration requests a cache that never expires.
const NoExpiration time.Duration = -1

type Cache[T any] struct {
	data           T
	lastFetch      time.Time
	expirationTime time.Duration
	fetchFunc      func() (T, error)
	mutex          sync.RWMutex
	refreshing     atomic.Bool
}

func NewCache[T any](name string, expiration time.Duration, fetchFunc func() (T, error)) *Cache[T] {
	c := &Cache[T]{
		expirationTime: expiration,
		fetchFunc:      fetchFunc,
	}
	// immediately load data after initialization
	log.Printf("initializing cache: %s with expiration %s...", name, expiration)
	_, err := c.refresh()
	if err != nil {
		log.Fatalf("error initializing cache: %s: %v", name, err)
	}
	log.Printf("successfully initialized cache: %s", name)

	return c
}

func (c *Cache[T]) Get() (T, error) {
	c.mutex.RLock()
	isExpired := c.lastFetch.IsZero() || (c.expirationTime >= 0 && time.Since(c.lastFetch) > c.expirationTime)
	hasData := !c.lastFetch.IsZero()
	data := c.data
	c.mutex.RUnlock()

	isVeryExpired := isExpired && time.Since(c.lastFetch) > 5*c.expirationTime
	// if very expired, block the thread
	if isVeryExpired {
		return c.refresh()
	}

	// if expired but we have data, trigger a refresh in the background and return old data
	if isExpired && hasData {
		if !c.refreshing.Load() {
			go c.refreshInBackground()
		}
		return data, nil
	}

	// if we have no data or we're expired without data, we must load synchronously
	if !hasData || isExpired {
		return c.refresh()
	}

	return data, nil
}

func (c *Cache[T]) refresh() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// double-check expiration to avoid duplicate work
	if !c.lastFetch.IsZero() && (c.expirationTime < 0 || time.Since(c.lastFetch) < c.expirationTime) {
		return c.data, nil
	}

	newData, err := c.fetchFunc()
	if err != nil {
		var zero T
		return zero, err
	}

	c.data = newData
	c.lastFetch = time.Now()
	return newData, nil
}

func (c *Cache[T]) refreshInBackground() {
	// if already refreshing, don't start another
	if !c.refreshing.CompareAndSwap(false, true) {
		return
	}
	defer c.refreshing.Store(false)

	newData, err := c.fetchFunc()
	if err != nil {
		return // fail silently in background refresh
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = newData
	c.lastFetch = time.Now()
}

// forces a synchronous refresh
func (c *Cache[T]) ForceRefresh() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	newData, err := c.fetchFunc()
	if err != nil {
		var zero T
		return zero, err
	}

	c.data = newData
	c.lastFetch = time.Now()
	return newData, nil
}
