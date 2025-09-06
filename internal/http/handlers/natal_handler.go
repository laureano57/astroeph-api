package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NatalHandler handles natal chart requests
type NatalHandler struct {
	natalService *service.NatalService
	logger       *logging.Logger
}

// NewNatalHandler creates a new natal chart handler
func NewNatalHandler(natalService *service.NatalService, logger *logging.Logger) *NatalHandler {
	return &NatalHandler{
		natalService: natalService,
		logger:       logger,
	}
}

// HandleNatalChart handles POST /api/v1/natal-chart
func (nh *NatalHandler) HandleNatalChart(c *gin.Context) {
	var req service.NatalChartRequest

	// Bind JSON request to struct with validation
	if err := c.ShouldBindJSON(&req); err != nil {
		nh.logger.Error().
			Err(err).
			Str("endpoint", "natal-chart").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := nh.natalService.ValidateNatalChartRequest(&req); err != nil {
		nh.logger.Error().
			Err(err).
			Str("endpoint", "natal-chart").
			Msg("Invalid request parameters")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request parameters",
			"details": err.Error(),
		})
		return
	}

	// Set default house system if not provided
	if req.HouseSystem == "" {
		req.HouseSystem = "Placidus"
	}

	// Set default SVG width if SVG is requested
	if req.DrawChart && req.SVGWidth <= 0 {
		req.SVGWidth = 600
	}

	// Calculate natal chart
	response, err := nh.natalService.CalculateNatalChart(&req)
	if err != nil {
		nh.logger.Error().
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

	// Add AI-formatted response if requested
	if req.AIResponse {
		nh.logger.Debug().
			Str("endpoint", "natal-chart").
			Msg("🤖 Generating LLM-optimized response")

		llmText, err := nh.natalService.GetNatalChartFormatted(&req)
		if err != nil {
			nh.logger.Error().
				Err(err).
				Str("endpoint", "natal-chart").
				Msg("Failed to generate LLM-formatted natal chart")
			// Continue without formatted response instead of failing
		} else {
			response.AIFormattedResponse = &llmText
		}
	}

	// Return structured JSON response
	c.JSON(http.StatusOK, response)
}

// GetSupportedHouseSystems handles GET /api/v1/house-systems
func (nh *NatalHandler) GetSupportedHouseSystems(c *gin.Context) {
	houseSystems := nh.natalService.GetSupportedHouseSystems()

	c.JSON(http.StatusOK, gin.H{
		"house_systems": houseSystems,
		"default":       "Placidus",
	})
}
