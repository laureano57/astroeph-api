package chart

import (
	"astroeph-api/models"
	"math"
	"strings"
)

// MovableBody represents a celestial body with position and movement
type MovableBody struct {
	Body
	Degree           float64 `json:"degree"`
	Speed            float64 `json:"speed"`
	NormalizedDegree float64 `json:"normalized_degree"`
	House            int     `json:"house,omitempty"`
}

// SignedDeg returns the degree within the current sign (0-29)
func (mb *MovableBody) SignedDeg() int {
	return int(math.Mod(mb.Degree, 30))
}

// Minute returns the arc minutes of the position
func (mb *MovableBody) Minute() int {
	minutes := (math.Mod(mb.Degree, 30) - float64(mb.SignedDeg())) * 60
	return int(math.Floor(minutes))
}

// IsRetrograde returns true if the body is moving in retrograde motion
func (mb *MovableBody) IsRetrograde() bool {
	return mb.Speed < 0
}

// GetSign returns the zodiac sign for this position
func (mb *MovableBody) GetSign() SignMember {
	idx := int(mb.Degree / 30.0)
	if idx >= 12 {
		idx = 11
	}
	if idx < 0 {
		idx = 0
	}
	return SIGN_MEMBERS[idx]
}

// Planet represents a planet in the chart
type Planet struct {
	MovableBody
}

// Vertex represents a chart vertex (ASC, IC, DSC, MC)
type Vertex struct {
	MovableBody
}

// House represents a house in the chart
type House struct {
	MovableBody
	Ruler             string `json:"ruler,omitempty"`
	RulerSign         string `json:"ruler_sign,omitempty"`
	RulerHouse        int    `json:"ruler_house,omitempty"`
	ClassicRuler      string `json:"classic_ruler,omitempty"`
	ClassicRulerSign  string `json:"classic_ruler_sign,omitempty"`
	ClassicRulerHouse int    `json:"classic_ruler_house,omitempty"`
}

// Sign represents a zodiac sign with its position
type Sign struct {
	SignMember
	Degree           float64 `json:"degree"`
	NormalizedDegree float64 `json:"normalized_degree"`
}

// Aspect represents an aspect between two celestial bodies
type Aspect struct {
	Body1        *MovableBody `json:"body1"`
	Body2        *MovableBody `json:"body2"`
	AspectMember AspectMember `json:"aspect_member"`
	Applying     *bool        `json:"applying,omitempty"`
	Orb          *float64     `json:"orb,omitempty"`
}

// ChartData contains all the astrological data for generating a chart
type ChartData struct {
	Name        string        `json:"name"`
	Lat         float64       `json:"lat"`
	Lon         float64       `json:"lon"`
	UTCTime     string        `json:"utc_time"`
	Config      Config        `json:"config"`
	Houses      []House       `json:"houses"`
	Planets     []Planet      `json:"planets"`
	Vertices    []Vertex      `json:"vertices"`
	Signs       []Sign        `json:"signs"`
	Aspects     []Aspect      `json:"aspects"`
	Aspectables []MovableBody `json:"aspectables"`
}

// NewChartData creates a new ChartData from models.NatalChartResponse
func NewChartData(response *models.NatalChartResponse, config Config) *ChartData {
	chartData := &ChartData{
		Name:    response.BirthInfo.City,
		Lat:     response.BirthInfo.Latitude,
		Lon:     response.BirthInfo.Longitude,
		UTCTime: response.UTCTime.Format("2006-01-02 15:04:05"),
		Config:  config,
	}

	// Create vertices from chart angles
	chartData.createVertices(response)

	// Create houses from house cusps
	chartData.createHouses(response)

	// Create planets from planet positions
	chartData.createPlanets(response)

	// Create signs
	chartData.createSigns()

	// Set normalized degrees
	chartData.setNormalizedDegrees()

	// Create aspectables list
	chartData.setAspectables()

	// Create aspects
	chartData.createAspects(response)

	return chartData
}

// createVertices creates vertex objects from chart angles
func (cd *ChartData) createVertices(response *models.NatalChartResponse) {
	// Get ASC degree from first house cusp
	ascDeg := 0.0
	mcDeg := 0.0

	if len(response.Houses) > 0 {
		ascDeg = parseDegreeDMS(response.Houses[0].Cusp)
	}
	if len(response.Houses) >= 10 {
		mcDeg = parseDegreeDMS(response.Houses[9].Cusp) // 10th house is MC
	}

	cd.Vertices = []Vertex{
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[0].Body, Degree: ascDeg}},
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[1].Body, Degree: math.Mod(mcDeg+180, 360)}},  // IC
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[2].Body, Degree: math.Mod(ascDeg+180, 360)}}, // DSC
		{MovableBody: MovableBody{Body: VERTEX_MEMBERS[3].Body, Degree: mcDeg}},                     // MC
	}
}

// createHouses creates house objects from house cusps
func (cd *ChartData) createHouses(response *models.NatalChartResponse) {
	cd.Houses = make([]House, len(response.Houses))

	for i, houseCusp := range response.Houses {
		degree := parseDegreeDMS(houseCusp.Cusp)
		houseBody := HOUSE_MEMBERS[i]

		cd.Houses[i] = House{
			MovableBody: MovableBody{
				Body:   houseBody.Body,
				Degree: degree,
			},
		}
	}
}

