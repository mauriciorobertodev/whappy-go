package eventbus

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

func New(cfg *config.EventBusConfig) events.EventBus {
	if cfg == nil {
		panic("EventBusConfig is nil")
	}

	switch cfg.Driver {
	case config.EventBusDriverInMemory:
		return NewInMemoryEventBus()
	case config.EventBusDriverRedis:
		return NewRedisEventBus(cfg)
	default:
		panic("Unsupported event bus driver: " + string(cfg.Driver))
	}
}
