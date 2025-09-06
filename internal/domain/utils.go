package domain

import (
	"fmt"
	"math"
)

// normalizeAngle normalizes an angle to the range [0, 360) degrees
func normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// GetZodiacSign returns the zodiac sign for a given longitude
func GetZodiacSign(longitude float64) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}

	normalizedLon := normalizeAngle(longitude)
	signIndex := int(normalizedLon / 30.0)
	if signIndex >= 12 {
		signIndex = 11
	}
	return signs[signIndex]
}

// GetDegreeInSign returns the degree within a zodiac sign (0-29.999...)
func GetDegreeInSign(longitude float64) float64 {
	normalizedLon := normalizeAngle(longitude)
	return math.Mod(normalizedLon, 30.0)
}

// FormatDegreeInSign converts decimal degrees within a sign to DMS format
func FormatDegreeInSign(longitude float64) string {
	degreeInSign := GetDegreeInSign(longitude)
	return formatDegreesMinutesSeconds(degreeInSign)
}

// FormatLongitude converts a full longitude to DMS format
func FormatLongitude(longitude float64) string {
	return formatDegreesMinutesSeconds(longitude)
}

// formatDegreesMinutesSeconds converts decimal degrees to degrees, minutes and seconds format
func formatDegreesMinutesSeconds(decimalDegrees float64) string {
	degrees := int(decimalDegrees)
	remainingMinutes := (decimalDegrees - float64(degrees)) * 60
	minutes := int(remainingMinutes)
	seconds := int((remainingMinutes - float64(minutes)) * 60)
	return fmt.Sprintf("%d°%02d'%02d\"", degrees, minutes, seconds)
}

