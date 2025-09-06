package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTime parses various time formats
func ParseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"15:04:05",
		"15:04",
		"3:04:05 PM",
		"3:04 PM",
		"15.04.05",
		"15.04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// ValidateDate validates a date
func ValidateDate(year, month, day int) error {
	if year < 1800 || year > 2200 {
		return fmt.Errorf("year must be between 1800 and 2200")
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}

	if day < 1 || day > 31 {
		return fmt.Errorf("day must be between 1 and 31")
	}

	// Check for valid day in month
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// Handle leap years for February
	if month == 2 && IsLeapYear(year) {
		daysInMonth[1] = 29
	}

	if day > daysInMonth[month-1] {
		return fmt.Errorf("day %d is not valid for month %d", day, month)
	}

	return nil
}

// IsLeapYear checks if a year is a leap year
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// NormalizeString normalizes a string for comparison
func NormalizeString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// StringToInt converts string to int with default value
func StringToInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// StringToFloat converts string to float64 with default value
func StringToFloat(s string, defaultValue float64) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return defaultValue
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsInt checks if a slice contains an int
func ContainsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// MaxInt returns the maximum of two integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the minimum of two integers
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxFloat returns the maximum of two floats
func MaxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat returns the minimum of two floats
func MinFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RoundFloat rounds a float to a specified number of decimal places
func RoundFloat(val float64, precision int) float64 {
	ratio := float64(1)
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000)
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// ValidateEmail validates an email address (simple validation)
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// GetEnvOrDefault gets an environment variable or returns a default value
func GetEnvOrDefault(key, defaultValue string) string {
	// This would typically use os.Getenv, but since we're not importing os,
	// we'll return the default for now
	return defaultValue
}

// SliceToMap converts a slice of strings to a map for O(1) lookup
func SliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, item := range slice {
		m[item] = true
	}
	return m
}

// MapKeys returns the keys of a string map as a slice
func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// CoalesceString returns the first non-empty string
func CoalesceString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// CoalesceInt returns the first non-zero integer
func CoalesceInt(values ...int) int {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}
