package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeInfo represents time-related information for astrological calculations
type TimeInfo struct {
	LocalTime    time.Time `json:"local_time"`
	UTCTime      time.Time `json:"utc_time"`
	JulianDay    float64   `json:"julian_day"`
	Timezone     string    `json:"timezone"`
	GMTOffset    float64   `json:"gmt_offset"` // Offset from GMT in hours
	DayOfYear    int       `json:"day_of_year"`
	SiderealTime float64   `json:"sidereal_time"` // Local sidereal time in hours
}

// ParseTime parses a date and time string into TimeInfo
func ParseTime(year, month, day int, timeStr, timezone string) (*TimeInfo, error) {
	// Parse the time string (HH:MM:SS)
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		// Try HH:MM format
		parsedTime, err = time.Parse("15:04", timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time format: %s", timeStr)
		}
	}

	// Load the timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", timezone)
	}

	// Create the local time
	localTime := time.Date(year, time.Month(month), day,
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc)

	// Convert to UTC
	utcTime := localTime.UTC()

	// Calculate Julian Day
	julianDay := CalculateJulianDay(utcTime)

	// Calculate GMT offset
	_, offset := localTime.Zone()
	gmtOffset := float64(offset) / 3600.0 // Convert seconds to hours

	timeInfo := &TimeInfo{
		LocalTime:    localTime,
		UTCTime:      utcTime,
		JulianDay:    julianDay,
		Timezone:     timezone,
		GMTOffset:    gmtOffset,
		DayOfYear:    utcTime.YearDay(),
		SiderealTime: CalculateLocalSiderealTime(julianDay, 0), // Will be updated with longitude
	}

	return timeInfo, nil
}

// CalculateJulianDay calculates the Julian Day Number for a given UTC time
func CalculateJulianDay(utcTime time.Time) float64 {
	year := utcTime.Year()
	month := int(utcTime.Month())
	day := utcTime.Day()
	hour := utcTime.Hour()
	minute := utcTime.Minute()
	second := utcTime.Second()

	// Convert time to decimal hours
	decimalHours := float64(hour) + float64(minute)/60.0 + float64(second)/3600.0

	// Adjust for January and February
	if month <= 2 {
		year--
		month += 12
	}

	// Calculate Julian Day Number
	a := year / 100
	b := 2 - a + a/4

	// Julian Day calculation
	jd := int(365.25*float64(year+4716)) + int(30.6001*float64(month+1)) + day + b - 1524

	// Add the time fraction
	return float64(jd) + (decimalHours-12.0)/24.0
}

// CalculateLocalSiderealTime calculates Local Sidereal Time
func CalculateLocalSiderealTime(julianDay, longitude float64) float64 {
	// Days since J2000.0
	t := (julianDay - 2451545.0) / 36525.0

	// Greenwich Mean Sidereal Time at 0h UT
	gmst := 280.46061837 + 360.98564736629*(julianDay-2451545.0) + 0.000387933*t*t - t*t*t/38710000.0

	// Normalize to 0-360 degrees
	gmst = normalizeAngle(gmst)

	// Convert to hours (divide by 15)
	gmstHours := gmst / 15.0

	// Calculate Local Sidereal Time
	lst := gmstHours + longitude/15.0 // longitude in degrees, convert to hours

	// Normalize to 0-24 hours
	for lst < 0 {
		lst += 24
	}
	for lst >= 24 {
		lst -= 24
	}

	return lst
}

// FormatTimeForDisplay formats time for human-readable display
func (ti TimeInfo) FormatTimeForDisplay() string {
	return ti.LocalTime.Format("2006-01-02 15:04:05 MST")
}

// FormatDateForDisplay formats just the date
func (ti TimeInfo) FormatDateForDisplay() string {
	return ti.LocalTime.Format("2006-01-02")
}

// FormatTimeOnly formats just the time
func (ti TimeInfo) FormatTimeOnly() string {
	return ti.LocalTime.Format("15:04:05")
}

// GetSeason returns the astronomical season for the date (Northern Hemisphere)
func (ti TimeInfo) GetSeason() string {
	month := ti.LocalTime.Month()
	day := ti.LocalTime.Day()

	switch {
	case (month == 3 && day >= 20) || month == 4 || month == 5 || (month == 6 && day < 21):
		return "Spring"
	case (month == 6 && day >= 21) || month == 7 || month == 8 || (month == 9 && day < 23):
		return "Summer"
	case (month == 9 && day >= 23) || month == 10 || month == 11 || (month == 12 && day < 21):
		return "Autumn"
	default:
		return "Winter"
	}
}

