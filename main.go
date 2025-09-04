package main

import (
	"log"
	"net/http"
	"time"

	"astroeph-api/api"
	"astroeph-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize structured logging
	services.InitializeLogger()
	logger := services.AppLogger
	
	logger.Info().
		Str("version", "v1.0.0").
		Str("service", "astroeph-api").
		Msg("🚀 Starting AstroEph API server")

	// Initialize geocoding service
	if err := services.InitializeGeocodingService(); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to initialize geocoding service")
		log.Fatalf("Failed to initialize geocoding service: %v", err)
	}
	logger.Info().Msg("🌍 Geocoding service initialized successfully")

	// Initialize astrology service
	astroService, err := services.NewAstrologyService()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to initialize astrology service")
		log.Fatalf("Failed to initialize astrology service: %v", err)
	}

	logger.Info().Msg("✅ Astrology service initialized successfully")

	// Set up the router
	router := gin.Default()

	// Add logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		
		// Process request
		c.Next()
		
		// Log request details
		logger.RequestLogger().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Str("ip", c.ClientIP()).
			Msg("HTTP Request")
	})

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Astrological calculation service is running",
			"version": "v1.0.0",
		})
	})

	// Register API routes
	api.RegisterRoutes(router, astroService)

	// Run the server
	logger.Info().
		Str("port", "8080").
		Str("health_endpoint", "http://localhost:8080/health").
		Str("api_endpoint", "http://localhost:8080/api/v1/natal-chart").
		Msg("🌟 Server starting on port 8080")
		
	if err := router.Run(":8080"); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to run server")
		log.Fatalf("Failed to run server: %v", err)
	}
}
