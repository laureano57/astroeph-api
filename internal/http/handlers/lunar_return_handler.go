package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LunarReturnHandler handles lunar return requests
type LunarReturnHandler struct {
	lunarReturnService *service.LunarReturnService
	logger             *logging.Logger
}

// NewLunarReturnHandler creates a new lunar return handler
func NewLunarReturnHandler(lunarReturnService *service.LunarReturnService, logger *logging.Logger) *LunarReturnHandler {
	return &LunarReturnHandler{
		lunarReturnService: lunarReturnService,
		logger:             logger,
	}
}

// HandleLunarReturn handles POST /api/v1/lunar-return
func (lrh *LunarReturnHandler) HandleLunarReturn(c *gin.Context) {
	var req service.LunarReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		lrh.logger.Error().
			Err(err).
			Str("endpoint", "lunar-return").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.AIResponse {
		llmText, err := lrh.lunarReturnService.GetLunarReturnFormatted(&req)
		if err != nil {
			lrh.logger.Error().
				Err(err).
				Str("endpoint", "lunar-return").
				Msg("Failed to generate LLM-formatted lunar return")

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate formatted lunar return",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	response, err := lrh.lunarReturnService.CalculateLunarReturn(&req)
	if err != nil {
		lrh.logger.Error().
			Err(err).
			Str("endpoint", "lunar-return").
			Msg("Failed to calculate lunar return")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate lunar return",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLunarReturnFormatted handles POST /api/v1/lunar-return/formatted
func (lrh *LunarReturnHandler) GetLunarReturnFormatted(c *gin.Context) {
	var req service.LunarReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		lrh.logger.Error().
			Err(err).
			Str("endpoint", "lunar-return-formatted").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	llmText, err := lrh.lunarReturnService.GetLunarReturnFormatted(&req)
	if err != nil {
		lrh.logger.Error().
			Err(err).
			Str("endpoint", "lunar-return-formatted").
			Msg("Failed to generate formatted lunar return")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate formatted lunar return",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"formatted_response": llmText,
	})
}
