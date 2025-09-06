package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
)

// CompositeService handles composite chart calculations
type CompositeService struct {
	synastryService *SynastryService
	chartDrawer     *astro.ChartDrawer
	logger          *logging.Logger
}

// NewCompositeService creates a new composite service
func NewCompositeService(logger *logging.Logger) *CompositeService {
	synastryService := NewSynastryService(logger)
	chartDrawer := astro.NewChartDrawer()

	return &CompositeService{
		synastryService: synastryService,
		chartDrawer:     chartDrawer,
		logger:          logger,
	}
}

// CompositeChartRequest represents a request for composite chart calculation
type CompositeChartRequest struct {
	Person1    PersonData `json:"person1" binding:"required"`
	Person2    PersonData `json:"person2" binding:"required"`
	DrawChart  bool       `json:"draw_chart,omitempty"`
	SVGWidth   int        `json:"svg_width,omitempty"`
	SVGTheme   string     `json:"svg_theme,omitempty"`
	AIResponse bool       `json:"ai_response,omitempty"`
}

// CompositeChartResponse represents the response from composite chart calculation
type CompositeChartResponse struct {
	CompositeChart      *domain.Chart `json:"composite_chart"`
	Person1Chart        *domain.Chart `json:"person1_chart"`
	Person2Chart        *domain.Chart `json:"person2_chart"`
	ChartDraw           string        `json:"chart_draw,omitempty"`
	AIFormattedResponse *string       `json:"ai_formatted_response,omitempty"`
}

// CalculateCompositeChart calculates a composite chart between two people
func (cs *CompositeService) CalculateCompositeChart(req *CompositeChartRequest) (*CompositeChartResponse, error) {
	cs.logger.CalculationLogger().
		Str("person1_city", req.Person1.City).
		Str("person2_city", req.Person2.City).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting composite chart calculation")

	// Get both natal charts first
	synastryReq := &SynastryRequest{
		Person1:   req.Person1,
		Person2:   req.Person2,
		DrawChart: false,
	}

	synastryResponse, err := cs.synastryService.CalculateSynastry(synastryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate individual charts: %w", err)
	}

	person1Chart := synastryResponse.Person1Chart
	person2Chart := synastryResponse.Person2Chart

	// Calculate composite chart
	compositeChart := cs.calculateComposite(person1Chart, person2Chart)

	// Create response
	response := &CompositeChartResponse{
		CompositeChart: compositeChart,
		Person1Chart:   person1Chart,
		Person2Chart:   person2Chart,
	}

	// Generate SVG chart if requested
	if req.DrawChart {
		theme := cs.chartDrawer.GetThemeFromString(req.SVGTheme)
		width := req.SVGWidth
		if width <= 0 {
			width = 600
		}

		svg, err := cs.chartDrawer.GenerateCompositeChart(person1Chart, person2Chart, width, theme)
		if err != nil {
			cs.logger.Error().
				Err(err).
				Msg("Failed to generate composite chart SVG")
		} else {
			response.ChartDraw = svg
		}
	}

	cs.logger.Info().
		Int("composite_planets", len(compositeChart.Planets)).
		Int("composite_houses", len(compositeChart.Houses)).
		Msg("✨ Composite chart calculation completed successfully")

	return response, nil
}

// calculateComposite creates a composite chart from two natal charts
func (cs *CompositeService) calculateComposite(chart1, chart2 *domain.Chart) *domain.Chart {
	// Create composite birth info (midpoint of locations and time)
	compositeBirthInfo := cs.createCompositeBirthInfo(chart1.BirthInfo, chart2.BirthInfo)

	// Create new composite chart
	composite := domain.NewChart(
		domain.ChartTypeComposite,
		fmt.Sprintf("Composite: %s & %s", chart1.Name, chart2.Name),
		compositeBirthInfo,
	)

	composite.HouseSystem = chart1.HouseSystem

	// Calculate midpoint planets
	compositePlanets := cs.calculateMidpointPlanets(chart1.Planets, chart2.Planets)
	for _, planet := range compositePlanets {
		composite.AddPlanet(planet)
	}

	// Calculate midpoint houses
	compositeHouses := cs.calculateMidpointHouses(chart1.Houses, chart2.Houses)
	for _, house := range compositeHouses {
		composite.AddHouse(house)
	}

	// Calculate midpoint angles
	ascMidpoint := cs.calculateMidpoint(chart1.Angles.Ascendant.Value, chart2.Angles.Ascendant.Value)
	mcMidpoint := cs.calculateMidpoint(chart1.Angles.Midheaven.Value, chart2.Angles.Midheaven.Value)
	composite.SetAngles(ascMidpoint, mcMidpoint)

	// Calculate aspects for composite planets
	aspectCalc := astro.NewAspectCalculator()
	aspects := aspectCalc.CalculateAspects(compositePlanets)
	for _, aspect := range aspects {
		composite.AddAspect(aspect)
	}

	return composite
}

