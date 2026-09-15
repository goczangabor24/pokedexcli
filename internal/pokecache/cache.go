package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry //has state, needs a struct
	mu    sync.Mutex            //map is not thread-safe, struct needs protection
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	} //creates a new cache

	go c.reapLoop(interval) //starts cleaning up the cache concurrently after interval

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock() //locks the map, so no other concurrent process can access it
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]

	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
