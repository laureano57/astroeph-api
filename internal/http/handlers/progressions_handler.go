package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProgressionsHandler handles progressions requests
type ProgressionsHandler struct {
	progressionsService *service.ProgressionsService
	logger              *logging.Logger
}

// NewProgressionsHandler creates a new progressions handler
func NewProgressionsHandler(progressionsService *service.ProgressionsService, logger *logging.Logger) *ProgressionsHandler {
	return &ProgressionsHandler{
		progressionsService: progressionsService,
		logger:              logger,
	}
}

// HandleProgressions handles POST /api/v1/progressions
func (ph *ProgressionsHandler) HandleProgressions(c *gin.Context) {
	var req service.ProgressionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ph.logger.Error().
			Err(err).
			Str("endpoint", "progressions").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.AIResponse {
		llmText, err := ph.progressionsService.GetProgressionsFormatted(&req)
		if err != nil {
			ph.logger.Error().
				Err(err).
				Str("endpoint", "progressions").
				Msg("Failed to generate LLM-formatted progressions")

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate formatted progressions",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	response, err := ph.progressionsService.CalculateProgressions(&req)
	if err != nil {
		ph.logger.Error().
			Err(err).
			Str("endpoint", "progressions").
			Msg("Failed to calculate progressions")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate progressions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProgressionsFormatted handles POST /api/v1/progressions/formatted
func (ph *ProgressionsHandler) GetProgressionsFormatted(c *gin.Context) {
	var req service.ProgressionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ph.logger.Error().
			Err(err).
			Str("endpoint", "progressions-formatted").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	llmText, err := ph.progressionsService.GetProgressionsFormatted(&req)
	if err != nil {
		ph.logger.Error().
			Err(err).
			Str("endpoint", "progressions-formatted").
			Msg("Failed to generate formatted progressions")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate formatted progressions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"formatted_response": llmText,
	})
}
