package cache

import (
	"sync"
	"time"
)

// CacheItem represents a single item in the cache.
type CacheItem[V any] struct {
	Value      V
	Expiration int64
}

// IsExpired checks if the cache item is expired.
func (item CacheItem[V]) IsExpired() bool {
	return time.Now().UnixNano() > item.Expiration
}

// Cache is the generic cache with expiration and mutex protection.
type Cache[K comparable, V any] struct {
	items             map[K]CacheItem[V]
	mutex             sync.RWMutex
	DefaultExpiration time.Duration
	cleanupInterval   time.Duration
	stopCleanup       chan bool
}

// NewCache creates a new cache instance.
func NewCache[K comparable, V any](defaultExpiration, cleanupInterval time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		items:             make(map[K]CacheItem[V]),
		DefaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}

	// Start the cleanup goroutine to remove expired items.
	go cache.startCleanup()

	return cache
}

// Set adds an item to the cache with the specified key and expiration duration.
func (c *Cache[K, V]) Set(key K, value V, duration time.Duration) {
	if duration == 0 {
		duration = c.DefaultExpiration
	}
	expiration := time.Now().Add(duration).UnixNano()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem[V]{
		Value:      value,
		Expiration: expiration,
	}
}

// Set adds an item to the cache with the specified key and expiration duration.
func (c *Cache[K, V]) SetWithDefaultExpiration(key K, value V) {
	expiration := time.Now().Add(c.DefaultExpiration).UnixNano()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem[V]{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache by key.
// It returns the value and a boolean indicating whether the key was found and not expired.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.items[key]
	if !found || item.IsExpired() {
		var zeroV V
		return zeroV, false
	}

	return item.Value, true
}

// Delete removes an item from the cache by key.
func (c *Cache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

// startCleanup starts a background goroutine that periodically removes expired items.
func (c *Cache[K, V]) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// deleteExpired removes all expired items from the cache.
func (c *Cache[K, V]) deleteExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
		}
	}
}

// Close stops the background cleanup goroutine.
func (c *Cache[K, V]) Close() {
	c.stopCleanup <- true
}
