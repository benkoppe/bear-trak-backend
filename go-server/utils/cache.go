package utils

import (
	"sync"
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
}

func NewCache[T any](expiration time.Duration, fetchFunc func() (T, error)) *Cache[T] {
	return &Cache[T]{
		expirationTime: expiration,
		fetchFunc:      fetchFunc,
	}
}

func (c *Cache[T]) Get() (T, error) {
	c.mutex.RLock()
	if !c.lastFetch.IsZero() && (c.expirationTime < 0 || time.Since(c.lastFetch) < c.expirationTime) {
		data := c.data
		c.mutex.RUnlock()
		return data, nil
	}
	c.mutex.RUnlock()
	// cache refresh needed
	return c.refresh()
}

func (c *Cache[T]) refresh() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
