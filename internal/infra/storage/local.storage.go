package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

type LocalStorage struct {
	cfg *config.StorageConfig
}

func NewLocalStorage(cfg *config.StorageConfig) *LocalStorage {
	err := os.MkdirAll(cfg.Path, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("failed to create storage directory: %v", err))
	}
	return &LocalStorage{cfg: cfg}
}

func (s *LocalStorage) Save(ctx context.Context, key string, r io.Reader) error {
	path := filepath.Join(s.cfg.Path, key)

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Load(ctx context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(s.cfg.Path, key)
	return os.Open(path)
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	path := filepath.Join(s.cfg.Path, key)
	return os.Remove(path)
}

func (s *LocalStorage) URL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", s.cfg.URL, "storage", key), nil
}

func (s *LocalStorage) Get(ctx context.Context, key string) ([]byte, error) {
	path := filepath.Join(s.cfg.Path, key)
	return os.ReadFile(path)
}

func (s *LocalStorage) Put(ctx context.Context, key string, data []byte) error {
	path := filepath.Join(s.cfg.Path, key)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	path := filepath.Join(s.cfg.Path, key)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (s *LocalStorage) Healthy(ctx context.Context) error {
	testKey := ".healthcheck"
	testData := []byte("ok")

	if err := s.Put(ctx, testKey, testData); err != nil {
		return fmt.Errorf("health failed on put: %w", err)
	}

	_, err := s.Get(ctx, testKey)
	if err != nil {
		return fmt.Errorf("health failed on get: %w", err)
	}

	_ = s.Delete(ctx, testKey)

	return nil
}
