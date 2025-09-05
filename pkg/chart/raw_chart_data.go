package chart

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// RawChartData contains raw numeric astrological data
type RawChartData struct {
	Name        string          `json:"name"`
	Lat         float64         `json:"lat"`
	Lon         float64         `json:"lon"`
	UTCTime     time.Time       `json:"utc_time"`
	Planets     []RawPlanetData `json:"planets"`
	HouseCusps  []float64       `json:"house_cusps"`
	Ascendant   float64         `json:"ascendant"`
	Midheaven   float64         `json:"midheaven"`
	HouseSystem string          `json:"house_system"`
}

// RawPlanetData contains raw planet position data
type RawPlanetData struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
}

// GenerateNatalChartSVGFromRawData generates SVG from raw numeric data
func GenerateNatalChartSVGFromRawData(rawData *RawChartData, width int, themeType *ThemeType) (*ChartResponse, error) {
	if rawData == nil {
		return nil, fmt.Errorf("raw chart data is required")
	}

	if width <= 0 {
		width = 600
	}

	// Debug SVG paths loading
	//DebugSVGPaths()
	//DebugSignMembers()

	// Create configuration
	config := DefaultConfig()
	if themeType != nil {
		config.ThemeType = *themeType
	}

	// Convert raw data to ChartData
	chartData := convertRawDataToChartData(rawData, config)

	// Create chart generator
	chart := NewChart(chartData, width, nil, nil)

	// Generate SVG
	svg := chart.GenerateSVG()

	response := &ChartResponse{
		SVG:    svg,
		Width:  chart.Width,
		Height: chart.Height,
	}

	return response, nil
}

// convertRawDataToChartData converts raw data to internal ChartData format
func convertRawDataToChartData(rawData *RawChartData, config Config) *ChartData {
	chartData := &ChartData{
		Name:    rawData.Name,
		Lat:     rawData.Lat,
		Lon:     rawData.Lon,
		UTCTime: rawData.UTCTime.Format("2006-01-02 15:04:05"),
		Config:  config,
	}

	// Create vertices from raw ascendant/midheaven
	chartData.Vertices = []Vertex{
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[0].Body, Degree: rawData.Ascendant}},                          // ASC
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[1].Body, Degree: normalizeAngle360(rawData.Midheaven + 180)}}, // IC
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[2].Body, Degree: normalizeAngle360(rawData.Ascendant + 180)}}, // DSC
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[3].Body, Degree: rawData.Midheaven}},                          // MC
	}

	// Create houses from raw cusps
	chartData.Houses = make([]House, 12)
	for i := 0; i < 12; i++ {
		degree := 0.0
		if i < len(rawData.HouseCusps) {
			degree = rawData.HouseCusps[i]
		}

		chartData.Houses[i] = House{
			MovableBody: MovableBody{
				Body:   HOUSE_MEMBERS[i].Body,
				Degree: degree,
			},
		}
	}

	// Create planets from raw data
	chartData.Planets = make([]Planet, 0, len(rawData.Planets))
	for _, rawPlanet := range rawData.Planets {
		// Find matching planet member
		var planetBody Body
		found := false
		normalizedName := normalizeBodyNameForRaw(rawPlanet.Name)

		for _, member := range PLANET_MEMBERS {
			if normalizeBodyNameForRaw(member.Name) == normalizedName {
				planetBody = member
				found = true
				break
			}
		}

		if !found {
			continue
		}

		// Calculate which house this planet is in
		house := getHouseForLongitude(rawPlanet.Longitude, rawData.HouseCusps)

		planet := Planet{
			MovableBody: MovableBody{
				Body:   planetBody,
				Degree: rawPlanet.Longitude,
				Speed:  rawPlanet.Speed,
				House:  house,
			},
		}

		chartData.Planets = append(chartData.Planets, planet)
	}

	// Create signs - always fixed positions relative to 0° Aries
	chartData.Signs = make([]Sign, 12)
	for i, signMember := range SIGN_MEMBERS {
		chartData.Signs[i] = Sign{
			SignMember: signMember,
			Degree:     float64(i * 30), // 0°, 30°, 60°, etc.
		}
	}

	// Set normalized degrees relative to Ascendant
	ascDegree := rawData.Ascendant

	// Normalize signs
	for i := range chartData.Signs {
		chartData.Signs[i].NormalizedDegree = normalizeRelativeToAsc(chartData.Signs[i].Degree, ascDegree)
	}

	// Normalize planets
	for i := range chartData.Planets {
		chartData.Planets[i].NormalizedDegree = normalizeRelativeToAsc(chartData.Planets[i].Degree, ascDegree)
	}

	// Normalize vertices
	for i := range chartData.Vertices {
		chartData.Vertices[i].NormalizedDegree = normalizeRelativeToAsc(chartData.Vertices[i].Degree, ascDegree)
	}

	// Normalize houses
	for i := range chartData.Houses {
		chartData.Houses[i].NormalizedDegree = normalizeRelativeToAsc(chartData.Houses[i].Degree, ascDegree)
	}

	// Set aspectables
	chartData.Aspectables = make([]MovableBody, 0)

	// Add planets if they should be displayed
	for _, planet := range chartData.Planets {
		if shouldDisplayBodyByName(planet.Name, config.Display) {
			chartData.Aspectables = append(chartData.Aspectables, planet.MovableBody)
		}
	}

	// Add vertices if they should be displayed
	for _, vertex := range chartData.Vertices {
		if shouldDisplayBodyByName(vertex.Name, config.Display) {
			chartData.Aspectables = append(chartData.Aspectables, vertex.MovableBody)
		}
	}

	// Calculate aspects between aspectables
	chartData.Aspects = calculateAspectsFromBodies(chartData.Aspectables, config)

	return chartData
}

