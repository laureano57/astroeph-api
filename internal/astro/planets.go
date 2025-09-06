package astro

import (
	"astroeph-api/internal/domain"
	"fmt"
	"math"
)

// PlanetCalculator handles planet-related calculations
type PlanetCalculator struct {
	ephemeris *Ephemeris
}

// NewPlanetCalculator creates a new planet calculator
func NewPlanetCalculator(ephemeris *Ephemeris) *PlanetCalculator {
	return &PlanetCalculator{
		ephemeris: ephemeris,
	}
}

// CalculateAllPlanets calculates positions for all planets
func (pc *PlanetCalculator) CalculateAllPlanets(
	timeInfo *domain.TimeInfo,
	houseCusps []float64,
) ([]domain.Planet, error) {

	// Convert to Julian Day
	julianDay := pc.ephemeris.GetJulianDay(timeInfo)

	// Calculate planet positions using ephemeris
	positions, err := pc.ephemeris.CalculateAllPlanets(julianDay)
	if err != nil {
		return nil, err
	}

	// Convert to domain planets
	var planets []domain.Planet
	houseCalc := NewHouseCalculator(pc.ephemeris)

	for _, pos := range positions {
		// Determine which house the planet is in
		houseNumber := houseCalc.DetermineHouseForPlanet(pos.Longitude, houseCusps)

		// Convert to domain planet
		planet := pos.ToDomainPlanet(pc.ephemeris, houseNumber)
		planets = append(planets, planet)
	}

	return planets, nil
}

// CalculateSinglePlanet calculates position for a single planet
func (pc *PlanetCalculator) CalculateSinglePlanet(
	planetName string,
	timeInfo *domain.TimeInfo,
	houseCusps []float64,
) (*domain.Planet, error) {

	// Get planet ID from name
	planetID := pc.getPlanetIDFromName(planetName)
	if planetID == -1 {
		return nil, fmt.Errorf("unknown planet: %s", planetName)
	}

	// Convert to Julian Day
	julianDay := pc.ephemeris.GetJulianDay(timeInfo)

	// Calculate planet position
	pos, err := pc.ephemeris.CalculatePlanetPosition(julianDay, planetID)
	if err != nil {
		return nil, err
	}

	// Determine house
	houseCalc := NewHouseCalculator(pc.ephemeris)
	houseNumber := houseCalc.DetermineHouseForPlanet(pos.Longitude, houseCusps)

	// Convert to domain planet
	planet := pos.ToDomainPlanet(pc.ephemeris, houseNumber)

	return &planet, nil
}

// CalculatePlanetaryDignities calculates dignities for all planets
func (pc *PlanetCalculator) CalculatePlanetaryDignities(planets []domain.Planet) map[string]PlanetaryDignity {
	dignities := make(map[string]PlanetaryDignity)

	for _, planet := range planets {
		dignity := pc.calculatePlanetDignity(planet)
		dignities[planet.Name] = dignity
	}

	return dignities
}

// calculatePlanetDignity calculates dignities for a single planet
func (pc *PlanetCalculator) calculatePlanetDignity(planet domain.Planet) PlanetaryDignity {
	dignity := PlanetaryDignity{
		Planet: planet.Name,
		Sign:   planet.Sign,
	}

	// Check various dignities
	dignity.IsInDomicile = pc.isInDomicile(planet)
	dignity.IsInExaltation = pc.isInExaltation(planet)
	dignity.IsInDetriment = pc.isInDetriment(planet)
	dignity.IsInFall = pc.isInFall(planet)
	dignity.IsPeregrine = !dignity.IsInDomicile && !dignity.IsInExaltation && !dignity.IsInDetriment && !dignity.IsInFall

	// Calculate dignity score (-2 to +2)
	if dignity.IsInExaltation {
		dignity.Score = 2
	} else if dignity.IsInDomicile {
		dignity.Score = 1
	} else if dignity.IsInDetriment {
		dignity.Score = -1
	} else if dignity.IsInFall {
		dignity.Score = -2
	} else {
		dignity.Score = 0 // Peregrine
	}

	return dignity
}

