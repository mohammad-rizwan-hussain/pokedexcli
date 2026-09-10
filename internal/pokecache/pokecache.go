package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu    sync.Mutex
	items map[string]*cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	ticker := time.NewTicker(interval)

	cache := Cache{
		items: make(map[string]*cacheEntry),
	}
	go func() {
		for range ticker.C {
			cache.reapLoop(interval)
		}
	}()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	data, ok := c.items[key]
	if !ok {
		return nil, false
	}
	return data.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.items {
		if time.Since(entry.createdAt) >= interval {
			delete(c.items, key)
		}
	}
}
