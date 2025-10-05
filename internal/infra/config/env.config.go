package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
)

func GetEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err != nil {
			panic(fmt.Sprintf("Invalid integer value for %s: %s", key, value))
		}

		return intValue
	}

	return defaultValue
}

func GetEnvURL(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		url, err := url.ParseRequestURI(value)
		if err != nil {
			panic(fmt.Sprintf("Invalid URL value for %s: %s", key, value))
		}

		return url.String()
	}

	return defaultValue
}

func GetEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "1":
			return true
		case "false", "0":
			return false
		default:
			panic(fmt.Sprintf("Invalid boolean value for %s: %s", key, value))
		}
	}

	return defaultValue
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}

	return defaultValue
}

func GetEnvLogLevel(key string, defaultValue logger.Level) logger.Level {
	if value := os.Getenv(key); value != "" {
		level, err := logger.ParseLogLevel(value)
		if err == nil {
			return level
		}
	}

	return defaultValue
}
