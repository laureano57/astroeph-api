package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"astroeph-api/pkg/chart"
	"fmt"
)

// NatalService handles natal chart calculations
type NatalService struct {
	ephemeris        *astro.Ephemeris
	planetCalculator *astro.PlanetCalculator
	houseCalculator  *astro.HouseCalculator
	aspectCalculator *astro.AspectCalculator
	chartDrawer      *astro.ChartDrawer
	logger           *logging.Logger
}

// NewNatalService creates a new natal chart service
func NewNatalService(logger *logging.Logger) *NatalService {
	ephemeris, err := astro.NewEphemeris(logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize ephemeris for natal service")
		return nil
	}

	planetCalc := astro.NewPlanetCalculator(ephemeris)
	houseCalc := astro.NewHouseCalculator(ephemeris)
	aspectCalc := astro.NewAspectCalculator()
	chartDrawer := astro.NewChartDrawer()

	return &NatalService{
		ephemeris:        ephemeris,
		planetCalculator: planetCalc,
		houseCalculator:  houseCalc,
		aspectCalculator: aspectCalc,
		chartDrawer:      chartDrawer,
		logger:           logger,
	}
}

// NatalChartRequest represents a request for natal chart calculation
type NatalChartRequest struct {
	Day         int    `json:"day" binding:"required,min=1,max=31"`
	Month       int    `json:"month" binding:"required,min=1,max=12"`
	Year        int    `json:"year" binding:"required"`
	LocalTime   string `json:"local_time" binding:"required"` // HH:MM:SS format
	City        string `json:"city" binding:"required"`
	HouseSystem string `json:"house_system,omitempty"` // defaults to "Placidus"
	DrawChart   bool   `json:"draw_chart,omitempty"`   // whether to generate SVG chart
	SVGWidth    int    `json:"svg_width,omitempty"`    // width of SVG chart (defaults to 600)
	SVGTheme    string `json:"svg_theme,omitempty"`    // theme for SVG chart ("light", "dark", "mono")
	AIResponse  bool   `json:"ai_response,omitempty"`  // whether to format response for LLM
}

// NatalChartResponse represents the response from natal chart calculation
type NatalChartResponse struct {
	*domain.Chart
	AIFormattedResponse *string `json:"ai_formatted_response,omitempty"`
}

// CalculateNatalChart calculates a complete natal chart
func (ns *NatalService) CalculateNatalChart(req *NatalChartRequest) (*NatalChartResponse, error) {
	ns.logger.CalculationLogger().
		Str("city", req.City).
		Int("year", req.Year).
		Int("month", req.Month).
		Int("day", req.Day).
		Str("house_system", req.HouseSystem).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting natal chart calculation")

	// Set defaults
	if req.HouseSystem == "" {
		req.HouseSystem = "Placidus"
	}
	if req.SVGWidth <= 0 && req.DrawChart {
		req.SVGWidth = 600
	}

	// Get location information
	geocodingService := astro.GetGeocodingService()
	if geocodingService == nil {
		return nil, fmt.Errorf("geocoding service not available")
	}

	location, err := geocodingService.GetCityInfo(req.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get location for %s: %w", req.City, err)
	}

	// Parse time information
	timeInfo, err := domain.ParseTime(req.Year, req.Month, req.Day, req.LocalTime, location.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	// Create birth info
	birthInfo := domain.BirthInfo{
		Date:     timeInfo.FormatDateForDisplay(),
		Time:     timeInfo.FormatTimeOnly(),
		Location: *location,
	}

	// Create new natal chart
	natalChart := domain.NewChart(domain.ChartTypeNatal, req.City, birthInfo)
	natalChart.HouseSystem = req.HouseSystem
	natalChart.Timezone = location.Timezone
	natalChart.UTCTime = timeInfo.UTCTime

	// Calculate houses first (needed for planet house assignments)
	houseSystem := domain.HouseSystem(req.HouseSystem)
	houses, err := ns.houseCalculator.CalculateHouses(timeInfo, location, houseSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate houses: %w", err)
	}

	// Add houses to chart
	for _, house := range houses {
		natalChart.AddHouse(house)
	}

	// Extract house cusps for planet calculations
	var houseCusps []float64
	for _, house := range houses {
		houseCusps = append(houseCusps, house.CuspValue)
	}

	// Calculate planets
	planets, err := ns.planetCalculator.CalculateAllPlanets(timeInfo, houseCusps)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate planets: %w", err)
	}

	// Add planets to chart
	for _, planet := range planets {
		natalChart.AddPlanet(planet)
	}

	// Calculate aspects
	aspects := ns.aspectCalculator.CalculateAspects(planets)
	for _, aspect := range aspects {
		natalChart.AddAspect(aspect)
	}

	// Set chart angles (Ascendant and Midheaven)
	if len(houseCusps) >= 10 {
		ascendant := houseCusps[0] // 1st house cusp
		midheaven := houseCusps[9] // 10th house cusp
		natalChart.SetAngles(ascendant, midheaven)
	}

	// Generate SVG chart if requested
	if req.DrawChart {
		theme := ns.parseTheme(req.SVGTheme)
		svg, err := ns.chartDrawer.GenerateNatalChart(natalChart, req.SVGWidth, theme)
		if err != nil {
			ns.logger.Error().
				Err(err).
				Msg("Failed to generate chart SVG")
			// Don't fail the entire request if SVG generation fails
		} else {
			natalChart.ChartDraw = svg
		}
	}

	ns.logger.Info().
		Str("endpoint", "natal-chart").
		Int("planets_calculated", len(natalChart.Planets)).
		Int("houses_calculated", len(natalChart.Houses)).
		Int("aspects_found", len(natalChart.Aspects)).
		Msg("✨ Natal chart calculation completed successfully")

	return &NatalChartResponse{Chart: natalChart}, nil
}