// CalculatePlanetaryReturns calculates when a planet returns to its natal position
func (pc *PlanetCalculator) CalculatePlanetaryReturns(
	natalPlanet domain.Planet,
	startDate *domain.TimeInfo,
) ([]PlanetaryReturn, error) {

	var returns []PlanetaryReturn

	// Get approximate orbital period in days
	orbitalPeriod := pc.getOrbitalPeriod(natalPlanet.Name)
	if orbitalPeriod == 0 {
		return returns, nil // No orbital period data
	}

	// Calculate approximate return dates
	currentDate := startDate
	for i := 1; i <= 10; i++ { // Calculate up to 10 returns
		returnDate := currentDate.AddDays(int(orbitalPeriod * float64(i)))

		returns = append(returns, PlanetaryReturn{
			Planet:     natalPlanet.Name,
			ReturnDate: *returnDate,
			ReturnType: getReturnType(natalPlanet.Name),
		})
	}

	return returns, nil
}

// CalculateProgressions calculates secondary progressions for planets
func (pc *PlanetCalculator) CalculateProgressions(
	natalPlanets []domain.Planet,
	natalTime *domain.TimeInfo,
	progressionDate *domain.TimeInfo,
) ([]domain.Planet, error) {

	// Secondary progressions: 1 day = 1 year
	yearsDiff := progressionDate.LocalTime.Sub(natalTime.LocalTime).Hours() / (24 * 365.25)
	daysDiff := yearsDiff // In secondary progressions, years become days

	// Create new time for progression calculation
	progressedTime := natalTime.AddDays(int(daysDiff))

	// Calculate planet positions for progressed time
	return pc.CalculateAllPlanets(progressedTime, []float64{})
}

// getPlanetIDFromName converts planet name to swephgo ID
func (pc *PlanetCalculator) getPlanetIDFromName(name string) int {
	planetIDs := map[string]int{
		"Sun":        SE_SUN,
		"Moon":       SE_MOON,
		"Mercury":    SE_MERCURY,
		"Venus":      SE_VENUS,
		"Mars":       SE_MARS,
		"Jupiter":    SE_JUPITER,
		"Saturn":     SE_SATURN,
		"Uranus":     SE_URANUS,
		"Neptune":    SE_NEPTUNE,
		"Pluto":      SE_PLUTO,
		"North Node": SE_MEAN_NODE,
		"Chiron":     SE_CHIRON,
	}

	if id, exists := planetIDs[name]; exists {
		return id
	}
	return -1
}

// isInDomicile checks if planet is in its domicile (ruling) sign
func (pc *PlanetCalculator) isInDomicile(planet domain.Planet) bool {
	domiciles := map[string][]string{
		"Sun":     {"Leo"},
		"Moon":    {"Cancer"},
		"Mercury": {"Gemini", "Virgo"},
		"Venus":   {"Taurus", "Libra"},
		"Mars":    {"Aries", "Scorpio"},      // Traditional rulership
		"Jupiter": {"Sagittarius", "Pisces"}, // Traditional rulership
		"Saturn":  {"Capricorn", "Aquarius"}, // Traditional rulership
		"Uranus":  {"Aquarius"},
		"Neptune": {"Pisces"},
		"Pluto":   {"Scorpio"},
	}

	if signs, exists := domiciles[planet.Name]; exists {
		for _, sign := range signs {
			if planet.Sign == sign {
				return true
			}
		}
	}
	return false
}

// isInExaltation checks if planet is in its exaltation sign
func (pc *PlanetCalculator) isInExaltation(planet domain.Planet) bool {
	exaltations := map[string]string{
		"Sun":     "Aries",
		"Moon":    "Taurus",
		"Mercury": "Virgo",
		"Venus":   "Pisces",
		"Mars":    "Capricorn",
		"Jupiter": "Cancer",
		"Saturn":  "Libra",
		"Uranus":  "Scorpio",
		"Neptune": "Aquarius",
		"Pluto":   "Aries",
	}

	if sign, exists := exaltations[planet.Name]; exists {
		return planet.Sign == sign
	}
	return false
}

// isInDetriment checks if planet is in its detriment sign
func (pc *PlanetCalculator) isInDetriment(planet domain.Planet) bool {
	return planet.IsInDetriment()
}

// isInFall checks if planet is in its fall sign
func (pc *PlanetCalculator) isInFall(planet domain.Planet) bool {
	return planet.IsInFall()
}

