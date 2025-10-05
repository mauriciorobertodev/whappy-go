package cache

import (
	"sync"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	value      []byte
	expiration int64
}

func NewInMemoryCache(config *config.CacheConfig) *InMemoryCache {
	return &InMemoryCache{store: make(map[string]item)}
}

func (c *InMemoryCache) Set(key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	c.store[key] = item{value: value, expiration: exp}
	return nil
}

func (c *InMemoryCache) Forever(key string, value []byte) error {
	return c.Set(key, value, 0)
}

func (c *InMemoryCache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	it, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return nil, cache.ErrNotFound
	}

	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		c.Delete(key)
		return nil, cache.ErrNotFound
	}

	if it.expiration == 0 {
		return it.value, nil
	}

	return it.value, nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

func (c *InMemoryCache) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]item)
	return nil
}