// normalizeAngle360 ensures angle is between 0 and 360 degrees
func normalizeAngle360(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// normalizeRelativeToAsc normalizes a degree relative to Ascendant
func normalizeRelativeToAsc(degree, ascDegree float64) float64 {
	return normalizeAngle360(degree - ascDegree)
}

// normalizeBodyNameForRaw normalizes body names for raw data comparison
func normalizeBodyNameForRaw(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "sun":
		return "sun"
	case "moon":
		return "moon"
	case "mercury":
		return "mercury"
	case "venus":
		return "venus"
	case "mars":
		return "mars"
	case "jupiter":
		return "jupiter"
	case "saturn":
		return "saturn"
	case "uranus":
		return "uranus"
	case "neptune":
		return "neptune"
	case "pluto":
		return "pluto"
	case "north node":
		return "asc_node"
	case "chiron":
		return "chiron"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

// shouldDisplayBodyByName checks if a body should be displayed
func shouldDisplayBodyByName(name string, display Display) bool {
	switch normalizeBodyNameForRaw(name) {
	case "sun":
		return display.Sun
	case "moon":
		return display.Moon
	case "mercury":
		return display.Mercury
	case "venus":
		return display.Venus
	case "mars":
		return display.Mars
	case "jupiter":
		return display.Jupiter
	case "saturn":
		return display.Saturn
	case "uranus":
		return display.Uranus
	case "neptune":
		return display.Neptune
	case "pluto":
		return display.Pluto
	case "asc_node":
		return display.AscNode
	case "chiron":
		return display.Chiron
	case "asc":
		return display.Asc
	case "ic":
		return display.IC
	case "dsc":
		return display.Dsc
	case "mc":
		return display.MC
	default:
		return false
	}
}

// getHouseForLongitude determines which house a longitude is in
func getHouseForLongitude(longitude float64, houseCusps []float64) int {
	if len(houseCusps) != 12 {
		return 1 // Default to first house
	}

	normLon := normalizeAngle360(longitude)

	for i := 0; i < 12; i++ {
		currentCusp := normalizeAngle360(houseCusps[i])
		nextCusp := normalizeAngle360(houseCusps[(i+1)%12])

		// Handle wrap-around case
		if currentCusp > nextCusp {
			if normLon >= currentCusp || normLon < nextCusp {
				return i + 1
			}
		} else {
			if normLon >= currentCusp && normLon < nextCusp {
				return i + 1
			}
		}
	}

	return 1 // Default to first house
}

// calculateAspectsFromBodies calculates aspects between aspectable bodies
func calculateAspectsFromBodies(bodies []MovableBody, config Config) []Aspect {
	var aspects []Aspect

	for i, body1 := range bodies {
		for j, body2 := range bodies {
			if i >= j {
				continue // Avoid duplicates and self-aspects
			}

			// Calculate angular difference
			angle := math.Abs(body1.Degree - body2.Degree)
			if angle > 180 {
				angle = 360 - angle
			}

			// Check all aspects
			for _, aspectMember := range ASPECT_MEMBERS {
				orbValue := config.GetOrbForAspect(aspectMember.Name)
				if orbValue == 0 {
					continue
				}

				diff := math.Abs(angle - float64(aspectMember.Value))
				if diff <= float64(orbValue) {
					applying := body1.Speed > body2.Speed
					if angle < float64(aspectMember.Value) {
						applying = !applying
					}

					aspect := Aspect{
						Body1:        &body1,
						Body2:        &body2,
						AspectMember: aspectMember,
						Applying:     &applying,
						Orb:          &diff,
					}
					aspects = append(aspects, aspect)
					break // Only one aspect per pair
				}
			}
		}
	}

	return aspects
}

// calculateCompositeAspects calculates aspects between bodies from two different charts (synastry/transits)
func calculateCompositeAspects(bodies1, bodies2 []MovableBody, config Config) []Aspect {
	var aspects []Aspect

	for _, body1 := range bodies1 {
		for _, body2 := range bodies2 {
			// Calculate angular difference
			angle := math.Abs(body1.Degree - body2.Degree)
			if angle > 180 {
				angle = 360 - angle
			}

			// Check all aspects
			for _, aspectMember := range ASPECT_MEMBERS {
				orbValue := config.GetOrbForAspect(aspectMember.Name)
				if orbValue == 0 {
					continue // Aspect not enabled
				}

				diff := math.Abs(angle - float64(aspectMember.Value))
				if diff <= float64(orbValue) {
					applying := body1.Speed > body2.Speed
					if angle < float64(aspectMember.Value) {
						applying = !applying
					}

					aspect := Aspect{
						Body1:        &body1,
						Body2:        &body2,
						AspectMember: aspectMember,
						Applying:     &applying,
						Orb:          &diff,
					}
					aspects = append(aspects, aspect)
					break // Only one aspect per pair
				}
			}
		}
	}

	return aspects
}

// calculateProgressedPositions calculates secondary progressed positions for a given date
func calculateProgressedPositions(natalBodies []MovableBody, natalDate, progressedDate string) ([]MovableBody, error) {
	// TODO: Implement secondary progressions calculation
	// This is a placeholder that would be expanded when implementing progressions
	// Secondary progressions: 1 day = 1 year
	progressedBodies := make([]MovableBody, len(natalBodies))
	copy(progressedBodies, natalBodies)

	// For now, return natal positions unchanged
	// Real implementation would calculate progressed positions based on date difference
	return progressedBodies, nil
}
