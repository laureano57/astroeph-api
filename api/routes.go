package api

import (
	"net/http"

	"astroeph-api/models"
	"astroeph-api/pkg/chart"
	"astroeph-api/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API routes with the given router and astrology service
func RegisterRoutes(router *gin.Engine, astroService *services.AstrologyService) {
	// API versioning group
	v1 := router.Group("/api/v1")
	{
		// Natal chart endpoint
		v1.POST("/natal-chart", func(c *gin.Context) {
			handleNatalChart(c, astroService)
		})

		// Natal chart SVG endpoint
		v1.POST("/natal-chart/svg", func(c *gin.Context) {
			handleNatalChartSVG(c, astroService)
		})

		// Transits endpoint (placeholder for now)
		v1.POST("/transits", func(c *gin.Context) {
			handleTransits(c, astroService)
		})

		// Synastry endpoint (placeholder for now)
		v1.POST("/synastry", func(c *gin.Context) {
			handleSynastry(c, astroService)
		})

		// Additional endpoints to be implemented in later phases
		v1.POST("/composite-chart", func(c *gin.Context) {
			handleCompositeChart(c, astroService)
		})

		v1.POST("/progressions", func(c *gin.Context) {
			handleProgressions(c, astroService)
		})

		v1.POST("/solar-return", func(c *gin.Context) {
			handleSolarReturn(c, astroService)
		})

		v1.POST("/lunar-return", func(c *gin.Context) {
			handleLunarReturn(c, astroService)
		})
	}
}

// handleNatalChart processes natal chart calculation requests
func handleNatalChart(c *gin.Context, astroService *services.AstrologyService) {
	logger := services.AppLogger

	var req models.NatalChartRequest

	// Bind JSON request to struct with validation
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Str("endpoint", "natal-chart").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set default house system if not provided
	if req.HouseSystem == "" {
		req.HouseSystem = "Placidus"
	}

	// Set default SVG width if SVG is requested
	if req.GenerateSVG && req.SVGWidth <= 0 {
		req.SVGWidth = 600
	}

	logger.CalculationLogger().
		Str("city", req.City).
		Int("year", req.Year).
		Int("month", req.Month).
		Int("day", req.Day).
		Str("house_system", req.HouseSystem).
		Bool("ai_response", req.AIResponse).
		Bool("generate_svg", req.GenerateSVG).
		Int("svg_width", req.SVGWidth).
		Str("svg_theme", req.SVGTheme).
		Msg("🔮 Starting natal chart calculation")

	// Parse SVG theme if provided
	var themeType *chart.ThemeType
	if req.SVGTheme != "" {
		switch req.SVGTheme {
		case "light":
			theme := chart.ThemeLight
			themeType = &theme
		case "dark":
			theme := chart.ThemeDark
			themeType = &theme
		case "mono":
			theme := chart.ThemeMono
			themeType = &theme
		default:
			// Use default theme
			themeType = nil
		}
	}

	// Call the astrology service to calculate the natal chart with optional SVG
	chartData, err := astroService.CalculateNatalChartWithSVG(&req, req.GenerateSVG, req.SVGWidth, themeType)
	if err != nil {
		logger.Error().
			Err(err).
			Str("endpoint", "natal-chart").
			Str("city", req.City).
			Msg("Failed to calculate natal chart")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate natal chart",
			"details": err.Error(),
		})
		return
	}

	logger.Info().
		Str("endpoint", "natal-chart").
		Int("planets_calculated", len(chartData.Planets)).
		Int("houses_calculated", len(chartData.Houses)).
		Int("aspects_found", len(chartData.Aspects)).
		Msg("✨ Natal chart calculation completed successfully")

	// Check if AI-formatted response is requested
	if req.AIResponse {
		logger.Debug().
			Str("endpoint", "natal-chart").
			Msg("🤖 Generating LLM-optimized response")

		llmText := models.FormatNatalChartForLLM(chartData)
		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	// Return structured JSON response
	c.JSON(http.StatusOK, chartData)
}

