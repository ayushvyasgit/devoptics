package router

import (
	"time"

	"github.com/ayushvyasgit/devoptics/services/api-gateway/internal/config"
	"github.com/ayushvyasgit/devoptics/services/api-gateway/internal/middleware"
	"github.com/ayushvyasgit/devoptics/services/api-gateway/internal/proxy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, logger *zap.Logger) (*gin.Engine, error) {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Logger middleware
	r.Use(middleware.Logger(logger))

	// Rate limiting
	rateLimiter, err := middleware.NewRateLimiter(cfg.RedisURL, cfg.RateLimitPerMinute)
	if err != nil {
		logger.Warn("Failed to create rate limiter, continuing without it", zap.Error(err))
	} else {
		r.Use(rateLimiter.Limit())
		logger.Info("Rate limiting enabled", zap.Int("max_requests", cfg.RateLimitPerMinute))
	}

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Service proxies
	authProxy := proxy.NewServiceProxy(cfg.AuthServiceURL)
	codeAnalyzerProxy := proxy.NewServiceProxy(cfg.CodeAnalyzerServiceURL)
	metricsProxy := proxy.NewServiceProxy(cfg.MetricsCollectorServiceURL)
	aiEngineProxy := proxy.NewServiceProxy(cfg.AIEngineServiceURL)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "api-gateway",
			"version": "1.0.0",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authProxy.Forward)
			auth.POST("/login", authProxy.Forward)
		}

		// Protected auth routes
		authProtected := v1.Group("/auth")
		authProtected.Use(authMiddleware.RequireAuth())
		{
			authProtected.GET("/me", authProxy.Forward)
		}

		// Code analyzer routes (protected)
		analyzer := v1.Group("/analyzer")
		analyzer.Use(authMiddleware.RequireAuth())
		{
			analyzer.POST("/scan", codeAnalyzerProxy.Forward)
			analyzer.GET("/reports/:id", codeAnalyzerProxy.Forward)
			analyzer.GET("/reports", codeAnalyzerProxy.Forward)
		}

		// Metrics routes (protected)
		metrics := v1.Group("/metrics")
		metrics.Use(authMiddleware.RequireAuth())
		{
			metrics.POST("/ingest", metricsProxy.Forward)
			metrics.GET("/query", metricsProxy.Forward)
		}

		// AI Engine routes (protected)
		ai := v1.Group("/ai")
		ai.Use(authMiddleware.RequireAuth())
		{
			ai.POST("/code/analyze", aiEngineProxy.Forward)
			ai.POST("/code/similarity", aiEngineProxy.Forward)
			ai.POST("/anomaly/detect", aiEngineProxy.Forward)
			ai.POST("/incident/summarize", aiEngineProxy.Forward)
		}
	}

	return r, nil
}