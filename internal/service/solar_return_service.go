package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
)

// SolarReturnService handles solar return calculations
type SolarReturnService struct {
	natalService *NatalService
	logger       *logging.Logger
}

// NewSolarReturnService creates a new solar return service
func NewSolarReturnService(logger *logging.Logger) *SolarReturnService {
	natalService := NewNatalService(logger)

	return &SolarReturnService{
		natalService: natalService,
		logger:       logger,
	}
}

// SolarReturnRequest represents a request for solar return calculation
type SolarReturnRequest struct {
	// Birth data
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`

	// Return year and location
	ReturnYear int    `json:"return_year" binding:"required"`
	ReturnCity string `json:"return_city,omitempty"` // If different from birth city

	// Chart options
	HouseSystem string `json:"house_system,omitempty"`
	DrawChart   bool   `json:"draw_chart,omitempty"`
	SVGWidth    int    `json:"svg_width,omitempty"`
	SVGTheme    string `json:"svg_theme,omitempty"`
	AIResponse  bool   `json:"ai_response,omitempty"`
}

// SolarReturnResponse represents the response from solar return calculation
type SolarReturnResponse struct {
	NatalChart          *domain.Chart `json:"natal_chart"`
	SolarReturnChart    *domain.Chart `json:"solar_return_chart"`
	ReturnYear          int           `json:"return_year"`
	ReturnDate          string        `json:"return_date"`
	ChartDraw           string        `json:"chart_draw,omitempty"`
	AIFormattedResponse *string       `json:"ai_formatted_response,omitempty"`
}

// CalculateSolarReturn calculates a solar return chart
func (srs *SolarReturnService) CalculateSolarReturn(req *SolarReturnRequest) (*SolarReturnResponse, error) {
	srs.logger.CalculationLogger().
		Int("birth_year", req.BirthYear).
		Int("return_year", req.ReturnYear).
		Str("birth_city", req.BirthCity).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting solar return calculation")

	// Calculate natal chart first
	natalReq := &NatalChartRequest{
		Day:         req.BirthDay,
		Month:       req.BirthMonth,
		Year:        req.BirthYear,
		LocalTime:   req.BirthTime,
		City:        req.BirthCity,
		HouseSystem: req.HouseSystem,
		DrawChart:   false,
	}

	natalResponse, err := srs.natalService.CalculateNatalChart(natalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate natal chart: %w", err)
	}

	// Calculate solar return date (when Sun returns to natal position)
	// This is a simplified calculation - in practice, you'd need precise ephemeris calculations
	returnDate := fmt.Sprintf("%d-%02d-%02d", req.ReturnYear, req.BirthMonth, req.BirthDay)

	// Determine return location
	returnCity := req.ReturnCity
	if returnCity == "" {
		returnCity = req.BirthCity // Use birth city if no return city specified
	}

	// Calculate solar return chart
	returnReq := &NatalChartRequest{
		Day:         req.BirthDay,
		Month:       req.BirthMonth,
		Year:        req.ReturnYear,
		LocalTime:   req.BirthTime, // Approximate - would need precise calculation
		City:        returnCity,
		HouseSystem: req.HouseSystem,
		DrawChart:   req.DrawChart,
		SVGWidth:    req.SVGWidth,
		SVGTheme:    req.SVGTheme,
	}

	returnResponse, err := srs.natalService.CalculateNatalChart(returnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate solar return chart: %w", err)
	}

	// Update chart type and name
	returnResponse.Chart.Type = domain.ChartTypeSolarReturn
	returnResponse.Chart.Name = fmt.Sprintf("Solar Return %d", req.ReturnYear)

	response := &SolarReturnResponse{
		NatalChart:       natalResponse.Chart,
		SolarReturnChart: returnResponse.Chart,
		ReturnYear:       req.ReturnYear,
		ReturnDate:       returnDate,
		ChartDraw:        returnResponse.Chart.ChartDraw,
	}

	srs.logger.Info().
		Int("return_year", req.ReturnYear).
		Msg("✨ Solar return calculation completed successfully")

	return response, nil
}

// GetSolarReturnFormatted returns formatted solar return for LLM consumption
func (srs *SolarReturnService) GetSolarReturnFormatted(req *SolarReturnRequest) (string, error) {
	response, err := srs.CalculateSolarReturn(req)
	if err != nil {
		return "", err
	}

	return srs.formatSolarReturnForLLM(response), nil
}

// formatSolarReturnForLLM formats solar return results for LLM consumption
func (srs *SolarReturnService) formatSolarReturnForLLM(response *SolarReturnResponse) string {
	formatted := fmt.Sprintf("SOLAR RETURN ANALYSIS - %d\n\n", response.ReturnYear)

	returnChart := response.SolarReturnChart

	// Basic information
	formatted += fmt.Sprintf("Return Date: %s\n", response.ReturnDate)
	formatted += fmt.Sprintf("Location: %s\n", returnChart.BirthInfo.Location.GetDisplayName())
	formatted += fmt.Sprintf("House System: %s\n\n", returnChart.HouseSystem)

	// Solar return planetary positions
	formatted += "SOLAR RETURN PLANETARY POSITIONS:\n"
	for _, planet := range returnChart.Planets {
		formatted += fmt.Sprintf("• %s: %s %s (House %d)\n",
			planet.Name, planet.Degree, planet.Sign, planet.House)
	}

	// Solar return angles
	formatted += "\nSOLAR RETURN ANGLES:\n"
	formatted += fmt.Sprintf("• Ascendant: %s %s\n", returnChart.Angles.Ascendant.Degree, returnChart.Angles.Ascendant.Sign)
	formatted += fmt.Sprintf("• Midheaven: %s %s\n", returnChart.Angles.Midheaven.Degree, returnChart.Angles.Midheaven.Sign)

	// Major aspects in solar return
	majorAspects := astro.FilterMajorAspects(returnChart.Aspects)
	if len(majorAspects) > 0 {
		formatted += "\nMAJOR ASPECTS IN SOLAR RETURN:\n"
		for _, aspect := range majorAspects {
			formatted += fmt.Sprintf("• %s %s %s - %.1f° orb\n",
				aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb)
		}
	}

	// Solar return interpretation
	formatted += "\nSOLAR RETURN INTERPRETATION:\n"
	formatted += fmt.Sprintf("This solar return chart for %d shows the energetic themes and potential experiences for the year ahead. ", response.ReturnYear)
	formatted += fmt.Sprintf("The chart contains %d planetary positions with %d major aspects, indicating the key areas of focus and development. ",
		len(returnChart.Planets), len(majorAspects))

	// Highlight solar return Ascendant
	srAsc := returnChart.Angles.Ascendant.Sign
	formatted += fmt.Sprintf("The Solar Return Ascendant in %s suggests themes of %s will be prominent this year.",
		srAsc, srs.getSignKeyword(srAsc))

	return formatted
}

// getSignKeyword returns a keyword for a zodiac sign
func (srs *SolarReturnService) getSignKeyword(sign string) string {
	keywords := map[string]string{
		"Aries":       "initiative and new beginnings",
		"Taurus":      "stability and material focus",
		"Gemini":      "communication and learning",
		"Cancer":      "emotional security and home",
		"Leo":         "creativity and self-expression",
		"Virgo":       "organization and service",
		"Libra":       "relationships and balance",
		"Scorpio":     "transformation and depth",
		"Sagittarius": "expansion and adventure",
		"Capricorn":   "achievement and responsibility",
		"Aquarius":    "innovation and independence",
		"Pisces":      "spirituality and compassion",
	}

	if keyword, exists := keywords[sign]; exists {
		return keyword
	}
	return "personal growth"
}
