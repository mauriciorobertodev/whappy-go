package config

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
)

type AppConfig struct {
	ENVIRONMENT  string
	APP_URL      string
	APP_PORT     string
	LOG_LEVEL    logger.Level
	TOKEN_HASHER string

	ADMIN_TOKEN string

	CACHE_FILE_UPLOAD_TTL time.Duration
}

func (c *AppConfig) IsProduction() bool {
	return c.ENVIRONMENT == "production"
}

func (c *AppConfig) IsDevelopment() bool {
	return c.ENVIRONMENT == "development"
}

func (c *AppConfig) HasAdminToken() bool {
	return c.ADMIN_TOKEN != ""
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		ENVIRONMENT:           GetEnvString("ENVIRONMENT", "development"),
		APP_URL:               GetEnvString("APP_URL", "http://localhost"),
		APP_PORT:              GetEnvString("APP_PORT", "8080"),
		LOG_LEVEL:             GetEnvLogLevel("LOG_LEVEL", logger.LevelInfo),
		TOKEN_HASHER:          GetEnvString("TOKEN_HASHER", "simple"), // bcrypt, simple
		ADMIN_TOKEN:           GetEnvString("ADMIN_TOKEN", ""),
		CACHE_FILE_UPLOAD_TTL: GetEnvDuration("CACHE_FILE_UPLOAD_TTL", 5*time.Minute),
	}
}
