package domain

import (
	"fmt"
	"math"
)

// Location represents a geographic location
type Location struct {
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Region    string  `json:"region,omitempty"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Elevation float64 `json:"elevation,omitempty"` // In meters above sea level
}

// NewLocation creates a new Location
func NewLocation(name, city, country string, lat, lon float64, timezone string) *Location {
	return &Location{
		Name:      name,
		City:      city,
		Country:   country,
		Latitude:  lat,
		Longitude: lon,
		Timezone:  timezone,
	}
}

// IsValidCoordinates checks if the coordinates are valid
func (l Location) IsValidCoordinates() bool {
	return l.Latitude >= -90 && l.Latitude <= 90 &&
		l.Longitude >= -180 && l.Longitude <= 180
}

// FormatLatitude formats latitude for display with direction
func (l Location) FormatLatitude() string {
	direction := "N"
	lat := l.Latitude
	if lat < 0 {
		direction = "S"
		lat = -lat
	}

	degrees := int(lat)
	minutes := (lat - float64(degrees)) * 60

	return fmt.Sprintf("%d°%02.0f'%s", degrees, minutes, direction)
}

// FormatLongitude formats longitude for display with direction
func (l Location) FormatLongitude() string {
	direction := "E"
	lon := l.Longitude
	if lon < 0 {
		direction = "W"
		lon = -lon
	}

	degrees := int(lon)
	minutes := (lon - float64(degrees)) * 60

	return fmt.Sprintf("%d°%02.0f'%s", degrees, minutes, direction)
}

// FormatCoordinates formats both latitude and longitude
func (l Location) FormatCoordinates() string {
	return fmt.Sprintf("%s, %s", l.FormatLatitude(), l.FormatLongitude())
}

// DistanceTo calculates the distance to another location in kilometers using Haversine formula
func (l Location) DistanceTo(other Location) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1 := l.Latitude * math.Pi / 180
	lon1 := l.Longitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	lon2 := other.Longitude * math.Pi / 180

	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// BearingTo calculates the initial bearing from this location to another
func (l Location) BearingTo(other Location) float64 {
	// Convert degrees to radians
	lat1 := l.Latitude * math.Pi / 180
	lon1 := l.Longitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	lon2 := other.Longitude * math.Pi / 180

	dlon := lon2 - lon1

	y := math.Sin(dlon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)

	bearing := math.Atan2(y, x)

	// Convert to degrees and normalize to 0-360
	bearing = bearing * 180 / math.Pi
	return normalizeAngle(bearing)
}

// IsNorthernHemisphere returns true if the location is in the northern hemisphere
func (l Location) IsNorthernHemisphere() bool {
	return l.Latitude > 0
}

// IsSouthernHemisphere returns true if the location is in the southern hemisphere
func (l Location) IsSouthernHemisphere() bool {
	return l.Latitude < 0
}

// IsEasternHemisphere returns true if the location is in the eastern hemisphere
func (l Location) IsEasternHemisphere() bool {
	return l.Longitude > 0
}

// IsWesternHemisphere returns true if the location is in the western hemisphere
func (l Location) IsWesternHemisphere() bool {
	return l.Longitude < 0
}

// GetQuadrant returns the quadrant of the location (NE, NW, SE, SW)
func (l Location) GetQuadrant() string {
	var ns, ew string

	if l.IsNorthernHemisphere() {
		ns = "N"
	} else {
		ns = "S"
	}

	if l.IsEasternHemisphere() {
		ew = "E"
	} else {
		ew = "W"
	}

	return ns + ew
}

// GetTimeZoneOffset attempts to parse timezone offset from timezone name
func (l Location) GetTimeZoneOffset() float64 {
	// This is a simplified approach. In a full implementation,
	// you would use proper timezone libraries to get accurate offsets
	// accounting for daylight saving time, etc.

	timezoneOffsets := map[string]float64{
		"America/New_York":    -5,
		"America/Chicago":     -6,
		"America/Denver":      -7,
		"America/Los_Angeles": -8,
		"Europe/London":       0,
		"Europe/Paris":        1,
		"Europe/Berlin":       1,
		"Asia/Tokyo":          9,
		"Asia/Shanghai":       8,
		"Australia/Sydney":    10,
		// Add more as needed
	}

	if offset, exists := timezoneOffsets[l.Timezone]; exists {
		return offset
	}

	return 0 // Default to GMT
}

// IsWithinDistance checks if another location is within the specified distance (km)
func (l Location) IsWithinDistance(other Location, maxDistance float64) bool {
	return l.DistanceTo(other) <= maxDistance
}

// GetCardinalDirection returns the cardinal direction to another location
func (l Location) GetCardinalDirection(other Location) string {
	bearing := l.BearingTo(other)

	directions := []string{
		"N", "NNE", "NE", "ENE",
		"E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}

	index := int((bearing+11.25)/22.5) % 16
	return directions[index]
}

// Validate checks if the location data is valid
func (l Location) Validate() error {
	if l.City == "" {
		return fmt.Errorf("city name is required")
	}

	if l.Country == "" {
		return fmt.Errorf("country is required")
	}

	if !l.IsValidCoordinates() {
		return fmt.Errorf("invalid coordinates: latitude must be between -90 and 90, longitude between -180 and 180")
	}

	if l.Timezone == "" {
		return fmt.Errorf("timezone is required")
	}

	return nil
}

// String returns a string representation of the location
func (l Location) String() string {
	if l.Region != "" {
		return fmt.Sprintf("%s, %s, %s (%s)", l.City, l.Region, l.Country, l.FormatCoordinates())
	}
	return fmt.Sprintf("%s, %s (%s)", l.City, l.Country, l.FormatCoordinates())
}

// Equals checks if two locations are approximately equal
func (l Location) Equals(other Location) bool {
	const tolerance = 0.001 // ~100 meters

	return math.Abs(l.Latitude-other.Latitude) < tolerance &&
		math.Abs(l.Longitude-other.Longitude) < tolerance
}

// IsCity checks if this location represents a city (as opposed to a region or country)
func (l Location) IsCity() bool {
	return l.City != ""
}

// GetDisplayName returns the most appropriate display name for the location
func (l Location) GetDisplayName() string {
	if l.Name != "" {
		return l.Name
	}
	if l.City != "" {
		return l.City
	}
	return l.String()
}
