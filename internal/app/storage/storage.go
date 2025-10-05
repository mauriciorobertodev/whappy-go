package storage

import (
	"context"
	"errors"
	"io"
)

var (
	ErrStorageNotConfigured = errors.New("storage not configured")
)

type Storage interface {
	Save(ctx context.Context, key string, r io.Reader) error
	Load(ctx context.Context, key string) (io.ReadCloser, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, data []byte) error
	Delete(ctx context.Context, key string) error
	URL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Healthy(ctx context.Context) error
}
