package astro

import (
	"astroeph-api/internal/domain"
	"astroeph-api/pkg/chart"
	"time"
)

// ChartDrawer handles SVG chart generation
type ChartDrawer struct {
	// Chart generation settings
	defaultWidth int
	defaultTheme chart.ThemeType
}

// NewChartDrawer creates a new chart drawer
func NewChartDrawer() *ChartDrawer {
	return &ChartDrawer{
		defaultWidth: 600,
		defaultTheme: chart.ThemeLight,
	}
}

// GenerateNatalChart generates an SVG chart for a natal chart
func (cd *ChartDrawer) GenerateNatalChart(
	natalChart *domain.Chart,
	width int,
	themeType *chart.ThemeType,
) (string, error) {

	if width <= 0 {
		width = cd.defaultWidth
	}

	theme := cd.defaultTheme
	if themeType != nil {
		theme = *themeType
	}

	// Convert domain chart to raw chart data format expected by pkg/chart
	rawData := cd.convertToRawChartData(natalChart)

	// Generate SVG using existing chart library
	response, err := chart.GenerateNatalChartSVGFromRawData(rawData, width, &theme)
	if err != nil {
		return "", err
	}

	return response.SVG, nil
}

// GenerateCompositeChart generates an SVG chart for a composite chart
func (cd *ChartDrawer) GenerateCompositeChart(
	chart1, chart2 *domain.Chart,
	width int,
	themeType *chart.ThemeType,
) (string, error) {

	if width <= 0 {
		width = cd.defaultWidth
	}

	theme := cd.defaultTheme
	if themeType != nil {
		theme = *themeType
	}

	// For composite charts, we would need to calculate midpoints
	// This is a simplified implementation
	compositeChart := cd.calculateCompositeChart(chart1, chart2)
	rawData := cd.convertToRawChartData(compositeChart)

	// Generate SVG
	response, err := chart.GenerateNatalChartSVGFromRawData(rawData, width, &theme)
	if err != nil {
		return "", err
	}

	return response.SVG, nil
}

// GenerateSynastryChart generates an SVG chart for synastry
func (cd *ChartDrawer) GenerateSynastryChart(
	chart1, chart2 *domain.Chart,
	width int,
	themeType *chart.ThemeType,
) (string, error) {

	// For synastry, we typically show both charts with aspects between them
	// This would require a more complex implementation
	// For now, return the first chart with a note about synastry aspects
	return cd.GenerateNatalChart(chart1, width, themeType)
}

// convertToRawChartData converts a domain chart to the format expected by pkg/chart
func (cd *ChartDrawer) convertToRawChartData(domainChart *domain.Chart) *chart.RawChartData {
	// Convert planets
	var rawPlanets []chart.RawPlanetData
	for _, planet := range domainChart.Planets {
		rawPlanet := chart.RawPlanetData{
			Name:      planet.Name,
			Longitude: planet.Longitude,
			Speed:     planet.Speed,
		}
		rawPlanets = append(rawPlanets, rawPlanet)
	}

	// Convert house cusps
	var houseCusps []float64
	for _, house := range domainChart.Houses {
		houseCusps = append(houseCusps, house.CuspValue)
	}

	// Parse UTC time
	utcTime := domainChart.UTCTime
	if utcTime.IsZero() {
		utcTime = time.Now().UTC()
	}

	rawData := &chart.RawChartData{
		Name:        domainChart.Name,
		Lat:         domainChart.BirthInfo.Location.Latitude,
		Lon:         domainChart.BirthInfo.Location.Longitude,
		UTCTime:     utcTime,
		Planets:     rawPlanets,
		HouseCusps:  houseCusps,
		Ascendant:   domainChart.Angles.Ascendant.Value,
		Midheaven:   domainChart.Angles.Midheaven.Value,
		HouseSystem: domainChart.HouseSystem,
	}

	return rawData
}

