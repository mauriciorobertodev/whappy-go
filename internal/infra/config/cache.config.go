package config

type CacheDriver string

const (
	CacheDriverInMemory CacheDriver = "memory"
	CacheDriverRedis    CacheDriver = "redis"
	CacheDriverNone     CacheDriver = "none"
)

type CacheConfig struct {
	Driver CacheDriver
	*RedisConfig
}

func (d CacheDriver) IsValid() bool {
	switch d {
	case CacheDriverInMemory, CacheDriverRedis, CacheDriverNone:
		return true
	default:
		return false
	}
}

func LoadCacheConfig() *CacheConfig {
	cfg := &CacheConfig{
		Driver:      CacheDriver(GetEnvString("CACHE_DRIVER", "memory")),
		RedisConfig: LoadRedisConfig(),
	}

	if !cfg.IsConfigured() {
		return nil
	}

	return cfg
}

func (c *CacheConfig) IsConfigured() bool {
	if !c.Driver.IsValid() {
		panic("Invalid Cache Driver: " + string(c.Driver))
	}

	if c.Driver == CacheDriverNone {
		return false
	}

	if c.Driver == CacheDriverRedis {
		return c.RedisConfig.IsConfigured()
	}

	return true
}

func (c *CacheConfig) GetRedisAddress() string {
	if c.Driver != CacheDriverRedis {
		panic("Cache driver is not Redis")
	}

	return c.RedisConfig.GetAddress()
}

func (c *CacheConfig) IsInMemory() bool {
	return c != nil && c.Driver == CacheDriverInMemory
}

func (c *CacheConfig) IsRedis() bool {
	return c != nil && c.Driver == CacheDriverRedis
}
