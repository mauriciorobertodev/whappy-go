package cache

import (
	"encoding/json"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
)

func Set[T any](c cache.Cache, key string, value T, ttl time.Duration) error {
	l := app.GetCacheLogger()

	data, err := serialize(value)
	if err != nil {
		l.Error("Error serializing cache value", "err", err)
		return err
	}
	return c.Set(key, data, ttl)
}

func Forever[T any](c cache.Cache, key string, value T) error {
	return Set(c, key, value, 0)
}

func Get[T any](c cache.Cache, key string) (T, error) {
	l := app.GetCacheLogger()

	var zero T
	data, err := c.Get(key)
	if err != nil {
		l.Error("Error getting cache value", "err", err)
		return zero, err
	}
	return unserialize[T](data)
}

func serialize[T any](value T) ([]byte, error) {
	return json.Marshal(value)
}

func unserialize[T any](data []byte) (T, error) {
	var zero T
	err := json.Unmarshal(data, &zero)
	return zero, err
}