// calculateCompositeChart calculates a composite chart from two natal charts
func (cd *ChartDrawer) calculateCompositeChart(chart1, chart2 *domain.Chart) *domain.Chart {
	// Create new composite chart
	composite := domain.NewChart(
		domain.ChartTypeComposite,
		"Composite: "+chart1.Name+" & "+chart2.Name,
		chart1.BirthInfo, // Use first chart's birth info as base
	)

	// Calculate midpoint planets
	var compositePlanets []domain.Planet
	for i, planet1 := range chart1.Planets {
		if i < len(chart2.Planets) {
			planet2 := chart2.Planets[i]
			if planet1.Name == planet2.Name {
				// Calculate midpoint longitude
				midpointLon := cd.calculateMidpoint(planet1.Longitude, planet2.Longitude)

				// Create composite planet
				compositePlanet := domain.NewPlanet(
					planet1.Name,
					midpointLon,
					0, // Latitude not needed for basic composite
					0, // Speed not applicable for composite
					1, // House will be recalculated
				)
				compositePlanets = append(compositePlanets, compositePlanet)
			}
		}
	}

	// Calculate midpoint houses
	var compositeHouses []domain.House
	for i, house1 := range chart1.Houses {
		if i < len(chart2.Houses) {
			house2 := chart2.Houses[i]
			if house1.Number == house2.Number {
				// Calculate midpoint cusp
				midpointCusp := cd.calculateMidpoint(house1.CuspValue, house2.CuspValue)

				compositeHouse := domain.NewHouse(house1.Number, midpointCusp)
				compositeHouses = append(compositeHouses, compositeHouse)
			}
		}
	}

	// Calculate midpoint angles
	ascMidpoint := cd.calculateMidpoint(chart1.Angles.Ascendant.Value, chart2.Angles.Ascendant.Value)
	mcMidpoint := cd.calculateMidpoint(chart1.Angles.Midheaven.Value, chart2.Angles.Midheaven.Value)

	// Set composite data
	for _, planet := range compositePlanets {
		composite.AddPlanet(planet)
	}

	for _, house := range compositeHouses {
		composite.AddHouse(house)
	}

	composite.SetAngles(ascMidpoint, mcMidpoint)
	composite.HouseSystem = chart1.HouseSystem

	return composite
}

// calculateMidpoint calculates the midpoint between two longitudes
func (cd *ChartDrawer) calculateMidpoint(lon1, lon2 float64) float64 {
	// Normalize angles to 0-360 range
	lon1 = cd.normalizeAngle(lon1)
	lon2 = cd.normalizeAngle(lon2)

	// Calculate the shorter arc between the two points
	diff := lon2 - lon1
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}

	midpoint := lon1 + diff/2
	return cd.normalizeAngle(midpoint)
}

// normalizeAngle normalizes an angle to 0-360 range
func (cd *ChartDrawer) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// SetDefaultWidth sets the default chart width
func (cd *ChartDrawer) SetDefaultWidth(width int) {
	cd.defaultWidth = width
}

// SetDefaultTheme sets the default chart theme
func (cd *ChartDrawer) SetDefaultTheme(theme chart.ThemeType) {
	cd.defaultTheme = theme
}

// GetAvailableThemes returns available chart themes
func (cd *ChartDrawer) GetAvailableThemes() []chart.ThemeType {
	return []chart.ThemeType{
		chart.ThemeLight,
		chart.ThemeDark,
		chart.ThemeMono,
	}
}

// GetThemeFromString converts string to theme type
func (cd *ChartDrawer) GetThemeFromString(themeStr string) *chart.ThemeType {
	switch themeStr {
	case "light":
		theme := chart.ThemeLight
		return &theme
	case "dark":
		theme := chart.ThemeDark
		return &theme
	case "mono":
		theme := chart.ThemeMono
		return &theme
	default:
		return nil
	}
}

// ChartOptions represents options for chart generation
type ChartOptions struct {
	Width            int              `json:"width"`
	Theme            *chart.ThemeType `json:"theme"`
	ShowAspects      bool             `json:"show_aspects"`
	ShowHouseNumbers bool             `json:"show_house_numbers"`
	ShowDegrees      bool             `json:"show_degrees"`
}

// GenerateChartWithOptions generates a chart with specific options
func (cd *ChartDrawer) GenerateChartWithOptions(
	domainChart *domain.Chart,
	options ChartOptions,
) (string, error) {

	// Apply default options
	if options.Width <= 0 {
		options.Width = cd.defaultWidth
	}

	if options.Theme == nil {
		theme := cd.defaultTheme
		options.Theme = &theme
	}

	// For now, use the basic generation method
	// In a full implementation, these options would affect the chart generation
	return cd.GenerateNatalChart(domainChart, options.Width, options.Theme)
}
