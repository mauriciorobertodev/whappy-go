package fake

import (
	"sync"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
)

type FakeCache struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	value      []byte
	expiration int64
}

func NewFakeCache() *FakeCache {
	return &FakeCache{store: make(map[string]item)}
}

func (c *FakeCache) Set(key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	c.store[key] = item{value: value, expiration: exp}
	return nil
}

func (c *FakeCache) Forever(key string, value []byte) error {
	return c.Set(key, value, 0)
}

func (c *FakeCache) Get(key string) ([]byte, error) {
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

func (c *FakeCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

func (c *FakeCache) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]item)
	return nil
}

func (c *FakeCache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.store[key]
	return ok
}

func (c *FakeCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]item)
}

func (c *FakeCache) CountItems() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

func (c *FakeCache) CountExpiredItems() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	now := time.Now().UnixNano()
	for _, it := range c.store {
		if it.expiration > 0 && now > it.expiration {
			count++
		}
	}
	return count
}

func (c *FakeCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]item)
	go func() {
		time.Sleep(100 * time.Millisecond)
		c.mu.Lock()
		defer c.mu.Unlock()
		c.store = make(map[string]item)
	}()
}
