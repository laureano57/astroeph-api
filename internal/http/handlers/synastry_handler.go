package handlers

import (
	"astroeph-api/internal/logging"
	"astroeph-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SynastryHandler handles synastry requests
type SynastryHandler struct {
	synastryService *service.SynastryService
	logger          *logging.Logger
}

// NewSynastryHandler creates a new synastry handler
func NewSynastryHandler(synastryService *service.SynastryService, logger *logging.Logger) *SynastryHandler {
	return &SynastryHandler{
		synastryService: synastryService,
		logger:          logger,
	}
}

// HandleSynastry handles POST /api/v1/synastry
func (sh *SynastryHandler) HandleSynastry(c *gin.Context) {
	var req service.SynastryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		sh.logger.Error().
			Err(err).
			Str("endpoint", "synastry").
			Msg("Invalid request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Calculate synastry
	response, err := sh.synastryService.CalculateSynastry(&req)
	if err != nil {
		sh.logger.Error().
			Err(err).
			Str("endpoint", "synastry").
			Msg("Failed to calculate synastry")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate synastry",
			"details": err.Error(),
		})
		return
	}

	// Add AI-formatted response if requested
	if req.AIResponse {
		sh.logger.Debug().
			Str("endpoint", "synastry").
			Msg("🤖 Generating LLM-optimized response")

		llmText, err := sh.synastryService.GetSynastryFormatted(&req)
		if err != nil {
			sh.logger.Error().
				Err(err).
				Str("endpoint", "synastry").
				Msg("Failed to generate LLM-formatted synastry")
			// Continue without formatted response instead of failing
		} else {
			response.AIFormattedResponse = &llmText
		}
	}

	c.JSON(http.StatusOK, response)
}
