package main

import (
	"database/sql"
	"log"

	"github.com/ayushvyasgit/devoptics/services/auth-service/internal/config"
	"github.com/ayushvyasgit/devoptics/services/auth-service/internal/handler"
	"github.com/ayushvyasgit/devoptics/services/auth-service/internal/repository"
	"github.com/ayushvyasgit/devoptics/services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry, cfg.JWTRefreshExpiry)
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	r := gin.Default()

	r.GET("/health", authHandler.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.GetMe)
		}
	}

	log.Printf("🚀 Auth Service starting on port %s", cfg.Port)
	log.Printf("📝 Environment: %s", cfg.Env)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//
