package config

import "fmt"

type RedisConfig struct {
	Host string
	Port int
}

func LoadRedisConfig() *RedisConfig {
	cfg := &RedisConfig{
		Host: GetEnvString("REDIS_HOST", ""),
		Port: GetEnvInt("REDIS_PORT", 0),
	}

	if !cfg.IsConfigured() {
		return nil
	}

	return cfg
}

func (c *RedisConfig) IsConfigured() bool {
	return c.Host != "" && c.Port != 0
}

func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
