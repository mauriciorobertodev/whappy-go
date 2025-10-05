package cache

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

func New(cfg *config.CacheConfig) cache.Cache {
	if cfg != nil && !cfg.IsConfigured() {
		return nil
	}

	switch cfg.Driver {
	case "memory":
		return NewInMemoryCache(cfg)
	case "redis":
		return NewRedisCache(cfg)
	default:
		return nil
	}
}
