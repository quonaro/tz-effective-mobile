package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddress string
	LogLevel    string
}

func Load(path string) (*Config, error) {
	if path != "" {
		_ = godotenv.Load(path)
	}

	dbURL := getEnv("DATABASE_URL", "postgres://subs:subs@localhost:5432/subscriptions?sslmode=disable")
	host := getEnv("BACKEND_HOST", "0.0.0.0")
	port := getEnv("BACKEND_PORT", "8080")
	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		DatabaseURL: dbURL,
		HTTPAddress:  host + ":" + port,
		LogLevel:    logLevel,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	d, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return time.Duration(d) * time.Second
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.HTTPAddress == "" {
		return fmt.Errorf("HTTP_ADDRESS is required")
	}
	return nil
}