// handleTransits processes transit calculation requests (placeholder)
func handleTransits(c *gin.Context, astroService *services.AstrologyService) {
	var req models.TransitsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call the astrology service (currently returns not implemented error)
	transitsData, err := astroService.CalculateTransits(&req)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error":   "Transits calculation not yet implemented",
			"details": err.Error(),
		})
		return
	}

	if req.AIResponse {
		llmText := models.FormatForLLM(transitsData)
		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	c.JSON(http.StatusOK, transitsData)
}

// handleSynastry processes synastry calculation requests (placeholder)
func handleSynastry(c *gin.Context, astroService *services.AstrologyService) {
	var req models.SynastryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call the astrology service (currently returns not implemented error)
	synastryData, err := astroService.CalculateSynastry(&req)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error":   "Synastry calculation not yet implemented",
			"details": err.Error(),
		})
		return
	}

	if req.AIResponse {
		llmText := models.FormatForLLM(synastryData)
		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	c.JSON(http.StatusOK, synastryData)
}

// Placeholder handlers for additional endpoints (to be implemented in later phases)

func handleCompositeChart(c *gin.Context, astroService *services.AstrologyService) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Composite chart calculation not yet implemented",
		"message": "This endpoint will be available in a future version",
	})
}

func handleProgressions(c *gin.Context, astroService *services.AstrologyService) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Progressions calculation not yet implemented",
		"message": "This endpoint will be available in a future version",
	})
}

func handleSolarReturn(c *gin.Context, astroService *services.AstrologyService) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Solar return calculation not yet implemented",
		"message": "This endpoint will be available in a future version",
	})
}

func handleLunarReturn(c *gin.Context, astroService *services.AstrologyService) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Lunar return calculation not yet implemented",
		"message": "This endpoint will be available in a future version",
	})
}

// handleNatalChartSVG processes natal chart SVG generation requests
func handleNatalChartSVG(c *gin.Context, astroService *services.AstrologyService) {
	logger := services.AppLogger

	var req models.NatalChartRequest

	// Bind JSON request to struct with validation
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Str("endpoint", "natal-chart-svg").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set defaults
	if req.HouseSystem == "" {
		req.HouseSystem = "Placidus"
	}
	if req.SVGWidth <= 0 {
		req.SVGWidth = 600
	}
	if req.SVGTheme == "" {
		req.SVGTheme = "dark"
	}

	// Force SVG generation
	req.GenerateSVG = true

	logger.CalculationLogger().
		Str("city", req.City).
		Int("year", req.Year).
		Int("month", req.Month).
		Int("day", req.Day).
		Str("house_system", req.HouseSystem).
		Int("svg_width", req.SVGWidth).
		Str("svg_theme", req.SVGTheme).
		Msg("🎨 Starting natal chart SVG generation")

	// Parse SVG theme
	var themeType *chart.ThemeType
	switch req.SVGTheme {
	case "light":
		theme := chart.ThemeLight
		themeType = &theme
	case "dark":
		theme := chart.ThemeDark
		themeType = &theme
	case "mono":
		theme := chart.ThemeMono
		themeType = &theme
	}

	// Calculate natal chart with SVG
	chartData, err := astroService.CalculateNatalChartWithSVG(&req, true, req.SVGWidth, themeType)
	if err != nil {
		logger.Error().
			Err(err).
			Str("endpoint", "natal-chart-svg").
			Str("city", req.City).
			Msg("Failed to generate natal chart SVG")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate natal chart SVG",
			"details": err.Error(),
		})
		return
	}

	if chartData.SVG == "" {
		logger.Error().
			Str("endpoint", "natal-chart-svg").
			Msg("SVG generation returned empty result")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "SVG generation failed",
			"details": "Generated SVG is empty",
		})
		return
	}

	logger.Info().
		Str("endpoint", "natal-chart-svg").
		Int("svg_length", len(chartData.SVG)).
		Msg("✨ Natal chart SVG generated successfully")

	// Return SVG with appropriate content type
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, chartData.SVG)
}
