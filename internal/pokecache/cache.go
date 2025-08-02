package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt 	time.Time
	val 		[]byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
	ttl 	time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	c := Cache{
		entries: make(map[string]cacheEntry),
		ttl: 	 ttl,
	}

	// Start the reaper loop to clean up expired entries.
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	// Lock the cache for writing to prevent concurrent writes.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add the entry to the cache with the current time.
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// Lock the cache for reading to prevent a concurrent write.
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get the entry from the cache.
	entry, exists := c.entries[key]
	if !exists || time.Since(entry.createdAt) > c.ttl {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	for {
		time.Sleep(c.ttl)

		// Lock the cache for writing to remove expired entries.
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.ttl {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}