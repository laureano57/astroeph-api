package http

import (
	"astroeph-api/internal/http/handlers"
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(
	router *gin.Engine,
	natalService *service.NatalService,
	synastryService *service.SynastryService,
	compositeService *service.CompositeService,
	solarReturnService *service.SolarReturnService,
	lunarReturnService *service.LunarReturnService,
	progressionsService *service.ProgressionsService,
	logger *logging.Logger,
) {
	// Add logging middleware
	router.Use(loggingMiddleware(logger))

	// Add CORS middleware
	router.Use(corsMiddleware())

	// API versioning group
	v1 := router.Group("/api/v1")
	{
		// Create handlers
		natalHandler := handlers.NewNatalHandler(natalService, logger)
		synastryHandler := handlers.NewSynastryHandler(synastryService, logger)
		compositeHandler := handlers.NewCompositeHandler(compositeService, logger)
		solarReturnHandler := handlers.NewSolarReturnHandler(solarReturnService, logger)
		lunarReturnHandler := handlers.NewLunarReturnHandler(lunarReturnService, logger)
		progressionsHandler := handlers.NewProgressionsHandler(progressionsService, logger)

		// Natal chart endpoints
		v1.POST("/natal-chart", natalHandler.HandleNatalChart)

		// Synastry endpoints
		v1.POST("/synastry", synastryHandler.HandleSynastry)

		// Composite chart endpoints
		v1.POST("/composite-chart", compositeHandler.HandleCompositeChart)

		// Solar return endpoints
		v1.POST("/solar-return", solarReturnHandler.HandleSolarReturn)

		// Lunar return endpoints
		v1.POST("/lunar-return", lunarReturnHandler.HandleLunarReturn)

		// Progressions endpoints
		v1.POST("/progressions", progressionsHandler.HandleProgressions)

		// Utility endpoints
		v1.GET("/house-systems", natalHandler.GetSupportedHouseSystems)
	}
}

// loggingMiddleware adds request logging
func loggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
