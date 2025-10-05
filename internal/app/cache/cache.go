package cache

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("key not found")
)

const (
	DefaultTTL = 5 * time.Minute

	CacheKeyFileUploadPrefix = "file:upload:"
	CacheKeyThumbnailPrefix  = "file:thumb:"
	CacheKeyGroupTypePrefix  = "group:type:"
	CacheKeyTokenPrefix      = "token:"
)

type Cache interface {
	Set(key string, value []byte, ttl time.Duration) error
	Forever(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Flush() error
}
