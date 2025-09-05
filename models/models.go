package models

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Common request fields used across multiple endpoints
type BaseRequest struct {
	Day        int    `json:"day" binding:"required,min=1,max=31"`
	Month      int    `json:"month" binding:"required,min=1,max=12"`
	Year       int    `json:"year" binding:"required"`
	LocalTime  string `json:"local_time" binding:"required"` // HH:MM:SS format
	City       string `json:"city" binding:"required"`
	AIResponse bool   `json:"ai_response,omitempty"`
}

// Natal Chart specific requests and responses
type NatalChartRequest struct {
	BaseRequest
	HouseSystem string `json:"house_system,omitempty"` // defaults to "Placidus"
	DrawChart   bool   `json:"draw_chart,omitempty"`   // whether to generate SVG chart
	SVGWidth    int    `json:"svg_width,omitempty"`    // width of SVG chart (defaults to 600)
	SVGTheme    string `json:"svg_theme,omitempty"`    // theme for SVG chart ("light", "dark", "mono")
}

type PlanetPosition struct {
	Name   string `json:"name"`
	Sign   string `json:"sign"`
	Degree string `json:"degree"`
	House  int    `json:"house"`
}

type HouseCusp struct {
	House int    `json:"house"`
	Cusp  string `json:"cusp"`
	Sign  string `json:"sign"`
}

type ChartAngle struct {
	Sign   string `json:"sign"`
	Degree string `json:"degree"`
}

type Aspect struct {
	Planet1    string `json:"planet1"`
	Planet2    string `json:"planet2"`
	Type       string `json:"type"` // e.g., "conjunction", "opposition", etc.
	Orb        string `json:"orb"`
	IsApplying bool   `json:"is_applying"`
}

type ChartData struct {
	Planets     []PlanetPosition `json:"planets"`
	Houses      []HouseCusp      `json:"houses"`
	Aspects     []Aspect         `json:"aspects"`
	Ascendant   ChartAngle       `json:"ascendant"`
	Midheaven   ChartAngle       `json:"midheaven"`
	HouseSystem string           `json:"house_system"`
}

type NatalChartResponse struct {
	ChartData
	BirthInfo BirthInfo `json:"birth_info"`
	Timezone  string    `json:"timezone"`
	UTCTime   time.Time `json:"utc_time"`
	ChartDraw string    `json:"chart_draw,omitempty"` // Generated SVG chart
}

type BirthInfo struct {
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Transits request and response
type TransitsRequest struct {
	BaseRequest
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`
}

type TransitAspect struct {
	TransitPlanet string  `json:"transit_planet"`
	NatalPlanet   string  `json:"natal_planet"`
	Type          string  `json:"type"`
	Angle         float64 `json:"angle"`
	Orb           float64 `json:"orb"`
	IsExact       bool    `json:"is_exact"`
}

type TransitsResponse struct {
	TransitDate    time.Time        `json:"transit_date"`
	TransitPlanets []PlanetPosition `json:"transit_planets"`
	Aspects        []TransitAspect  `json:"aspects"`
	BirthChart     ChartData        `json:"birth_chart"`
}

// Synastry request and response
type SynastryRequest struct {
	Person1    PersonData `json:"person1" binding:"required"`
	Person2    PersonData `json:"person2" binding:"required"`
	AIResponse bool       `json:"ai_response,omitempty"`
}

type PersonData struct {
	Day       int    `json:"day" binding:"required,min=1,max=31"`
	Month     int    `json:"month" binding:"required,min=1,max=12"`
	Year      int    `json:"year" binding:"required"`
	LocalTime string `json:"local_time" binding:"required"`
	City      string `json:"city" binding:"required"`
	Name      string `json:"name,omitempty"`
}

type SynastryAspect struct {
	Person1Planet string  `json:"person1_planet"`
	Person2Planet string  `json:"person2_planet"`
	Type          string  `json:"type"`
	Angle         float64 `json:"angle"`
	Orb           float64 `json:"orb"`
}

type SynastryResponse struct {
	Person1Chart  ChartData        `json:"person1_chart"`
	Person2Chart  ChartData        `json:"person2_chart"`
	Aspects       []SynastryAspect `json:"aspects"`
	Compatibility string           `json:"compatibility,omitempty"`
}

// Additional endpoint requests (to be implemented in later phases)
type CompositeChartRequest struct {
	Person1    PersonData `json:"person1" binding:"required"`
	Person2    PersonData `json:"person2" binding:"required"`
	AIResponse bool       `json:"ai_response,omitempty"`
}

type ProgressionsRequest struct {
	BaseRequest
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`
}

type SolarReturnRequest struct {
	BaseRequest
	BirthDay   int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear  int    `json:"birth_year" binding:"required"`
	BirthTime  string `json:"birth_time" binding:"required"`
	BirthCity  string `json:"birth_city" binding:"required"`
	ReturnYear int    `json:"return_year" binding:"required"`
}

type LunarReturnRequest struct {
	BaseRequest
	BirthDay    int    `json:"birth_day" binding:"required,min=1,max=31"`
	BirthMonth  int    `json:"birth_month" binding:"required,min=1,max=12"`
	BirthYear   int    `json:"birth_year" binding:"required"`
	BirthTime   string `json:"birth_time" binding:"required"`
	BirthCity   string `json:"birth_city" binding:"required"`
	ReturnMonth int    `json:"return_month" binding:"required,min=1,max=12"`
	ReturnYear  int    `json:"return_year" binding:"required"`
}

// Utility functions
func GetZodiacSign(longitude float64) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	signIndex := int(longitude / 30.0)
	if signIndex >= 12 {
		signIndex = 11
	}
	return signs[signIndex]
}