// createPlanets creates planet objects from planet positions
func (cd *ChartData) createPlanets(response *models.NatalChartResponse) {
	cd.Planets = make([]Planet, 0, len(response.Planets))

	for _, planetPos := range response.Planets {
		// Find matching planet member
		var planetBody Body
		found := false
		for _, member := range PLANET_MEMBERS {
			if normalizeBodyName(member.Name) == normalizeBodyName(planetPos.Name) {
				planetBody = member
				found = true
				break
			}
		}

		if !found {
			continue // Skip unknown planets
		}

		degree := parsePlanetPosition(planetPos.Degree, planetPos.Sign)

		planet := Planet{
			MovableBody: MovableBody{
				Body:   planetBody,
				Degree: degree,
				House:  planetPos.House,
			},
		}

		cd.Planets = append(cd.Planets, planet)
	}
}

// createSigns creates sign objects
func (cd *ChartData) createSigns() {
	cd.Signs = make([]Sign, 12)

	for i, signMember := range SIGN_MEMBERS {
		cd.Signs[i] = Sign{
			SignMember: signMember,
			Degree:     float64(i * 30),
		}
	}
}

// setNormalizedDegrees sets normalized degrees relative to Ascendant
func (cd *ChartData) setNormalizedDegrees() {
	if len(cd.Vertices) == 0 {
		return
	}

	ascDegree := cd.Vertices[0].Degree // ASC is first vertex

	// Normalize signs
	for i := range cd.Signs {
		cd.Signs[i].NormalizedDegree = cd.normalize(cd.Signs[i].Degree, ascDegree)
	}

	// Normalize planets
	for i := range cd.Planets {
		cd.Planets[i].NormalizedDegree = cd.normalize(cd.Planets[i].Degree, ascDegree)
	}

	// Normalize vertices
	for i := range cd.Vertices {
		cd.Vertices[i].NormalizedDegree = cd.normalize(cd.Vertices[i].Degree, ascDegree)
	}

	// Normalize houses
	for i := range cd.Houses {
		cd.Houses[i].NormalizedDegree = cd.normalize(cd.Houses[i].Degree, ascDegree)
	}
}

// setAspectables creates the list of aspectable bodies
func (cd *ChartData) setAspectables() {
	cd.Aspectables = make([]MovableBody, 0)

	// Add planets if they should be displayed
	for _, planet := range cd.Planets {
		if cd.shouldDisplayBody(planet.Name) {
			cd.Aspectables = append(cd.Aspectables, planet.MovableBody)
		}
	}

	// Add vertices if they should be displayed
	for _, vertex := range cd.Vertices {
		if cd.shouldDisplayBody(vertex.Name) {
			cd.Aspectables = append(cd.Aspectables, vertex.MovableBody)
		}
	}
}

// createAspects creates aspect objects from model aspects
func (cd *ChartData) createAspects(response *models.NatalChartResponse) {
	cd.Aspects = make([]Aspect, 0, len(response.Aspects))

	for _, modelAspect := range response.Aspects {
		// Find the two bodies in our aspectables list
		var body1, body2 *MovableBody

		for i := range cd.Aspectables {
			bodyName := normalizeBodyName(cd.Aspectables[i].Name)
			if bodyName == normalizeBodyName(modelAspect.Planet1) {
				body1 = &cd.Aspectables[i]
			}
			if bodyName == normalizeBodyName(modelAspect.Planet2) {
				body2 = &cd.Aspectables[i]
			}
		}

		if body1 == nil || body2 == nil {
			continue // Skip if we can't find both bodies
		}

		// Find aspect member
		aspectMember := GetAspectMember(modelAspect.Type)
		if aspectMember == nil {
			continue
		}

		// Parse orb
		orb := parseOrbValue(modelAspect.Orb)

		aspect := Aspect{
			Body1:        body1,
			Body2:        body2,
			AspectMember: *aspectMember,
			Applying:     &modelAspect.IsApplying,
			Orb:          &orb,
		}

		cd.Aspects = append(cd.Aspects, aspect)
	}
}

// normalize normalizes a degree relative to the Ascendant
func (cd *ChartData) normalize(degree, ascDegree float64) float64 {
	return math.Mod(degree-ascDegree+360, 360)
}

// shouldDisplayBody checks if a body should be displayed based on configuration
func (cd *ChartData) shouldDisplayBody(name string) bool {
	display := cd.Config.Display

	switch normalizeBodyName(name) {
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
	case "asc_node", "north_node":
		return display.AscNode
	case "chiron":
		return display.Chiron
	case "ceres":
		return display.Ceres
	case "pallas":
		return display.Pallas
	case "juno":
		return display.Juno
	case "vesta":
		return display.Vesta
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

// Helper functions for parsing model data

// parseDegreeDMS parses a degree string in DMS format to decimal degrees
func parseDegreeDMS(dmsStr string) float64 {
	// For now, return 0 - this would need proper DMS parsing
	// The existing service provides calculated positions already
	return 0.0
}

// parsePlanetPosition calculates the absolute degree from sign and degree within sign
func parsePlanetPosition(degreeStr, signStr string) float64 {
	// Find the sign index
	signIndex := 0
	signName := strings.ToLower(signStr)
	for i, sign := range SIGN_MEMBERS {
		if strings.ToLower(sign.Name) == signName {
			signIndex = i
			break
		}
	}

	// For now, use a placeholder degree within sign
	// TODO: Parse the actual degree string (e.g., "15°23'45\"")
	degreeInSign := 15.0

	return float64(signIndex*30) + degreeInSign
}

// parseOrbValue parses orb string to float64
func parseOrbValue(orbStr string) float64 {
	// Placeholder implementation - would need proper parsing
	return 1.0
}

// normalizeBodyName normalizes body names for comparison
func normalizeBodyName(name string) string {
	switch name {
	case "North Node":
		return "asc_node"
	default:
		return name
	}
}

// normalizeSignName normalizes sign names for comparison
func normalizeSignName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