// parseTheme converts theme string to chart theme type
func (ns *NatalService) parseTheme(themeStr string) *chart.ThemeType {
	return ns.chartDrawer.GetThemeFromString(themeStr)
}

// GetNatalChartFormatted returns a formatted natal chart for LLM consumption
func (ns *NatalService) GetNatalChartFormatted(req *NatalChartRequest) (string, error) {
	response, err := ns.CalculateNatalChart(req)
	if err != nil {
		return "", err
	}

	return ns.formatChartForLLM(response.Chart), nil
}

// formatChartForLLM formats a natal chart for LLM consumption
func (ns *NatalService) formatChartForLLM(chart *domain.Chart) string {
	formatted := fmt.Sprintf("NATAL CHART ANALYSIS\n")
	formatted += fmt.Sprintf("Birth Date: %s at %s\n", chart.BirthInfo.Date, chart.BirthInfo.Time)
	formatted += fmt.Sprintf("Birth Location: %s\n", chart.BirthInfo.Location.GetDisplayName())
	formatted += fmt.Sprintf("Coordinates: %s\n", chart.BirthInfo.Location.FormatCoordinates())
	formatted += fmt.Sprintf("House System: %s\n\n", chart.HouseSystem)

	// Planetary Positions
	formatted += "PLANETARY POSITIONS:\n"
	for _, planet := range chart.Planets {
		formatted += fmt.Sprintf("• %s: %s %s (House %d)\n",
			planet.Name, planet.Degree, planet.Sign, planet.House)
	}

	// Chart Angles
	formatted += "\nCHART ANGLES:\n"
	formatted += fmt.Sprintf("• Ascendant: %s %s\n", chart.Angles.Ascendant.Degree, chart.Angles.Ascendant.Sign)
	formatted += fmt.Sprintf("• Midheaven: %s %s\n", chart.Angles.Midheaven.Degree, chart.Angles.Midheaven.Sign)

	// House Cusps
	formatted += "\nHOUSE CUSPS:\n"
	for _, house := range chart.Houses {
		formatted += fmt.Sprintf("• House %d: %s %s\n", house.Number, house.Cusp, house.Sign)
	}

	// Major Aspects
	if len(chart.Aspects) > 0 {
		formatted += "\nMAJOR ASPECTS:\n"
		for _, aspect := range chart.Aspects {
			if aspect.IsMajorAspect() {
				formatted += fmt.Sprintf("• %s %s %s - %.1f° orb\n",
					aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb)
			}
		}
	}

	// Summary
	formatted += fmt.Sprintf("\nASTROLOGICAL SUMMARY:\n")
	formatted += fmt.Sprintf("This natal chart shows %d planetary positions across %d houses, with %d major aspects. ",
		len(chart.Planets), len(chart.Houses), len(astro.FilterMajorAspects(chart.Aspects)))

	// Element emphasis
	elementCounts := ns.getElementCounts(chart.Planets)
	dominantElement, maxCount := ns.getDominantElement(elementCounts)
	if maxCount > 1 {
		formatted += fmt.Sprintf("There is a notable emphasis in %s with %d planetary placements. ", dominantElement, maxCount)
	}

	formatted += "This chart provides a comprehensive astrological foundation for interpretation."

	return formatted
}

// getElementCounts counts planets by element
func (ns *NatalService) getElementCounts(planets []domain.Planet) map[string]int {
	counts := make(map[string]int)
	for _, planet := range planets {
		counts[planet.Element]++
	}
	return counts
}

// getDominantElement finds the element with the most planets
func (ns *NatalService) getDominantElement(elementCounts map[string]int) (string, int) {
	maxCount := 0
	dominantElement := ""

	for element, count := range elementCounts {
		if count > maxCount {
			maxCount = count
			dominantElement = element
		}
	}

	return dominantElement, maxCount
}

// GetSupportedHouseSystems returns available house systems
func (ns *NatalService) GetSupportedHouseSystems() []string {
	systems := astro.GetAvailableHouseSystems()
	var systemNames []string
	for _, system := range systems {
		systemNames = append(systemNames, string(system))
	}
	return systemNames
}

// ValidateNatalChartRequest validates a natal chart request
func (ns *NatalService) ValidateNatalChartRequest(req *NatalChartRequest) error {
	if req.Day < 1 || req.Day > 31 {
		return fmt.Errorf("day must be between 1 and 31")
	}

	if req.Month < 1 || req.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}

	if req.Year < 1800 || req.Year > 2200 {
		return fmt.Errorf("year must be between 1800 and 2200")
	}

	if req.LocalTime == "" {
		return fmt.Errorf("local_time is required")
	}

	if req.City == "" {
		return fmt.Errorf("city is required")
	}

	if req.HouseSystem != "" && !astro.IsValidHouseSystem(req.HouseSystem) {
		return fmt.Errorf("invalid house system: %s", req.HouseSystem)
	}

	return nil
}