func GetDegreeInSign(longitude float64) float64 {
	return longitude - (float64(int(longitude/30.0)) * 30.0)
}

// formatDegreesMinutes converts decimal degrees to degrees and minutes format
func formatDegreesMinutes(decimalDegrees float64) string {
	degrees := int(decimalDegrees)
	minutes := int((decimalDegrees - float64(degrees)) * 60)
	return fmt.Sprintf("%d°%02d'", degrees, minutes)
}

// formatDegreesMinutesSeconds converts decimal degrees to degrees, minutes and seconds format
func formatDegreesMinutesSeconds(decimalDegrees float64) string {
	degrees := int(decimalDegrees)
	remainingMinutes := (decimalDegrees - float64(degrees)) * 60
	minutes := int(remainingMinutes)
	seconds := int((remainingMinutes - float64(minutes)) * 60)
	return fmt.Sprintf("%d°%02d'%02d\"", degrees, minutes, seconds)
}

// FormatDegreeInSign converts decimal degrees within a sign to DMS format
func FormatDegreeInSign(longitude float64) string {
	degreeInSign := GetDegreeInSign(longitude)
	return formatDegreesMinutesSeconds(degreeInSign)
}

// FormatLongitude converts a full longitude to DMS format (for house cusps)
func FormatLongitude(longitude float64) string {
	return formatDegreesMinutesSeconds(longitude)
}

// FormatNatalChartForLLM converts natal chart data to LLM-friendly text format
func FormatNatalChartForLLM(chartData *NatalChartResponse) string {
	var result strings.Builder

	// Header
	result.WriteString(fmt.Sprintf("NATAL CHART ANALYSIS\n"))
	result.WriteString(fmt.Sprintf("Birth Date: %s at %s\n", chartData.BirthInfo.Date, chartData.BirthInfo.Time))
	// Format coordinates in degrees and minutes
	latDirection := "N"
	lonDirection := "E"
	if chartData.BirthInfo.Latitude < 0 {
		latDirection = "S"
	}
	if chartData.BirthInfo.Longitude < 0 {
		lonDirection = "W"
	}

	latDegMin := formatDegreesMinutes(math.Abs(chartData.BirthInfo.Latitude))
	lonDegMin := formatDegreesMinutes(math.Abs(chartData.BirthInfo.Longitude))

	result.WriteString(fmt.Sprintf("Birth Location: %s (%s%s, %s%s)\n",
		chartData.BirthInfo.City, latDegMin, latDirection, lonDegMin, lonDirection))
	result.WriteString(fmt.Sprintf("House System: %s\n\n", chartData.HouseSystem))

	// Planetary Positions
	result.WriteString("PLANETARY POSITIONS:\n")
	for _, planet := range chartData.Planets {
		result.WriteString(fmt.Sprintf("• %s: %s %s (House %d)\n",
			planet.Name, planet.Degree, planet.Sign, planet.House))
	}

	// Chart Angles
	result.WriteString("\nCHART ANGLES:\n")
	result.WriteString(fmt.Sprintf("• Ascendant: %s %s\n", chartData.Ascendant.Degree, chartData.Ascendant.Sign))
	result.WriteString(fmt.Sprintf("• Midheaven: %s %s\n", chartData.Midheaven.Degree, chartData.Midheaven.Sign))

	// House Cusps
	result.WriteString("\nHOUSE CUSPS:\n")
	for _, house := range chartData.Houses {
		result.WriteString(fmt.Sprintf("• House %d: %s %s\n", house.House, house.Cusp, house.Sign))
	}

	// Major Aspects
	if len(chartData.Aspects) > 0 {
		result.WriteString("\nMAJOR ASPECTS:\n")
		for _, aspect := range chartData.Aspects {
			result.WriteString(fmt.Sprintf("• %s %s %s - %s orb\n",
				aspect.Planet1, aspect.Type, aspect.Planet2, aspect.Orb))
		}
	}

	// Astrological Summary
	result.WriteString(fmt.Sprintf("\nASTROLOGICAL SUMMARY:\n"))
	result.WriteString(fmt.Sprintf("This natal chart shows %d planetary positions across %d houses, with %d major aspects formed. ",
		len(chartData.Planets), len(chartData.Houses), len(chartData.Aspects)))

	// Sign emphasis
	signCount := make(map[string]int)
	for _, planet := range chartData.Planets {
		signCount[planet.Sign]++
	}

	maxCount := 0
	dominantSign := ""
	for sign, count := range signCount {
		if count > maxCount {
			maxCount = count
			dominantSign = sign
		}
	}

	if maxCount > 1 {
		result.WriteString(fmt.Sprintf("There is a notable emphasis in %s with %d planetary placements. ", dominantSign, maxCount))
	}

	result.WriteString("This chart provides a comprehensive astrological foundation for interpretation.")

	return result.String()
}

// FormatForLLM provides backward compatibility - delegates to specific formatters
func FormatForLLM(chartData interface{}) string {
	switch data := chartData.(type) {
	case *NatalChartResponse:
		return FormatNatalChartForLLM(data)
	default:
		return "Astrological data formatted for LLM analysis - specific formatting not yet implemented for this data type."
	}
}
