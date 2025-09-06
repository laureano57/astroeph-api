package chart

import (
	"astroeph-api/internal/domain"
	"errors"
)

// ChartRequest represents a request to generate a chart SVG
type ChartRequest struct {
	NatalChartResponse *domain.Chart `json:"natal_chart_response"`
	Width              int           `json:"width"`
	Height             *int          `json:"height,omitempty"`
	Config             *Config       `json:"config,omitempty"`
	ThemeType          *ThemeType    `json:"theme_type,omitempty"`
}

// ChartResponse represents the generated chart response
type ChartResponse struct {
	SVG    string `json:"svg"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// GenerateNatalChartSVG generates an SVG natal chart from chart data
func GenerateNatalChartSVG(request *ChartRequest) (*ChartResponse, error) {
	if request.NatalChartResponse == nil {
		return nil, errors.New("natal chart response is required")
	}

	if request.Width <= 0 {
		request.Width = 600 // Default width
	}

	// Use provided config or create default
	config := DefaultConfig()
	if request.Config != nil {
		config = *request.Config
	}

	// Override theme type if provided
	if request.ThemeType != nil {
		config.ThemeType = *request.ThemeType
	}

	// Create chart data from natal chart response
	chartData := NewChartData(request.NatalChartResponse, config)

	// Create chart generator
	chart := NewChart(chartData, request.Width, request.Height, nil)

	// Generate SVG
	svg := chart.GenerateSVG()

	response := &ChartResponse{
		SVG:    svg,
		Width:  chart.Width,
		Height: chart.Height,
	}

	return response, nil
}

// GenerateCompositeChartSVG generates an SVG composite chart from two chart datasets
func GenerateCompositeChartSVG(chart1, chart2 *domain.Chart, width int, height *int, config *Config) (*ChartResponse, error) {
	if chart1 == nil || chart2 == nil {
		return nil, errors.New("both natal chart responses are required for composite chart")
	}

	if width <= 0 {
		width = 600 // Default width
	}

	// Use provided config or create default
	if config == nil {
		defaultConfig := DefaultConfig()
		config = &defaultConfig
	}

	// Create chart data for both charts
	chartData1 := NewChartData(chart1, *config)
	chartData2 := NewChartData(chart2, *config)

	// Create chart generator for composite
	chart := NewChart(chartData1, width, height, chartData2)

	// Generate SVG
	svg := chart.GenerateSVG()

	response := &ChartResponse{
		SVG:    svg,
		Width:  chart.Width,
		Height: chart.Height,
	}

	return response, nil
}

// GenerateSynastryChartSVG generates an SVG synastry chart (relationship comparison)
func GenerateSynastryChartSVG(chart1, chart2 *domain.Chart, width int, height *int, config *Config) (*ChartResponse, error) {
	// Synastry uses the same visual structure as composite but may have different settings
	return GenerateCompositeChartSVG(chart1, chart2, width, height, config)
}

// GenerateTransitChartSVG generates an SVG transit chart (current planets vs natal)
func GenerateTransitChartSVG(natalChart *domain.Chart, transitPositions *domain.Chart, width int, height *int, config *Config) (*ChartResponse, error) {
	// Transits are technically synastry between natal and current positions
	return GenerateCompositeChartSVG(natalChart, transitPositions, width, height, config)
}

// PrepareProgressionData prepares chart data for secondary progressions
func PrepareProgressionData(natalChart *domain.Chart, progressedDate string, config Config) (*ChartData, error) {
	// Convert natal chart data
	natalData := NewChartData(natalChart, config)

	// Calculate progressed positions (placeholder for now)
	progressedBodies, err := calculateProgressedPositions(natalData.Aspectables, "", progressedDate)
	if err != nil {
		return nil, err
	}

	// Create progressed chart data with same houses but progressed planet positions
	progressedData := &ChartData{
		Signs:       natalData.Signs,                                      // Signs don't change
		Houses:      natalData.Houses,                                     // Houses usually stay the same for progressions
		Aspectables: progressedBodies,                                     // Use progressed planet positions
		Aspects:     calculateAspectsFromBodies(progressedBodies, config), // Recalculate aspects
		Config:      config,
	}

	return progressedData, nil
}

// GetAvailableThemes returns available theme types
func GetAvailableThemes() []ThemeType {
	return []ThemeType{ThemeLight, ThemeDark, ThemeMono}
}

// GetAvailableHouseSystems returns available house systems
func GetAvailableHouseSystems() []HouseSystem {
	return []HouseSystem{
		HousePlacidus,
		HouseKoch,
		HouseEqual,
		HouseCampanus,
		HouseRegiomontanus,
		HousePorphyry,
		HouseWholeSign,
	}
}

// ValidateChartRequest validates a chart generation request
func ValidateChartRequest(request *ChartRequest) error {
	if request == nil {
		return errors.New("chart request is required")
	}

	if request.NatalChartResponse == nil {
		return errors.New("natal chart response is required")
	}

	if request.Width <= 0 {
		return errors.New("width must be greater than 0")
	}

	if request.Height != nil && *request.Height <= 0 {
		return errors.New("height must be greater than 0 if provided")
	}

	// Validate theme type if provided
	if request.ThemeType != nil {
		validThemes := GetAvailableThemes()
		isValid := false
		for _, theme := range validThemes {
			if *request.ThemeType == theme {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("invalid theme type")
		}
	}

	return nil
}

// ChartPresets provides common chart configurations

// GetLightThemeConfig returns a light theme configuration
func GetLightThemeConfig() Config {
	config := DefaultConfig()
	config.ThemeType = ThemeLight
	return config
}

// GetDarkThemeConfig returns a dark theme configuration
func GetDarkThemeConfig() Config {
	config := DefaultConfig()
	config.ThemeType = ThemeDark
	return config
}

// GetMonoThemeConfig returns a monochrome theme configuration
func GetMonoThemeConfig() Config {
	config := DefaultConfig()
	config.ThemeType = ThemeMono
	return config
}

// GetMinimalDisplayConfig returns a configuration with minimal body display
func GetMinimalDisplayConfig() Config {
	config := DefaultConfig()
	config.Display = Display{
		Sun:     true,
		Moon:    true,
		Mercury: true,
		Venus:   true,
		Mars:    true,
		Jupiter: true,
		Saturn:  true,
		Uranus:  false,
		Neptune: false,
		Pluto:   false,
		AscNode: false,
		Chiron:  false,
		Ceres:   false,
		Pallas:  false,
		Juno:    false,
		Vesta:   false,
		Asc:     true,
		IC:      false,
		Dsc:     false,
		MC:      true,
	}
	return config
}

// GetFullDisplayConfig returns a configuration with all bodies displayed
func GetFullDisplayConfig() Config {
	config := DefaultConfig()
	config.Display = Display{
		Sun:     true,
		Moon:    true,
		Mercury: true,
		Venus:   true,
		Mars:    true,
		Jupiter: true,
		Saturn:  true,
		Uranus:  true,
		Neptune: true,
		Pluto:   true,
		AscNode: true,
		Chiron:  true,
		Ceres:   true,
		Pallas:  true,
		Juno:    true,
		Vesta:   true,
		Asc:     true,
		IC:      true,
		Dsc:     true,
		MC:      true,
	}
	return config
}