// createCompositeBirthInfo creates composite birth information
func (cs *CompositeService) createCompositeBirthInfo(birth1, birth2 domain.BirthInfo) domain.BirthInfo {
	// Calculate midpoint location
	midpointLat := (birth1.Location.Latitude + birth2.Location.Latitude) / 2
	midpointLon := (birth1.Location.Longitude + birth2.Location.Longitude) / 2

	// Use first location's city name with "Composite" prefix
	compositeName := fmt.Sprintf("Composite Location")

	compositeLocation := domain.NewLocation(
		compositeName,
		compositeName,
		birth1.Location.Country, // Use first person's country
		midpointLat,
		midpointLon,
		birth1.Location.Timezone, // Use first person's timezone
	)

	// Create composite birth info
	return domain.BirthInfo{
		Date:     "Composite Date",
		Time:     "Composite Time",
		Location: *compositeLocation,
	}
}

// calculateMidpointPlanets calculates midpoint positions for planets
func (cs *CompositeService) calculateMidpointPlanets(planets1, planets2 []domain.Planet) []domain.Planet {
	var compositePlanets []domain.Planet

	// Create map for quick lookup of planets by name
	planetMap2 := make(map[string]domain.Planet)
	for _, planet := range planets2 {
		planetMap2[planet.Name] = planet
	}

	// Calculate midpoints for matching planets
	for _, planet1 := range planets1 {
		if planet2, exists := planetMap2[planet1.Name]; exists {
			midpointLon := cs.calculateMidpoint(planet1.Longitude, planet2.Longitude)

			compositePlanet := domain.NewPlanet(
				planet1.Name,
				midpointLon,
				0, // Latitude not calculated for composite
				0, // Speed not applicable
				1, // House will be recalculated if needed
			)
			compositePlanets = append(compositePlanets, compositePlanet)
		}
	}

	return compositePlanets
}

// calculateMidpointHouses calculates midpoint house cusps
func (cs *CompositeService) calculateMidpointHouses(houses1, houses2 []domain.House) []domain.House {
	var compositeHouses []domain.House

	// Create map for quick lookup
	houseMap2 := make(map[int]domain.House)
	for _, house := range houses2 {
		houseMap2[house.Number] = house
	}

	// Calculate midpoint cusps
	for _, house1 := range houses1 {
		if house2, exists := houseMap2[house1.Number]; exists {
			midpointCusp := cs.calculateMidpoint(house1.CuspValue, house2.CuspValue)

			compositeHouse := domain.NewHouse(house1.Number, midpointCusp)
			compositeHouses = append(compositeHouses, compositeHouse)
		}
	}

	return compositeHouses
}

// calculateMidpoint calculates the midpoint between two longitudes
func (cs *CompositeService) calculateMidpoint(lon1, lon2 float64) float64 {
	// Normalize angles to 0-360 range
	lon1 = cs.normalizeAngle(lon1)
	lon2 = cs.normalizeAngle(lon2)

	// Calculate the shorter arc between the two points
	diff := lon2 - lon1
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}

	midpoint := lon1 + diff/2
	return cs.normalizeAngle(midpoint)
}

// normalizeAngle normalizes an angle to 0-360 range
func (cs *CompositeService) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// GetCompositeChartFormatted returns formatted composite chart for LLM consumption
func (cs *CompositeService) GetCompositeChartFormatted(req *CompositeChartRequest) (string, error) {
	response, err := cs.CalculateCompositeChart(req)
	if err != nil {
		return "", err
	}

	return cs.formatCompositeForLLM(response), nil
}

// formatCompositeForLLM formats composite chart results for LLM consumption
func (cs *CompositeService) formatCompositeForLLM(response *CompositeChartResponse) string {
	formatted := "COMPOSITE CHART ANALYSIS\n\n"

	composite := response.CompositeChart

	// Basic information
	formatted += fmt.Sprintf("Composite Chart: %s\n", composite.Name)
	formatted += fmt.Sprintf("House System: %s\n\n", composite.HouseSystem)

	// Composite planetary positions
	formatted += "COMPOSITE PLANETARY POSITIONS:\n"
	for _, planet := range composite.Planets {
		formatted += fmt.Sprintf("• %s: %s %s (House %d)\n",
			planet.Name, planet.Degree, planet.Sign, planet.House)
	}

	// Composite angles
	formatted += "\nCOMPOSITE ANGLES:\n"
	formatted += fmt.Sprintf("• Ascendant: %s %s\n", composite.Angles.Ascendant.Degree, composite.Angles.Ascendant.Sign)
	formatted += fmt.Sprintf("• Midheaven: %s %s\n", composite.Angles.Midheaven.Degree, composite.Angles.Midheaven.Sign)

	// Composite aspects
	if len(composite.Aspects) > 0 {
		formatted += "\nCOMPOSITE ASPECTS:\n"
		majorAspects := astro.FilterMajorAspects(composite.Aspects)
		for _, aspect := range majorAspects {
			formatted += fmt.Sprintf("• %s %s %s - %.1f° orb\n",
				aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb)
		}
	}

	// Summary
	formatted += fmt.Sprintf("\nCOMPOSITE SUMMARY:\n")
	formatted += fmt.Sprintf("This composite chart represents the combined energies and purposes of the relationship. ")
	formatted += fmt.Sprintf("It shows %d planetary midpoints across %d houses with %d major aspects, ",
		len(composite.Planets), len(composite.Houses), len(astro.FilterMajorAspects(composite.Aspects)))
	formatted += "revealing the relationship's inherent dynamics and potential."

	return formatted
}
