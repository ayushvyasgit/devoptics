package main

import (
	"log"

	"github.com/ayushvyasgit/devoptics/services/api-gateway/internal/config"
	"github.com/ayushvyasgit/devoptics/services/api-gateway/internal/router"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("Configuration loaded",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
	)

	// Setup router
	r, err := router.SetupRouter(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to setup router", zap.Error(err))
	}

	logger.Info("✅ API Gateway starting",
		zap.String("port", cfg.Port),
		zap.String("auth_service", cfg.AuthServiceURL),
		zap.String("code_analyzer", cfg.CodeAnalyzerServiceURL),
		zap.String("ai_engine", cfg.AIEngineServiceURL),
	)

	// Start server
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}