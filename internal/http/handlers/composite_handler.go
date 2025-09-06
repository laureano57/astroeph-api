package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CompositeHandler handles composite chart requests
type CompositeHandler struct {
	compositeService *service.CompositeService
	logger           *logging.Logger
}

// NewCompositeHandler creates a new composite handler
func NewCompositeHandler(compositeService *service.CompositeService, logger *logging.Logger) *CompositeHandler {
	return &CompositeHandler{
		compositeService: compositeService,
		logger:           logger,
	}
}

// HandleCompositeChart handles POST /api/v1/composite-chart
func (ch *CompositeHandler) HandleCompositeChart(c *gin.Context) {
	var req service.CompositeChartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ch.logger.Error().
			Err(err).
			Str("endpoint", "composite-chart").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Calculate composite chart
	response, err := ch.compositeService.CalculateCompositeChart(&req)
	if err != nil {
		ch.logger.Error().
			Err(err).
			Str("endpoint", "composite-chart").
			Msg("Failed to calculate composite chart")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate composite chart",
			"details": err.Error(),
		})
		return
	}

	// Add AI-formatted response if requested
	if req.AIResponse {
		ch.logger.Debug().
			Str("endpoint", "composite-chart").
			Msg("🤖 Generating LLM-optimized response")

		llmText, err := ch.compositeService.GetCompositeChartFormatted(&req)
		if err != nil {
			ch.logger.Error().
				Err(err).
				Str("endpoint", "composite-chart").
				Msg("Failed to generate LLM-formatted composite chart")
			// Continue without formatted response
		} else {
			response.AIFormattedResponse = &llmText
		}
	}

	c.JSON(http.StatusOK, response)
}