// AngularDistance calculates the shortest angular distance between two longitudes
func AngularDistance(lon1, lon2 float64) float64 {
	diff := math.Abs(normalizeAngle(lon1) - normalizeAngle(lon2))
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// GetElementForSign returns the element for a zodiac sign
func GetElementForSign(sign string) string {
	fireElements := map[string]bool{
		"Aries": true, "Leo": true, "Sagittarius": true,
	}
	earthElements := map[string]bool{
		"Taurus": true, "Virgo": true, "Capricorn": true,
	}
	airElements := map[string]bool{
		"Gemini": true, "Libra": true, "Aquarius": true,
	}

	if fireElements[sign] {
		return "fire"
	}
	if earthElements[sign] {
		return "earth"
	}
	if airElements[sign] {
		return "air"
	}
	return "water" // Cancer, Scorpio, Pisces
}

// GetModalityForSign returns the modality for a zodiac sign
func GetModalityForSign(sign string) string {
	cardinalSigns := map[string]bool{
		"Aries": true, "Cancer": true, "Libra": true, "Capricorn": true,
	}
	fixedSigns := map[string]bool{
		"Taurus": true, "Leo": true, "Scorpio": true, "Aquarius": true,
	}

	if cardinalSigns[sign] {
		return "cardinal"
	}
	if fixedSigns[sign] {
		return "fixed"
	}
	return "mutable" // Gemini, Virgo, Sagittarius, Pisces
}

// GetRulerForSign returns the traditional ruler for a zodiac sign
func GetRulerForSign(sign string) string {
	rulers := map[string]string{
		"Aries":       "Mars",
		"Taurus":      "Venus",
		"Gemini":      "Mercury",
		"Cancer":      "Moon",
		"Leo":         "Sun",
		"Virgo":       "Mercury",
		"Libra":       "Venus",
		"Scorpio":     "Mars", // Traditional ruler (modern: Pluto)
		"Sagittarius": "Jupiter",
		"Capricorn":   "Saturn",
		"Aquarius":    "Saturn",  // Traditional ruler (modern: Uranus)
		"Pisces":      "Jupiter", // Traditional ruler (modern: Neptune)
	}

	if ruler, exists := rulers[sign]; exists {
		return ruler
	}
	return ""
}

// GetModernRulerForSign returns the modern ruler for a zodiac sign
func GetModernRulerForSign(sign string) string {
	rulers := map[string]string{
		"Aries":       "Mars",
		"Taurus":      "Venus",
		"Gemini":      "Mercury",
		"Cancer":      "Moon",
		"Leo":         "Sun",
		"Virgo":       "Mercury",
		"Libra":       "Venus",
		"Scorpio":     "Pluto",
		"Sagittarius": "Jupiter",
		"Capricorn":   "Saturn",
		"Aquarius":    "Uranus",
		"Pisces":      "Neptune",
	}

	if ruler, exists := rulers[sign]; exists {
		return ruler
	}
	return ""
}

// GetPolarityForSign returns the polarity (positive/negative) for a zodiac sign
func GetPolarityForSign(sign string) string {
	positiveSigns := map[string]bool{
		"Aries": true, "Gemini": true, "Leo": true,
		"Libra": true, "Sagittarius": true, "Aquarius": true,
	}

	if positiveSigns[sign] {
		return "positive"
	}
	return "negative"
}

// IsFireSign returns true if the sign is a fire sign
func IsFireSign(sign string) bool {
	return GetElementForSign(sign) == "fire"
}

// IsEarthSign returns true if the sign is an earth sign
func IsEarthSign(sign string) bool {
	return GetElementForSign(sign) == "earth"
}

// IsAirSign returns true if the sign is an air sign
func IsAirSign(sign string) bool {
	return GetElementForSign(sign) == "air"
}

// IsWaterSign returns true if the sign is a water sign
func IsWaterSign(sign string) bool {
	return GetElementForSign(sign) == "water"
}

// IsCardinalSign returns true if the sign is cardinal
func IsCardinalSign(sign string) bool {
	return GetModalityForSign(sign) == "cardinal"
}

// IsFixedSign returns true if the sign is fixed
func IsFixedSign(sign string) bool {
	return GetModalityForSign(sign) == "fixed"
}

// IsMutableSign returns true if the sign is mutable
func IsMutableSign(sign string) bool {
	return GetModalityForSign(sign) == "mutable"
}

// IsPositiveSign returns true if the sign has positive polarity
func IsPositiveSign(sign string) bool {
	return GetPolarityForSign(sign) == "positive"
}

// IsNegativeSign returns true if the sign has negative polarity
func IsNegativeSign(sign string) bool {
	return GetPolarityForSign(sign) == "negative"
}

// GetSignNumber returns the sign number (1-12 for Aries-Pisces)
func GetSignNumber(sign string) int {
	signs := map[string]int{
		"Aries": 1, "Taurus": 2, "Gemini": 3, "Cancer": 4,
		"Leo": 5, "Virgo": 6, "Libra": 7, "Scorpio": 8,
		"Sagittarius": 9, "Capricorn": 10, "Aquarius": 11, "Pisces": 12,
	}

	if num, exists := signs[sign]; exists {
		return num
	}
	return 0
}

// GetSignByNumber returns the sign name for a given number (1-12)
func GetSignByNumber(number int) string {
	if number < 1 || number > 12 {
		return ""
	}

	signs := []string{
		"", "Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}

	return signs[number]
}

// GetOppositeSign returns the opposite sign
func GetOppositeSign(sign string) string {
	opposites := map[string]string{
		"Aries":       "Libra",
		"Taurus":      "Scorpio",
		"Gemini":      "Sagittarius",
		"Cancer":      "Capricorn",
		"Leo":         "Aquarius",
		"Virgo":       "Pisces",
		"Libra":       "Aries",
		"Scorpio":     "Taurus",
		"Sagittarius": "Gemini",
		"Capricorn":   "Cancer",
		"Aquarius":    "Leo",
		"Pisces":      "Virgo",
	}

	if opposite, exists := opposites[sign]; exists {
		return opposite
	}
	return ""
}

// AreCompatibleSigns checks if two signs are traditionally compatible
func AreCompatibleSigns(sign1, sign2 string) bool {
	// Same element signs are generally compatible
	if GetElementForSign(sign1) == GetElementForSign(sign2) {
		return true
	}

	// Fire and Air are compatible
	if (IsFireSign(sign1) && IsAirSign(sign2)) || (IsAirSign(sign1) && IsFireSign(sign2)) {
		return true
	}

	// Earth and Water are compatible
	if (IsEarthSign(sign1) && IsWaterSign(sign2)) || (IsWaterSign(sign1) && IsEarthSign(sign2)) {
		return true
	}

	return false
}

// GetSignKeywords returns keywords associated with a zodiac sign
func GetSignKeywords(sign string) []string {
	keywords := map[string][]string{
		"Aries":       {"pioneer", "leader", "energetic", "impulsive", "brave"},
		"Taurus":      {"stable", "practical", "sensual", "stubborn", "reliable"},
		"Gemini":      {"curious", "versatile", "communicative", "restless", "adaptable"},
		"Cancer":      {"nurturing", "emotional", "intuitive", "protective", "sensitive"},
		"Leo":         {"creative", "dramatic", "generous", "proud", "confident"},
		"Virgo":       {"analytical", "practical", "perfectionist", "helpful", "modest"},
		"Libra":       {"balanced", "harmonious", "diplomatic", "indecisive", "charming"},
		"Scorpio":     {"intense", "passionate", "mysterious", "transformative", "powerful"},
		"Sagittarius": {"adventurous", "philosophical", "optimistic", "freedom-loving", "honest"},
		"Capricorn":   {"ambitious", "disciplined", "practical", "responsible", "traditional"},
		"Aquarius":    {"innovative", "humanitarian", "independent", "eccentric", "intellectual"},
		"Pisces":      {"intuitive", "compassionate", "artistic", "dreamy", "spiritual"},
	}

	if words, exists := keywords[sign]; exists {
		return words
	}
	return []string{}
}