// getOrbitalPeriod returns the orbital period of a planet in days
func (pc *PlanetCalculator) getOrbitalPeriod(planetName string) float64 {
	orbitalPeriods := map[string]float64{
		"Sun":        365.25,   // Solar return
		"Moon":       27.32,    // Lunar return (monthly)
		"Mercury":    87.97,    // ~88 days
		"Venus":      224.70,   // ~225 days
		"Mars":       686.98,   // ~687 days
		"Jupiter":    4332.59,  // ~12 years
		"Saturn":     10759.22, // ~29.5 years
		"Uranus":     30688.5,  // ~84 years
		"Neptune":    60182,    // ~165 years
		"Pluto":      90560,    // ~248 years
		"North Node": -6798.38, // ~18.6 years (retrograde)
	}

	if period, exists := orbitalPeriods[planetName]; exists {
		return math.Abs(period) // Return absolute value
	}
	return 0
}

// getReturnType determines the type of planetary return
func getReturnType(planetName string) string {
	returnTypes := map[string]string{
		"Sun":        "Solar Return",
		"Moon":       "Lunar Return",
		"Mercury":    "Mercury Return",
		"Venus":      "Venus Return",
		"Mars":       "Mars Return",
		"Jupiter":    "Jupiter Return",
		"Saturn":     "Saturn Return",
		"Uranus":     "Uranus Return",
		"Neptune":    "Neptune Return",
		"Pluto":      "Pluto Return",
		"North Node": "Nodal Return",
	}

	if returnType, exists := returnTypes[planetName]; exists {
		return returnType
	}
	return "Planetary Return"
}

// PlanetaryDignity represents the dignity of a planet in a sign
type PlanetaryDignity struct {
	Planet         string `json:"planet"`
	Sign           string `json:"sign"`
	IsInDomicile   bool   `json:"is_in_domicile"`
	IsInExaltation bool   `json:"is_in_exaltation"`
	IsInDetriment  bool   `json:"is_in_detriment"`
	IsInFall       bool   `json:"is_in_fall"`
	IsPeregrine    bool   `json:"is_peregrine"`
	Score          int    `json:"score"` // -2 to +2
}

// PlanetaryReturn represents a planetary return
type PlanetaryReturn struct {
	Planet     string          `json:"planet"`
	ReturnDate domain.TimeInfo `json:"return_date"`
	ReturnType string          `json:"return_type"`
}

// GetPlanetaryStrengths calculates overall strength scores for planets
func (pc *PlanetCalculator) GetPlanetaryStrengths(
	planets []domain.Planet,
	dignities map[string]PlanetaryDignity,
) map[string]PlanetaryStrength {

	strengths := make(map[string]PlanetaryStrength)

	for _, planet := range planets {
		dignity := dignities[planet.Name]

		strength := PlanetaryStrength{
			Planet:       planet.Name,
			DignityScore: dignity.Score,
			HouseScore:   pc.getHouseScore(planet.House),
			AspectScore:  0, // Would be calculated from aspects
			OverallScore: 0,
		}

		// Calculate overall score (this is simplified)
		strength.OverallScore = strength.DignityScore + strength.HouseScore + strength.AspectScore

		strengths[planet.Name] = strength
	}

	return strengths
}

// getHouseScore returns a score for house position
func (pc *PlanetCalculator) getHouseScore(houseNumber int) int {
	// Angular houses are stronger
	if houseNumber == 1 || houseNumber == 4 || houseNumber == 7 || houseNumber == 10 {
		return 2
	}
	// Succedent houses are moderate
	if houseNumber == 2 || houseNumber == 5 || houseNumber == 8 || houseNumber == 11 {
		return 1
	}
	// Cadent houses are weaker
	return 0
}

// PlanetaryStrength represents the overall strength of a planet
type PlanetaryStrength struct {
	Planet       string `json:"planet"`
	DignityScore int    `json:"dignity_score"`
	HouseScore   int    `json:"house_score"`
	AspectScore  int    `json:"aspect_score"`
	OverallScore int    `json:"overall_score"`
}

// Helper function to create a domain error (assuming this exists)
func NewError(code, message string) error {
	return &DomainError{Code: code, Message: message}
}

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}
