package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	DatabaseURL  string
	WorkDir      string
	AIEngineURL  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:         getEnv("CODE_ANALYZER_PORT", "8082"),
		Env:          getEnv("GO_ENV", "development"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		WorkDir:      getEnv("WORK_DIR", "C:\\devoptics\\temp\\code-analyzer"),
		AIEngineURL:  getEnv("AI_ENGINE_URL", "http://localhost:8083"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
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