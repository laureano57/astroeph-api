package service

import (
	"astroeph-api/internal/astro"
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
)

// SynastryService handles synastry calculations
type SynastryService struct {
	natalService     *NatalService
	aspectCalculator *astro.AspectCalculator
	chartDrawer      *astro.ChartDrawer
	logger           *logging.Logger
}

// NewSynastryService creates a new synastry service
func NewSynastryService(logger *logging.Logger) *SynastryService {
	natalService := NewNatalService(logger)
	aspectCalc := astro.NewAspectCalculator()
	chartDrawer := astro.NewChartDrawer()

	return &SynastryService{
		natalService:     natalService,
		aspectCalculator: aspectCalc,
		chartDrawer:      chartDrawer,
		logger:           logger,
	}
}

// SynastryRequest represents a request for synastry calculation
type SynastryRequest struct {
	Person1    PersonData `json:"person1" binding:"required"`
	Person2    PersonData `json:"person2" binding:"required"`
	DrawChart  bool       `json:"draw_chart,omitempty"`
	SVGWidth   int        `json:"svg_width,omitempty"`
	SVGTheme   string     `json:"svg_theme,omitempty"`
	AIResponse bool       `json:"ai_response,omitempty"`
}

// PersonData represents birth data for one person
type PersonData struct {
	Day         int    `json:"day" binding:"required,min=1,max=31"`
	Month       int    `json:"month" binding:"required,min=1,max=12"`
	Year        int    `json:"year" binding:"required"`
	LocalTime   string `json:"local_time" binding:"required"`
	City        string `json:"city" binding:"required"`
	Name        string `json:"name,omitempty"`
	HouseSystem string `json:"house_system,omitempty"`
}

// SynastryResponse represents the response from synastry calculation
type SynastryResponse struct {
	Person1Chart        *domain.Chart   `json:"person1_chart"`
	Person2Chart        *domain.Chart   `json:"person2_chart"`
	SynastryAspects     []domain.Aspect `json:"synastry_aspects"`
	ChartDraw           string          `json:"chart_draw,omitempty"`
	AIFormattedResponse *string         `json:"ai_formatted_response,omitempty"`
}

// CalculateSynastry calculates synastry between two charts
func (ss *SynastryService) CalculateSynastry(req *SynastryRequest) (*SynastryResponse, error) {
	ss.logger.CalculationLogger().
		Str("person1_city", req.Person1.City).
		Str("person2_city", req.Person2.City).
		Bool("draw_chart", req.DrawChart).
		Msg("🔮 Starting synastry calculation")

	// Calculate natal charts for both people
	person1Chart, err := ss.calculatePersonChart(req.Person1)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate chart for person 1: %w", err)
	}

	person2Chart, err := ss.calculatePersonChart(req.Person2)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate chart for person 2: %w", err)
	}

	// Calculate synastry aspects
	synastryAspects := ss.aspectCalculator.CalculateAspectsBetweenCharts(
		person1Chart.Planets,
		person2Chart.Planets,
	)

	// Create response
	response := &SynastryResponse{
		Person1Chart:    person1Chart,
		Person2Chart:    person2Chart,
		SynastryAspects: synastryAspects,
	}

	// Generate SVG chart if requested
	if req.DrawChart {
		theme := ss.chartDrawer.GetThemeFromString(req.SVGTheme)
		width := req.SVGWidth
		if width <= 0 {
			width = 600
		}

		svg, err := ss.chartDrawer.GenerateSynastryChart(person1Chart, person2Chart, width, theme)
		if err != nil {
			ss.logger.Error().
				Err(err).
				Msg("Failed to generate synastry chart SVG")
		} else {
			response.ChartDraw = svg
		}
	}

	ss.logger.Info().
		Int("synastry_aspects", len(synastryAspects)).
		Msg("✨ Synastry calculation completed successfully")

	return response, nil
}

// calculatePersonChart converts PersonData to a natal chart
func (ss *SynastryService) calculatePersonChart(person PersonData) (*domain.Chart, error) {
	// Convert PersonData to NatalChartRequest
	natalReq := &NatalChartRequest{
		Day:         person.Day,
		Month:       person.Month,
		Year:        person.Year,
		LocalTime:   person.LocalTime,
		City:        person.City,
		HouseSystem: person.HouseSystem,
		DrawChart:   false, // Don't generate SVG for individual charts
		AIResponse:  false,
	}

	// Calculate natal chart
	response, err := ss.natalService.CalculateNatalChart(natalReq)
	if err != nil {
		return nil, err
	}

	// Set name if provided
	if person.Name != "" {
		response.Chart.Name = person.Name
	}

	return response.Chart, nil
}

// GetSynastryFormatted returns formatted synastry for LLM consumption
func (ss *SynastryService) GetSynastryFormatted(req *SynastryRequest) (string, error) {
	response, err := ss.CalculateSynastry(req)
	if err != nil {
		return "", err
	}

	return ss.formatSynastryForLLM(response), nil
}

// formatSynastryForLLM formats synastry results for LLM consumption
func (ss *SynastryService) formatSynastryForLLM(response *SynastryResponse) string {
	formatted := "SYNASTRY ANALYSIS\n\n"

	// Basic information
	formatted += fmt.Sprintf("Person 1: %s\n", response.Person1Chart.Name)
	formatted += fmt.Sprintf("Person 2: %s\n\n", response.Person2Chart.Name)

	// Person 1 Chart Summary
	formatted += "PERSON 1 CHART:\n"
	for _, planet := range response.Person1Chart.Planets {
		formatted += fmt.Sprintf("• %s: %s %s (House %d)\n", planet.Name, planet.Degree, planet.Sign, planet.House)
	}
	formatted += "\n"

	// Person 2 Chart Summary
	formatted += "PERSON 2 CHART:\n"
	for _, planet := range response.Person2Chart.Planets {
		formatted += fmt.Sprintf("• %s: %s %s (House %d)\n", planet.Name, planet.Degree, planet.Sign, planet.House)
	}
	formatted += "\n"

	// Synastry Aspects
	if len(response.SynastryAspects) > 0 {
		formatted += "SYNASTRY ASPECTS:\n"
		for _, aspect := range response.SynastryAspects {
			formatted += fmt.Sprintf("• %s %s %s (%.1f° orb) - %s\n",
				aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb, aspect.Nature)
		}
		formatted += "\n"
	}

	formatted += "SYNASTRY INTERPRETATION:\n"
	formatted += "This synastry analysis shows the astrological connections between these two individuals. "
	formatted += "The aspects between the planets reveal areas of harmony, tension, and growth potential in the relationship. "
	formatted += "Each aspect contributes to the overall dynamic between the two people."

	return formatted
}
