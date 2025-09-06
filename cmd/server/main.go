package main

import (
	"log"
	"net/http"

	"astroeph-api/internal/astro"
	"astroeph-api/internal/config"
	httpRouter "astroeph-api/internal/http"
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	logger := logging.NewLogger()
	logger.Info().
		Str("version", "v1.0.0").
		Str("service", "astroeph-api").
		Msg("🚀 Starting AstroEph API server")

	// Initialize geocoding service
	if err := astro.Initialize(); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to initialize geocoding service")
		log.Fatalf("Failed to initialize geocoding service: %v", err)
	}
	logger.Info().Msg("🌍 Geocoding service initialized successfully")

	// Initialize services
	natalService := service.NewNatalService(logger)
	synastryService := service.NewSynastryService(logger)
	compositeService := service.NewCompositeService(logger)
	solarReturnService := service.NewSolarReturnService(logger)
	lunarReturnService := service.NewLunarReturnService(logger)
	progressionsService := service.NewProgressionsService(logger)

	logger.Info().Msg("✅ All services initialized successfully")

	// Set up HTTP router
	ginRouter := gin.Default()

	// Add health check endpoint
	ginRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Astrological calculation service is running",
			"version": "v1.0.0",
		})
	})

	// Register API routes
	httpRouter.RegisterRoutes(
		ginRouter,
		natalService,
		synastryService,
		compositeService,
		solarReturnService,
		lunarReturnService,
		progressionsService,
		logger,
	)

	// Start server
	port := cfg.Server.Port
	logger.Info().
		Str("port", port).
		Str("health_endpoint", "http://localhost:"+port+"/health").
		Str("api_endpoint", "http://localhost:"+port+"/api/v1/natal-chart").
		Msg("🌟 Server starting")

	if err := ginRouter.Run(":" + port); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to run server")
		log.Fatalf("Failed to run server: %v", err)
	}
}
