package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
	"strings"
)

// ProgressionsService handles secondary progressions calculations
type ProgressionsService struct {
	natalService *NatalService
	logger       *logging.Logger
}

// NewProgressionsService creates a new progressions service
func NewProgressionsService(logger *logging.Logger) *ProgressionsService {
	natalService := NewNatalService(logger)

	return &ProgressionsService{
		natalService: natalService,
		logger:       logger,
	}
}

// ProgressionsRequest represents a request for progressions calculation
type ProgressionsRequest struct {
	// Birth data
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`

	// Progression date
	ProgressionDay   int `json:"progression_day" binding:"required,min=1,max=31"`
	ProgressionMonth int `json:"progression_month" binding:"required,min=1,max=12"`
	ProgressionYear  int `json:"progression_year" binding:"required"`

	// Chart options
	HouseSystem string `json:"house_system,omitempty"`
	DrawChart   bool   `json:"draw_chart,omitempty"`
	SVGWidth    int    `json:"svg_width,omitempty"`
	SVGTheme    string `json:"svg_theme,omitempty"`
	AIResponse  bool   `json:"ai_response,omitempty"`
}

// ProgressionsResponse represents the response from progressions calculation
type ProgressionsResponse struct {
	NatalChart          *domain.Chart `json:"natal_chart"`
	ProgressedChart     *domain.Chart `json:"progressed_chart"`
	ProgressionDate     string        `json:"progression_date"`
	YearsProgressed     float64       `json:"years_progressed"`
	DaysProgressed      float64       `json:"days_progressed"`
	ChartDraw           string        `json:"chart_draw,omitempty"`
	AIFormattedResponse *string       `json:"ai_formatted_response,omitempty"`
}

// CalculateProgressions calculates secondary progressions
func (ps *ProgressionsService) CalculateProgressions(req *ProgressionsRequest) (*ProgressionsResponse, error) {
	ps.logger.CalculationLogger().
		Int("birth_year", req.BirthYear).
		Int("progression_year", req.ProgressionYear).
		Str("birth_city", req.BirthCity).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting progressions calculation")

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

	natalResponse, err := ps.natalService.CalculateNatalChart(natalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate natal chart: %w", err)
	}

	// Calculate years and days progressed
	yearsProgressed, daysProgressed := ps.calculateProgressionTime(
		req.BirthYear, req.BirthMonth, req.BirthDay,
		req.ProgressionYear, req.ProgressionMonth, req.ProgressionDay,
	)

	// Calculate progressed chart
	// In secondary progressions, 1 day = 1 year, so we add days equal to years progressed
	progressedDate := ps.calculateProgressedDate(
		req.BirthYear, req.BirthMonth, req.BirthDay, daysProgressed,
	)

	progressedReq := &NatalChartRequest{
		Day:         progressedDate.Day,
		Month:       progressedDate.Month,
		Year:        progressedDate.Year,
		LocalTime:   req.BirthTime, // Keep same birth time
		City:        req.BirthCity, // Keep same birth location
		HouseSystem: req.HouseSystem,
		DrawChart:   req.DrawChart,
		SVGWidth:    req.SVGWidth,
		SVGTheme:    req.SVGTheme,
	}

	progressedResponse, err := ps.natalService.CalculateNatalChart(progressedReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate progressed chart: %w", err)
	}

	// Update chart type and name
	progressedResponse.Chart.Type = domain.ChartTypeProgressions
	progressedResponse.Chart.Name = fmt.Sprintf("Progressions for %d-%02d-%02d",
		req.ProgressionYear, req.ProgressionMonth, req.ProgressionDay)

	progressionDate := fmt.Sprintf("%d-%02d-%02d", req.ProgressionYear, req.ProgressionMonth, req.ProgressionDay)

	response := &ProgressionsResponse{
		NatalChart:      natalResponse.Chart,
		ProgressedChart: progressedResponse.Chart,
		ProgressionDate: progressionDate,
		YearsProgressed: yearsProgressed,
		DaysProgressed:  daysProgressed,
		ChartDraw:       progressedResponse.Chart.ChartDraw,
	}

	ps.logger.Info().
		Float64("years_progressed", yearsProgressed).
		Float64("days_progressed", daysProgressed).
		Msg("✨ Progressions calculation completed successfully")

	return response, nil
}

// calculateProgressionTime calculates years and days progressed
func (ps *ProgressionsService) calculateProgressionTime(
	birthYear, birthMonth, birthDay int,
	progressionYear, progressionMonth, progressionDay int,
) (float64, float64) {

	// Create birth and progression dates for comparison
	birthDate := fmt.Sprintf("%d-%02d-%02d", birthYear, birthMonth, birthDay)
	progressionDate := fmt.Sprintf("%d-%02d-%02d", progressionYear, progressionMonth, progressionDay)

	// Simple calculation - in practice you'd use proper date arithmetic
	yearsProgressed := float64(progressionYear - birthYear)

	// Adjust for months and days (simplified)
	monthDiff := progressionMonth - birthMonth
	dayDiff := progressionDay - birthDay

	yearsProgressed += float64(monthDiff) / 12.0
	yearsProgressed += float64(dayDiff) / 365.25

	// In secondary progressions, days progressed = years progressed
	daysProgressed := yearsProgressed

	ps.logger.Debug().
		Str("birth_date", birthDate).
		Str("progression_date", progressionDate).
		Float64("years_progressed", yearsProgressed).
		Msg("Calculated progression time")

	return yearsProgressed, daysProgressed
}

// ProgressedDate represents a calculated progressed date
type ProgressedDate struct {
	Year  int
	Month int
	Day   int
}

// calculateProgressedDate calculates the progressed date for ephemeris calculation
func (ps *ProgressionsService) calculateProgressedDate(
	birthYear, birthMonth, birthDay int, daysProgressed float64,
) ProgressedDate {

	// Add the progressed days to the birth date
	// This is a simplified calculation - in practice you'd use proper date arithmetic
	totalDays := int(daysProgressed)

	year := birthYear
	month := birthMonth
	day := birthDay + totalDays

	// Handle day overflow (simplified)
	for day > 31 {
		day -= 30 // Simplified - should use actual days in month
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	return ProgressedDate{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

// GetProgressionsFormatted returns formatted progressions for LLM consumption
func (ps *ProgressionsService) GetProgressionsFormatted(req *ProgressionsRequest) (string, error) {
	response, err := ps.CalculateProgressions(req)
	if err != nil {
		return "", err
	}

	return ps.formatProgressionsForLLM(response), nil
}

// formatProgressionsForLLM formats progressions results for LLM consumption
func (ps *ProgressionsService) formatProgressionsForLLM(response *ProgressionsResponse) string {
	formatted := "SECONDARY PROGRESSIONS ANALYSIS\n\n"

	natalChart := response.NatalChart
	progressedChart := response.ProgressedChart

	// Basic information
	formatted += fmt.Sprintf("Progression Date: %s\n", response.ProgressionDate)
	formatted += fmt.Sprintf("Years Progressed: %.2f\n", response.YearsProgressed)
	formatted += fmt.Sprintf("Days Progressed: %.2f\n\n", response.DaysProgressed)

	// Compare natal and progressed planets
	formatted += "NATAL vs PROGRESSED POSITIONS:\n"
	for i, natalPlanet := range natalChart.Planets {
		if i < len(progressedChart.Planets) {
			progressedPlanet := progressedChart.Planets[i]
			if natalPlanet.Name == progressedPlanet.Name {
				formatted += fmt.Sprintf("• %s:\n", natalPlanet.Name)
				formatted += fmt.Sprintf("  Natal: %s %s (House %d)\n",
					natalPlanet.Degree, natalPlanet.Sign, natalPlanet.House)
				formatted += fmt.Sprintf("  Progressed: %s %s (House %d)\n",
					progressedPlanet.Degree, progressedPlanet.Sign, progressedPlanet.House)

				// Check for sign changes
				if natalPlanet.Sign != progressedPlanet.Sign {
					formatted += fmt.Sprintf("  ** SIGN CHANGE: From %s to %s **\n",
						natalPlanet.Sign, progressedPlanet.Sign)
				}
				formatted += "\n"
			}
		}
	}

	// Progressed angles
	formatted += "PROGRESSED ANGLES:\n"
	formatted += fmt.Sprintf("• Progressed Ascendant: %s %s\n",
		progressedChart.Angles.Ascendant.Degree, progressedChart.Angles.Ascendant.Sign)
	formatted += fmt.Sprintf("• Progressed Midheaven: %s %s\n\n",
		progressedChart.Angles.Midheaven.Degree, progressedChart.Angles.Midheaven.Sign)

	// Major progressed aspects
	majorAspects := astro.FilterMajorAspects(progressedChart.Aspects)
	if len(majorAspects) > 0 {
		formatted += "MAJOR PROGRESSED ASPECTS:\n"
		for _, aspect := range majorAspects {
			formatted += fmt.Sprintf("• %s %s %s - %.1f° orb\n",
				aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb)
		}
		formatted += "\n"
	}

	// Progressions interpretation
	formatted += "PROGRESSIONS INTERPRETATION:\n"
	formatted += fmt.Sprintf("These secondary progressions show your inner development and evolving consciousness over %.1f years. ",
		response.YearsProgressed)
	formatted += "Secondary progressions reveal the unfolding of your natal potential and inner growth patterns. "

	// Check for significant progressed movements
	significantChanges := ps.findSignificantProgressedChanges(natalChart.Planets, progressedChart.Planets)
	if len(significantChanges) > 0 {
		formatted += fmt.Sprintf("Key developments include: %s.", strings.Join(significantChanges, ", "))
	}

	return formatted
}

// findSignificantProgressedChanges identifies significant changes in progressions
func (ps *ProgressionsService) findSignificantProgressedChanges(natalPlanets, progressedPlanets []domain.Planet) []string {
	var changes []string

	// Create maps for quick lookup
	natalMap := make(map[string]domain.Planet)
	for _, planet := range natalPlanets {
		natalMap[planet.Name] = planet
	}

	// Check for sign changes in important planets
	importantPlanets := []string{"Sun", "Moon", "Mercury", "Venus", "Mars"}

	for _, progressedPlanet := range progressedPlanets {
		if natalPlanet, exists := natalMap[progressedPlanet.Name]; exists {
			// Check if it's an important planet and if it changed signs
			for _, important := range importantPlanets {
				if progressedPlanet.Name == important && natalPlanet.Sign != progressedPlanet.Sign {
					changes = append(changes, fmt.Sprintf("%s progression from %s to %s",
						important, natalPlanet.Sign, progressedPlanet.Sign))
				}
			}
		}
	}

	return changes
}
