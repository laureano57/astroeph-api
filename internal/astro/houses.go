package astro

import (
	"astroeph-api/internal/domain"
	"math"
)

// HouseCalculator handles house-related calculations
type HouseCalculator struct {
	ephemeris *Ephemeris
}

// NewHouseCalculator creates a new house calculator
func NewHouseCalculator(ephemeris *Ephemeris) *HouseCalculator {
	return &HouseCalculator{
		ephemeris: ephemeris,
	}
}

// CalculateHouses calculates houses for a chart
func (hc *HouseCalculator) CalculateHouses(
	timeInfo *domain.TimeInfo,
	location *domain.Location,
	houseSystem domain.HouseSystem,
) ([]domain.House, error) {

	// Convert to Julian Day
	julianDay := hc.ephemeris.GetJulianDay(timeInfo)

	// Get house system code
	systemCode := hc.ephemeris.GetHouseSystemCode(string(houseSystem))

	// Calculate houses using ephemeris
	housesData, err := hc.ephemeris.CalculateHouses(
		julianDay,
		location.Latitude,
		location.Longitude,
		systemCode,
	)
	if err != nil {
		return nil, err
	}

	// Convert to domain houses
	houses := make([]domain.House, 12)
	houseSizes := CalculateHouseSizes(housesData.Cusps)

	for i := 0; i < 12; i++ {
		house := domain.NewHouse(i+1, housesData.Cusps[i])
		house.Size = houseSizes[i]
		houses[i] = house
	}

	return houses, nil
}

// DetermineHouseForPlanet determines which house a planet is in
func (hc *HouseCalculator) DetermineHouseForPlanet(
	planetLongitude float64,
	houseCusps []float64,
) int {
	if len(houseCusps) != 12 {
		return 1 // Default to first house
	}

	// Normalize planet longitude to 0-360 range
	planetLon := normalizeAngle360(planetLongitude)

	// Check each house
	for i := 0; i < 12; i++ {
		currentCusp := normalizeAngle360(houseCusps[i])
		nextCusp := normalizeAngle360(houseCusps[(i+1)%12])

		// Handle the wrap-around case (house 12 to house 1)
		if currentCusp > nextCusp {
			// We cross 0° boundary
			if planetLon >= currentCusp || planetLon < nextCusp {
				return i + 1
			}
		} else {
			// Normal case - no 0° crossing
			if planetLon >= currentCusp && planetLon < nextCusp {
				return i + 1
			}
		}
	}

	// If we get here, find closest house
	return hc.findClosestHouse(planetLon, houseCusps)
}

// findClosestHouse finds the closest house when normal calculation fails
func (hc *HouseCalculator) findClosestHouse(planetLon float64, houseCusps []float64) int {
	minDiff := 360.0
	closestHouse := 1

	for i := 0; i < 12; i++ {
		cuspLon := normalizeAngle360(houseCusps[i])
		diff := math.Abs(planetLon - cuspLon)
		if diff > 180 {
			diff = 360 - diff // Handle wrap-around
		}
		if diff < minDiff {
			minDiff = diff
			closestHouse = i + 1
		}
	}

	return closestHouse
}

// CalculateHouseSizes calculates the size of each house given cusps
func CalculateHouseSizes(cusps []float64) []float64 {
	if len(cusps) != 12 {
		return make([]float64, 12) // Return zeros if invalid input
	}

	sizes := make([]float64, 12)
	for i := 0; i < 12; i++ {
		nextIndex := (i + 1) % 12
		size := cusps[nextIndex] - cusps[i]

		// Handle wrap-around at 360 degrees
		if size < 0 {
			size += 360
		}

		sizes[i] = size
	}

	return sizes
}

// GetHouseSystemName returns the full name for a house system
func GetHouseSystemName(system domain.HouseSystem) string {
	switch system {
	case domain.HousePlacidus:
		return "Placidus"
	case domain.HouseKoch:
		return "Koch"
	case domain.HousePorphyrius:
		return "Porphyrius"
	case domain.HouseRegiomontanus:
		return "Regiomontanus"
	case domain.HouseCampanus:
		return "Campanus"
	case domain.HouseEqual:
		return "Equal"
	case domain.HouseWholeSign:
		return "Whole Sign"
	default:
		return "Placidus"
	}
}

// GetAvailableHouseSystems returns all available house systems
func GetAvailableHouseSystems() []domain.HouseSystem {
	return []domain.HouseSystem{
		domain.HousePlacidus,
		domain.HouseKoch,
		domain.HousePorphyrius,
		domain.HouseRegiomontanus,
		domain.HouseCampanus,
		domain.HouseEqual,
		domain.HouseWholeSign,
	}
}

// IsValidHouseSystem checks if a house system name is valid
func IsValidHouseSystem(system string) bool {
	validSystems := map[string]bool{
		"Placidus":      true,
		"Koch":          true,
		"Porphyrius":    true,
		"Regiomontanus": true,
		"Campanus":      true,
		"Equal":         true,
		"Whole Sign":    true,
	}

	return validSystems[system]
}

// GetDefaultHouseSystem returns the default house system
func GetDefaultHouseSystem() domain.HouseSystem {
	return domain.HousePlacidus
}

// CalculateHouseRulers calculates the ruling planets for each house
func (hc *HouseCalculator) CalculateHouseRulers(houses []domain.House) map[int]string {
	rulers := make(map[int]string)

	for _, house := range houses {
		rulers[house.Number] = house.Ruler
	}

	return rulers
}

// GetHouseKeywords returns keywords for all houses
func GetHouseKeywords() map[int][]string {
	keywords := make(map[int][]string)

	infos := domain.GetHouseInfos()
	for _, info := range infos {
		keywords[info.Number] = info.Keywords
	}

	return keywords
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
