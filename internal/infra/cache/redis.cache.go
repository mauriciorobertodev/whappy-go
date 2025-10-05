package cache

import (
	"context"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(config *config.CacheConfig) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{Addr: config.GetRedisAddress()}),
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Forever(key string, value []byte) error {
	return r.Set(key, value, 0)
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	result := r.client.Get(r.ctx, key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, cache.ErrNotFound
		}
		return nil, result.Err()
	}

	return result.Bytes()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) Flush() error {
	return r.client.FlushAll(r.ctx).Err()
}
