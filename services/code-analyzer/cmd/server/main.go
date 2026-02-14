package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/config"
	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create work directory
	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
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

	// Initialize handler
	h := handler.NewHandler(db, cfg.WorkDir)

	// Setup router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", h.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		analyzer := v1.Group("/analyzer")
		{
			analyzer.POST("/scan", h.ScanRepository)
			analyzer.GET("/reports/:id", h.GetReport)
			analyzer.GET("/reports/:id/issues", h.GetIssues)
			analyzer.GET("/reports", h.GetReports)
		}
	}

	log.Printf("🚀 Code Analyzer starting on port %s", cfg.Port)
	log.Printf("📁 Work directory: %s", cfg.WorkDir)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}