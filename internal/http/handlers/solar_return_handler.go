package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SolarReturnHandler handles solar return requests
type SolarReturnHandler struct {
	solarReturnService *service.SolarReturnService
	logger             *logging.Logger
}

// NewSolarReturnHandler creates a new solar return handler
func NewSolarReturnHandler(solarReturnService *service.SolarReturnService, logger *logging.Logger) *SolarReturnHandler {
	return &SolarReturnHandler{
		solarReturnService: solarReturnService,
		logger:             logger,
	}
}

// HandleSolarReturn handles POST /api/v1/solar-return
func (srh *SolarReturnHandler) HandleSolarReturn(c *gin.Context) {
	var req service.SolarReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		srh.logger.Error().
			Err(err).
			Str("endpoint", "solar-return").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.AIResponse {
		llmText, err := srh.solarReturnService.GetSolarReturnFormatted(&req)
		if err != nil {
			srh.logger.Error().
				Err(err).
				Str("endpoint", "solar-return").
				Msg("Failed to generate LLM-formatted solar return")

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate formatted solar return",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ai_response": llmText,
		})
		return
	}

	response, err := srh.solarReturnService.CalculateSolarReturn(&req)
	if err != nil {
		srh.logger.Error().
			Err(err).
			Str("endpoint", "solar-return").
			Msg("Failed to calculate solar return")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate solar return",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSolarReturnFormatted handles POST /api/v1/solar-return/formatted
func (srh *SolarReturnHandler) GetSolarReturnFormatted(c *gin.Context) {
	var req service.SolarReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		srh.logger.Error().
			Err(err).
			Str("endpoint", "solar-return-formatted").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	llmText, err := srh.solarReturnService.GetSolarReturnFormatted(&req)
	if err != nil {
		srh.logger.Error().
			Err(err).
			Str("endpoint", "solar-return-formatted").
			Msg("Failed to generate formatted solar return")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate formatted solar return",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"formatted_response": llmText,
	})
}
