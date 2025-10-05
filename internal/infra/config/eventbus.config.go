package config

type EventBusDriver string

const (
	EventBusDriverInMemory EventBusDriver = "memory"
	EventBusDriverRedis    EventBusDriver = "redis"
)

func (d EventBusDriver) IsValid() bool {
	switch d {
	case EventBusDriverInMemory, EventBusDriverRedis:
		return true
	}
	return false
}

type EventBusConfig struct {
	Driver EventBusDriver
	*RedisConfig
}

func LoadEventBusConfig() *EventBusConfig {
	cfg := &EventBusConfig{
		Driver:      EventBusDriver(GetEnvString("EVENTBUS_DRIVER", "memory")),
		RedisConfig: LoadRedisConfig(),
	}

	if !cfg.Driver.IsValid() {
		panic("Invalid EVENTBUS_DRIVER: " + string(cfg.Driver))
	}

	if cfg.Driver == EventBusDriverRedis && !cfg.RedisConfig.IsConfigured() {
		panic("Redis configuration is required for Redis event bus")
	}

	return cfg
}

func (c *EventBusConfig) GetRedisAddress() string {
	if c.Driver != EventBusDriverRedis {
		panic("Event bus driver is not Redis")
	}

	return c.RedisConfig.GetAddress()
}
