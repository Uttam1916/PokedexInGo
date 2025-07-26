package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cac      map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(intv time.Duration) *Cache {
	c := Cache{
		cac:      make(map[string]cacheEntry),
		interval: intv,
	}
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cac[key] = cacheEntry{
		val:       value,
		createdAt: time.Now(),
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cac[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range c.cac {
			if time.Since(v.createdAt) > c.interval {
				delete(c.cac, k)
			}
		}

	}
}
