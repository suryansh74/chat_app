package cache

import (
	"sync"
	"time"
)

type CachePort interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
}

type cacheItem struct {
	value   interface{}
	expires time.Time
}

type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]cacheItem
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, nil
	}

	if time.Now().After(item.expires) {
		return nil, nil
	}

	return item.value, nil
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expires := time.Now().Add(ttl)
	if ttl <= 0 {
		expires = time.Now().Add(24 * time.Hour)
	}

	c.data[key] = cacheItem{
		value:   value,
		expires: expires,
	}

	return nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

func (c *InMemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cacheItem)
	return nil
}
