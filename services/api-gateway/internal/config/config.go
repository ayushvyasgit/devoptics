package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                       string
	Env                        string
	RedisURL                   string
	AuthServiceURL             string
	CodeAnalyzerServiceURL     string
	MetricsCollectorServiceURL string
	AIEngineServiceURL         string
	JWTSecret                  string
	RateLimitPerMinute         int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                       getEnv("API_GATEWAY_PORT", "8080"),
		Env:                        getEnv("GO_ENV", "development"),
		RedisURL:                   getEnv("REDIS_URL", "redis://localhost:6379"),
		AuthServiceURL:             getEnv("AUTH_SERVICE_URL", "http://localhost:8091"),
		CodeAnalyzerServiceURL:     getEnv("CODE_ANALYZER_URL", "http://localhost:8082"),
		MetricsCollectorServiceURL: getEnv("METRICS_COLLECTOR_URL", "http://localhost:8004"),
		AIEngineServiceURL:         getEnv("AI_ENGINE_URL", "http://localhost:8083"),
		JWTSecret:                  getEnv("JWT_SECRET", ""),
		RateLimitPerMinute:         60,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}