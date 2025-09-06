package domain

import (
	"time"
)

// ChartType represents the type of astrological chart
type ChartType string

const (
	ChartTypeNatal        ChartType = "natal"
	ChartTypeSynastry     ChartType = "synastry"
	ChartTypeComposite    ChartType = "composite"
	ChartTypeSolarReturn  ChartType = "solar_return"
	ChartTypeLunarReturn  ChartType = "lunar_return"
	ChartTypeProgressions ChartType = "progressions"
	ChartTypeTransits     ChartType = "transits"
)

// Chart represents a complete astrological chart
type Chart struct {
	ID          string      `json:"id"`
	Type        ChartType   `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	BirthInfo   BirthInfo   `json:"birth_info"`
	Planets     []Planet    `json:"planets"`
	Houses      []House     `json:"houses"`
	Aspects     []Aspect    `json:"aspects"`
	Angles      ChartAngles `json:"angles"`
	HouseSystem string      `json:"house_system"`
	Timezone    string      `json:"timezone"`
	UTCTime     time.Time   `json:"utc_time"`
	ChartDraw   string      `json:"chart_draw,omitempty"` // SVG chart
	CreatedAt   time.Time   `json:"created_at"`
}

// ChartAngles represents the main chart angles
type ChartAngles struct {
	Ascendant  ChartAngle `json:"ascendant"`
	Midheaven  ChartAngle `json:"midheaven"`
	IC         ChartAngle `json:"ic,omitempty"`
	Descendant ChartAngle `json:"descendant,omitempty"`
}

// ChartAngle represents an important chart angle
type ChartAngle struct {
	Sign   string  `json:"sign"`
	Degree string  `json:"degree"`
	Value  float64 `json:"value"` // Raw degree value
}

// BirthInfo represents birth information for a chart
type BirthInfo struct {
	Date     string   `json:"date"`
	Time     string   `json:"time"`
	Location Location `json:"location"`
}

// NewChart creates a new chart with the given parameters
func NewChart(chartType ChartType, name string, birthInfo BirthInfo) *Chart {
	return &Chart{
		Type:      chartType,
		Name:      name,
		BirthInfo: birthInfo,
		CreatedAt: time.Now(),
		Planets:   make([]Planet, 0),
		Houses:    make([]House, 0),
		Aspects:   make([]Aspect, 0),
	}
}

// AddPlanet adds a planet to the chart
func (c *Chart) AddPlanet(planet Planet) {
	c.Planets = append(c.Planets, planet)
}

// AddHouse adds a house to the chart
func (c *Chart) AddHouse(house House) {
	c.Houses = append(c.Houses, house)
}

// AddAspect adds an aspect to the chart
func (c *Chart) AddAspect(aspect Aspect) {
	c.Aspects = append(c.Aspects, aspect)
}

// SetAngles sets the main chart angles
func (c *Chart) SetAngles(ascendant, midheaven float64) {
	c.Angles.Ascendant = ChartAngle{
		Sign:   GetZodiacSign(ascendant),
		Degree: FormatDegreeInSign(ascendant),
		Value:  ascendant,
	}

	c.Angles.Midheaven = ChartAngle{
		Sign:   GetZodiacSign(midheaven),
		Degree: FormatDegreeInSign(midheaven),
		Value:  midheaven,
	}

	// Calculate IC and Descendant
	ic := normalizeAngle(midheaven + 180)
	descendant := normalizeAngle(ascendant + 180)

	c.Angles.IC = ChartAngle{
		Sign:   GetZodiacSign(ic),
		Degree: FormatDegreeInSign(ic),
		Value:  ic,
	}

	c.Angles.Descendant = ChartAngle{
		Sign:   GetZodiacSign(descendant),
		Degree: FormatDegreeInSign(descendant),
		Value:  descendant,
	}
}

// GetPlanetByName returns a planet by its name
func (c *Chart) GetPlanetByName(name string) *Planet {
	for i, planet := range c.Planets {
		if planet.Name == name {
			return &c.Planets[i]
		}
	}
	return nil
}

// GetHouseByNumber returns a house by its number
func (c *Chart) GetHouseByNumber(number int) *House {
	for i, house := range c.Houses {
		if house.Number == number {
			return &c.Houses[i]
		}
	}
	return nil
}
