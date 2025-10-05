package storage

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

func New(cfg *config.StorageConfig) storage.Storage {
	if cfg != nil && cfg.IsConfigured() {
		if cfg.Driver == config.StorageDriverS3 {
			return NewS3Storage(cfg)
		}

		return NewLocalStorage(cfg)
	}

	return nil
}
