package utils

import (
	"sync"
	"time"
)

// NewCache creates a new cache with the given timeout in seconds.
func NewCache(log Logger, timeout int) *Cache {
	return (&Cache{
		timeout: timeout,
		log:     log,
		data:    make(map[string]cacheItem),
		mutex:   &sync.RWMutex{},
	}).runCleaner()
}

type Cache struct {
	data    map[string]cacheItem
	timeout int
	mutex   *sync.RWMutex
	log     Logger
}

type cacheItem struct {
	Timeout time.Time
	Data    []byte
}

// Get returns the data from the cache if it exists and is not expired.
// returns nil if the data is not in the cache or is expired.
func (c *Cache) Get(key string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.log.Logf("[DEBUG] Getting cache item [key: %v]", key)

	item, ok := c.data[key]
	if !ok {
		c.log.Logf("[DEBUG] Cache item not found [key: %v]", key)
		return nil
	}

	return item.Data
}

// Set adds the data to the cache with the given key.
func (c *Cache) Set(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.log.Logf("[DEBUG] Setting cache item [key: %v]", key)

	c.data[key] = cacheItem{
		Timeout: time.Now().Add(time.Duration(c.timeout) * time.Second),
		Data:    data,
	}
}

func (c *Cache) runCleaner() *Cache {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			c.clean()
		}
	}()

	return c
}

func (c *Cache) clean() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.log.Logf("[DEBUG] Cleaning cache [size: %v]", len(c.data))

	for key, item := range c.data {
		if time.Now().After(item.Timeout) {
			delete(c.data, key)
		}
	}

	c.log.Logf("[DEBUG] Cache cleaned [size: %v]", len(c.data))
}