// IsLeapYear returns true if the year is a leap year
func (ti TimeInfo) IsLeapYear() bool {
	year := ti.LocalTime.Year()
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// DaysInMonth returns the number of days in the current month
func (ti TimeInfo) DaysInMonth() int {
	month := ti.LocalTime.Month()
	year := ti.LocalTime.Year()
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetMoonPhase returns an approximate moon phase (0-7, where 0=New, 4=Full)
func (ti TimeInfo) GetMoonPhase() int {
	// Simplified moon phase calculation
	// This is approximate and should be replaced with more accurate ephemeris data
	synodic := 29.53058868 // Average synodic month length in days

	// Known new moon date (approximate)
	knownNewMoon := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	daysSinceNewMoon := ti.UTCTime.Sub(knownNewMoon).Hours() / 24.0

	// Calculate current lunar cycle position
	cyclePosition := daysSinceNewMoon / synodic
	cyclePosition = cyclePosition - float64(int(cyclePosition)) // Get fractional part

	// Convert to phase (0-7)
	phase := int(cyclePosition * 8)
	if phase < 0 {
		phase += 8
	}

	return phase % 8
}

// GetMoonPhaseName returns the name of the current moon phase
func (ti TimeInfo) GetMoonPhaseName() string {
	phases := []string{
		"New Moon",
		"Waxing Crescent",
		"First Quarter",
		"Waxing Gibbous",
		"Full Moon",
		"Waning Gibbous",
		"Third Quarter",
		"Waning Crescent",
	}

	phase := ti.GetMoonPhase()
	return phases[phase]
}

// AddDays adds the specified number of days to the time
func (ti TimeInfo) AddDays(days int) *TimeInfo {
	newLocalTime := ti.LocalTime.AddDate(0, 0, days)
	newUTCTime := newLocalTime.UTC()
	newJulianDay := CalculateJulianDay(newUTCTime)

	return &TimeInfo{
		LocalTime:    newLocalTime,
		UTCTime:      newUTCTime,
		JulianDay:    newJulianDay,
		Timezone:     ti.Timezone,
		GMTOffset:    ti.GMTOffset,
		DayOfYear:    newUTCTime.YearDay(),
		SiderealTime: CalculateLocalSiderealTime(newJulianDay, 0),
	}
}

// AddYears adds the specified number of years to the time
func (ti TimeInfo) AddYears(years int) *TimeInfo {
	newLocalTime := ti.LocalTime.AddDate(years, 0, 0)
	newUTCTime := newLocalTime.UTC()
	newJulianDay := CalculateJulianDay(newUTCTime)

	return &TimeInfo{
		LocalTime:    newLocalTime,
		UTCTime:      newUTCTime,
		JulianDay:    newJulianDay,
		Timezone:     ti.Timezone,
		GMTOffset:    ti.GMTOffset,
		DayOfYear:    newUTCTime.YearDay(),
		SiderealTime: CalculateLocalSiderealTime(newJulianDay, 0),
	}
}

// ParseDateString parses various date string formats
func ParseDateString(dateStr string) (year, month, day int, err error) {
	// Try different formats
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"02/01/2006", // European format
		"2006/01/02",
		"01-02-2006",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, e := time.Parse(format, dateStr); e == nil {
			return t.Year(), int(t.Month()), t.Day(), nil
		}
	}

	// Try manual parsing for formats like "15 March 1990"
	parts := strings.Fields(dateStr)
	if len(parts) == 3 {
		if d, e1 := strconv.Atoi(parts[0]); e1 == nil {
			if y, e2 := strconv.Atoi(parts[2]); e2 == nil {
				if m := parseMonthName(parts[1]); m > 0 {
					return y, m, d, nil
				}
			}
		}
	}

	return 0, 0, 0, fmt.Errorf("unable to parse date string: %s", dateStr)
}

// parseMonthName converts month names to numbers
func parseMonthName(monthStr string) int {
	months := map[string]int{
		"january": 1, "jan": 1,
		"february": 2, "feb": 2,
		"march": 3, "mar": 3,
		"april": 4, "apr": 4,
		"may":  5,
		"june": 6, "jun": 6,
		"july": 7, "jul": 7,
		"august": 8, "aug": 8,
		"september": 9, "sep": 9,
		"october": 10, "oct": 10,
		"november": 11, "nov": 11,
		"december": 12, "dec": 12,
	}

	return months[strings.ToLower(monthStr)]
}
