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
	createdAt time.Time //keeps track of when the entry was created
	val       []byte    //stores the retrieved data in bytes
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	} //creates a new cache

	go c.reapLoop(interval) //cleans up the cache concurrently after elapsed interval

	return c
}

func (c *Cache) Add(key string, val []byte) { //add method to Cache struct
	c.mu.Lock()         //locks the map, so no other concurrent process can access it
	defer c.mu.Unlock() //unlocks right before returning

	c.cache[key] = cacheEntry{ //adds a cache entry
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) { //get method Cache struct
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key] //checks if key exists and reads the bytes stored in the cache map at 'key'

	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval) //sends a tick after each interval to the ticker.C channel

	for range ticker.C { //waits for tick to arrive at ticker.C channel, then execute
		c.mu.Lock()

		for key, entry := range c.cache { //loops through the cache map
			if time.Since(entry.createdAt) > interval { //checks whether this entry is older than the cache interval
				delete(c.cache, key) //deletes cache map entry
			}
		}
		c.mu.Unlock()
	}
}
