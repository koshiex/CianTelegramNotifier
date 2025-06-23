package config

import (
	"os"
	"strconv"
)

type Config struct {
	TelegramToken      string
	CianAPIURL         string
	DatabasePath       string
	LogLevel           string
	CheckInterval      string
	HealthCheckEnabled bool
	HealthCheckPort    string
}

func New() *Config {
	return &Config{
		TelegramToken:      getEnv("TELEGRAM_TOKEN", ""),
		CianAPIURL:         getEnv("CIAN_API_URL", "http://localhost:5000"),
		DatabasePath:       getEnv("DATABASE_PATH", "./bot.db"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		CheckInterval:      getEnv("CHECK_INTERVAL", "10m"),
		HealthCheckEnabled: getBoolEnv("HEALTH_CHECK_ENABLED", true),
		HealthCheckPort:    getEnv("HEALTH_CHECK_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
