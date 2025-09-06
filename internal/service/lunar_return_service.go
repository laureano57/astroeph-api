package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
)

// LunarReturnService handles lunar return calculations
type LunarReturnService struct {
	natalService *NatalService
	logger       *logging.Logger
}

// NewLunarReturnService creates a new lunar return service
func NewLunarReturnService(logger *logging.Logger) *LunarReturnService {
	natalService := NewNatalService(logger)

	return &LunarReturnService{
		natalService: natalService,
		logger:       logger,
	}
}

// LunarReturnRequest represents a request for lunar return calculation
type LunarReturnRequest struct {
	// Birth data
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`

	// Return month/year and location
	ReturnMonth int    `json:"return_month" binding:"required,min=1,max=12"`
	ReturnYear  int    `json:"return_year" binding:"required"`
	ReturnCity  string `json:"return_city,omitempty"` // If different from birth city

	// Chart options
	HouseSystem string `json:"house_system,omitempty"`
	DrawChart   bool   `json:"draw_chart,omitempty"`
	SVGWidth    int    `json:"svg_width,omitempty"`
	SVGTheme    string `json:"svg_theme,omitempty"`
	AIResponse  bool   `json:"ai_response,omitempty"`
}

// LunarReturnResponse represents the response from lunar return calculation
type LunarReturnResponse struct {
	NatalChart          *domain.Chart `json:"natal_chart"`
	LunarReturnChart    *domain.Chart `json:"lunar_return_chart"`
	ReturnMonth         int           `json:"return_month"`
	ReturnYear          int           `json:"return_year"`
	ReturnDate          string        `json:"return_date"`
	ChartDraw           string        `json:"chart_draw,omitempty"`
	AIFormattedResponse *string       `json:"ai_formatted_response,omitempty"`
}

// CalculateLunarReturn calculates a lunar return chart
func (lrs *LunarReturnService) CalculateLunarReturn(req *LunarReturnRequest) (*LunarReturnResponse, error) {
	lrs.logger.CalculationLogger().
		Int("birth_year", req.BirthYear).
		Int("return_month", req.ReturnMonth).
		Int("return_year", req.ReturnYear).
		Str("birth_city", req.BirthCity).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting lunar return calculation")

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

	natalResponse, err := lrs.natalService.CalculateNatalChart(natalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate natal chart: %w", err)
	}

	// Estimate lunar return date (approximately every 27-28 days)
	// This is a simplified calculation - in practice, you'd need precise ephemeris calculations
	returnDay := lrs.estimateLunarReturnDay(req.ReturnMonth, req.ReturnYear, req.BirthDay)
	returnDate := fmt.Sprintf("%d-%02d-%02d", req.ReturnYear, req.ReturnMonth, returnDay)

	// Determine return location
	returnCity := req.ReturnCity
	if returnCity == "" {
		returnCity = req.BirthCity
	}

	// Calculate lunar return chart
	returnReq := &NatalChartRequest{
		Day:         returnDay,
		Month:       req.ReturnMonth,
		Year:        req.ReturnYear,
		LocalTime:   req.BirthTime, // Approximate - would need precise calculation
		City:        returnCity,
		HouseSystem: req.HouseSystem,
		DrawChart:   req.DrawChart,
		SVGWidth:    req.SVGWidth,
		SVGTheme:    req.SVGTheme,
	}

	returnResponse, err := lrs.natalService.CalculateNatalChart(returnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate lunar return chart: %w", err)
	}

	// Update chart type and name
	returnResponse.Chart.Type = domain.ChartTypeLunarReturn
	returnResponse.Chart.Name = fmt.Sprintf("Lunar Return %d-%02d", req.ReturnYear, req.ReturnMonth)

	response := &LunarReturnResponse{
		NatalChart:       natalResponse.Chart,
		LunarReturnChart: returnResponse.Chart,
		ReturnMonth:      req.ReturnMonth,
		ReturnYear:       req.ReturnYear,
		ReturnDate:       returnDate,
		ChartDraw:        returnResponse.Chart.ChartDraw,
	}

	lrs.logger.Info().
		Int("return_month", req.ReturnMonth).
		Int("return_year", req.ReturnYear).
		Msg("✨ Lunar return calculation completed successfully")

	return response, nil
}

// estimateLunarReturnDay estimates the day of lunar return (simplified)
func (lrs *LunarReturnService) estimateLunarReturnDay(month, year, birthDay int) int {
	// This is a very simplified estimation
	// In a real implementation, you would calculate when the Moon returns to its natal position

	// Lunar cycle is approximately 29.5 days
	// This is just a placeholder estimation
	estimatedDay := birthDay

	// Ensure day is valid for the month
	daysInMonth := lrs.getDaysInMonth(month, year)
	if estimatedDay > daysInMonth {
		estimatedDay = daysInMonth
	}

	return estimatedDay
}

// getDaysInMonth returns the number of days in a given month/year
func (lrs *LunarReturnService) getDaysInMonth(month, year int) int {
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if month == 2 && lrs.isLeapYear(year) {
		return 29
	}

	return daysInMonth[month-1]
}

// isLeapYear checks if a year is a leap year
func (lrs *LunarReturnService) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// GetLunarReturnFormatted returns formatted lunar return for LLM consumption
func (lrs *LunarReturnService) GetLunarReturnFormatted(req *LunarReturnRequest) (string, error) {
	response, err := lrs.CalculateLunarReturn(req)
	if err != nil {
		return "", err
	}

	return lrs.formatLunarReturnForLLM(response), nil
}

// formatLunarReturnForLLM formats lunar return results for LLM consumption
func (lrs *LunarReturnService) formatLunarReturnForLLM(response *LunarReturnResponse) string {
	formatted := fmt.Sprintf("LUNAR RETURN ANALYSIS - %d-%02d\n\n", response.ReturnYear, response.ReturnMonth)

	returnChart := response.LunarReturnChart

	// Basic information
	formatted += fmt.Sprintf("Return Date: %s\n", response.ReturnDate)
	formatted += fmt.Sprintf("Location: %s\n", returnChart.BirthInfo.Location.GetDisplayName())
	formatted += fmt.Sprintf("House System: %s\n\n", returnChart.HouseSystem)

	// Find Moon position in lunar return
	var moonPlanet *domain.Planet
	for _, planet := range returnChart.Planets {
		if planet.Name == "Moon" {
			moonPlanet = &planet
			break
		}
	}

	// Lunar return Moon position
	if moonPlanet != nil {
		formatted += "LUNAR RETURN MOON POSITION:\n"
		formatted += fmt.Sprintf("• Moon: %s %s (House %d)\n\n",
			moonPlanet.Degree, moonPlanet.Sign, moonPlanet.House)
	}

	// Key planetary positions
	formatted += "KEY PLANETARY POSITIONS:\n"
	personalPlanets := []string{"Sun", "Mercury", "Venus", "Mars"}
	for _, planetName := range personalPlanets {
		for _, planet := range returnChart.Planets {
			if planet.Name == planetName {
				formatted += fmt.Sprintf("• %s: %s %s (House %d)\n",
					planet.Name, planet.Degree, planet.Sign, planet.House)
				break
			}
		}
	}

	// Lunar return angles
	formatted += "\nLUNAR RETURN ANGLES:\n"
	formatted += fmt.Sprintf("• Ascendant: %s %s\n", returnChart.Angles.Ascendant.Degree, returnChart.Angles.Ascendant.Sign)
	formatted += fmt.Sprintf("• Midheaven: %s %s\n", returnChart.Angles.Midheaven.Degree, returnChart.Angles.Midheaven.Sign)

	// Major aspects involving the Moon
	majorAspects := astro.FilterMajorAspects(returnChart.Aspects)
	moonAspects := lrs.filterMoonAspects(majorAspects)
	if len(moonAspects) > 0 {
		formatted += "\nMAJOR ASPECTS TO THE MOON:\n"
		for _, aspect := range moonAspects {
			formatted += fmt.Sprintf("• Moon %s %s - %.1f° orb\n",
				aspect.Type, lrs.getOtherPlanet(aspect, "Moon"), aspect.Orb)
		}
	}

	// Lunar return interpretation
	formatted += "\nLUNAR RETURN INTERPRETATION:\n"
	monthName := lrs.getMonthName(response.ReturnMonth)
	formatted += fmt.Sprintf("This lunar return for %s %d shows the emotional themes and monthly experiences ahead. ",
		monthName, response.ReturnYear)

	if moonPlanet != nil {
		formatted += fmt.Sprintf("The Moon in %s suggests a focus on %s during this lunar cycle.",
			moonPlanet.Sign, lrs.getMoonSignKeyword(moonPlanet.Sign))
	}

	return formatted
}

// filterMoonAspects returns aspects involving the Moon
func (lrs *LunarReturnService) filterMoonAspects(aspects []domain.Aspect) []domain.Aspect {
	var moonAspects []domain.Aspect
	for _, aspect := range aspects {
		if aspect.Planet1 == "Moon" || aspect.Planet2 == "Moon" {
			moonAspects = append(moonAspects, aspect)
		}
	}
	return moonAspects
}

// getOtherPlanet returns the other planet in an aspect (not Moon)
func (lrs *LunarReturnService) getOtherPlanet(aspect domain.Aspect, exclude string) string {
	if aspect.Planet1 == exclude {
		return aspect.Planet2
	}
	return aspect.Planet1
}

// getMonthName returns the name of a month
func (lrs *LunarReturnService) getMonthName(month int) string {
	monthNames := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	if month >= 1 && month <= 12 {
		return monthNames[month]
	}
	return "Unknown"
}

// getMoonSignKeyword returns a keyword for Moon in a zodiac sign
func (lrs *LunarReturnService) getMoonSignKeyword(sign string) string {
	keywords := map[string]string{
		"Aries":       "emotional independence and new initiatives",
		"Taurus":      "emotional stability and comfort-seeking",
		"Gemini":      "emotional communication and mental stimulation",
		"Cancer":      "deep emotional needs and nurturing",
		"Leo":         "emotional creativity and self-expression",
		"Virgo":       "emotional organization and practical service",
		"Libra":       "emotional harmony and relationship focus",
		"Scorpio":     "emotional intensity and transformation",
		"Sagittarius": "emotional adventure and philosophical growth",
		"Capricorn":   "emotional responsibility and achievement",
		"Aquarius":    "emotional independence and humanitarian ideals",
		"Pisces":      "emotional sensitivity and spiritual connection",
	}

	if keyword, exists := keywords[sign]; exists {
		return keyword
	}
	return "emotional growth and development"
}
