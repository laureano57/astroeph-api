package astro

import (
	"astroeph-api/internal/domain"
	"astroeph-api/internal/logging"
	"fmt"
	"math"

	"github.com/mshafiee/swephgo"
)

// Ephemeris provides a wrapper around Swiss Ephemeris (swephgo)
type Ephemeris struct {
	logger      *logging.Logger
	initialized bool
}

// Planet constants for swephgo
const (
	SE_SUN       = 0
	SE_MOON      = 1
	SE_MERCURY   = 2
	SE_VENUS     = 3
	SE_MARS      = 4
	SE_JUPITER   = 5
	SE_SATURN    = 6
	SE_URANUS    = 7
	SE_NEPTUNE   = 8
	SE_PLUTO     = 9
	SE_MEAN_NODE = 10
	SE_TRUE_NODE = 11
	SE_CHIRON    = 15
)

// NewEphemeris creates a new Ephemeris instance
func NewEphemeris(logger *logging.Logger) (*Ephemeris, error) {
	eph := &Ephemeris{
		logger: logger,
	}

	if err := eph.initialize(); err != nil {
		return nil, err
	}

	return eph, nil
}

// initialize initializes the Swiss Ephemeris
func (e *Ephemeris) initialize() error {
	// Set ephemeris path - in production, you might want to set a specific path
	swephgo.SetEphePath([]byte("")) // Use built-in ephemeris data

	e.logger.Info().Msg("🔮 Initializing Swiss Ephemeris")

	// Test Swiss Ephemeris initialization
	testJD := swephgo.Julday(2000, 1, 1, 12.0, 1)
	xx := make([]float64, 6)
	serr := make([]byte, 256)
	result := swephgo.Calc(testJD, 0, 0, xx, serr)

	if result < 0 {
		e.logger.Error().
			Int("result_code", int(result)).
			Str("error", string(serr)).
			Msg("Swiss Ephemeris test calculation failed")
		return fmt.Errorf("Swiss Ephemeris initialization failed: %s", string(serr))
	}

	e.logger.Info().
		Float64("test_sun_longitude", xx[0]).
		Str("ephemeris_status", string(serr)).
		Msg("✅ Swiss Ephemeris initialized successfully")

	e.initialized = true
	return nil
}

// CalculatePlanetPosition calculates the position of a planet for a given Julian Day
func (e *Ephemeris) CalculatePlanetPosition(julianDay float64, planetID int) (*PlanetPosition, error) {
	if !e.initialized {
		return nil, fmt.Errorf("ephemeris not initialized")
	}

	xx := make([]float64, 6)
	serr := make([]byte, 256)
	result := swephgo.Calc(julianDay, planetID, 0, xx, serr)

	if result < 0 {
		return nil, fmt.Errorf("failed to calculate position for planet %d: %s", planetID, string(serr))
	}

	pos := &PlanetPosition{
		PlanetID:  planetID,
		Longitude: xx[0],
		Latitude:  xx[1],
		Distance:  xx[2],
		LongSpeed: xx[3],
		LatSpeed:  xx[4],
		DistSpeed: xx[5],
	}

	return pos, nil
}

// CalculateAllPlanets calculates positions for all main planets
func (e *Ephemeris) CalculateAllPlanets(julianDay float64) ([]PlanetPosition, error) {
	if !e.initialized {
		return nil, fmt.Errorf("ephemeris not initialized")
	}

	mainPlanets := []int{
		SE_SUN, SE_MOON, SE_MERCURY, SE_VENUS, SE_MARS,
		SE_JUPITER, SE_SATURN, SE_URANUS, SE_NEPTUNE, SE_PLUTO,
		SE_MEAN_NODE, SE_CHIRON,
	}

	var positions []PlanetPosition

	for _, planetID := range mainPlanets {
		pos, err := e.CalculatePlanetPosition(julianDay, planetID)
		if err != nil {
			e.logger.Warn().
				Err(err).
				Int("planet_id", planetID).
				Msg("Failed to calculate planet position, skipping")
			continue
		}
		positions = append(positions, *pos)
	}

	return positions, nil
}

// CalculateHouses calculates house cusps using Swiss Ephemeris
func (e *Ephemeris) CalculateHouses(julianDay, latitude, longitude float64, houseSystem rune) (*HousesData, error) {
	if !e.initialized {
		return nil, fmt.Errorf("ephemeris not initialized")
	}

	// Calculate houses using swephgo
	cusps := make([]float64, 13) // 0-12, where 1-12 are the house cusps
	ascmc := make([]float64, 10) // Ascendant, MC, etc.
	result := swephgo.Houses(julianDay, latitude, longitude, int(houseSystem), cusps, ascmc)

	if result < 0 {
		return nil, fmt.Errorf("failed to calculate houses: house system not supported or invalid parameters")
	}

	housesData := &HousesData{
		Cusps:         cusps[1:13], // Houses 1-12
		Ascendant:     cusps[1],    // 1st house cusp is the Ascendant
		Midheaven:     cusps[10],   // 10th house cusp is the Midheaven
		IC:            cusps[4],    // 4th house cusp is the IC
		Descendant:    cusps[7],    // 7th house cusp is the Descendant
		ARMC:          ascmc[2],    // Right Ascension of MC
		Vertex:        ascmc[3],    // Vertex
		EquatorialAsc: ascmc[4],    // Equatorial Ascendant
		CoAscendant1:  ascmc[5],    // Co-ascendant (Koch)
		CoAscendant2:  ascmc[6],    // Co-ascendant (Munkasey)
		PolarAsc:      ascmc[7],    // Polar ascendant
	}

	return housesData, nil
}

// GetJulianDay converts a date/time to Julian Day Number
func (e *Ephemeris) GetJulianDay(timeInfo *domain.TimeInfo) float64 {
	utc := timeInfo.UTCTime
	hour := float64(utc.Hour()) + float64(utc.Minute())/60.0 + float64(utc.Second())/3600.0
	return swephgo.Julday(utc.Year(), int(utc.Month()), utc.Day(), hour, 1)
}

// GetPlanetName returns the name of a planet by its ID
func (e *Ephemeris) GetPlanetName(planetID int) string {
	planetNames := map[int]string{
		SE_SUN:       "Sun",
		SE_MOON:      "Moon",
		SE_MERCURY:   "Mercury",
		SE_VENUS:     "Venus",
		SE_MARS:      "Mars",
		SE_JUPITER:   "Jupiter",
		SE_SATURN:    "Saturn",
		SE_URANUS:    "Uranus",
		SE_NEPTUNE:   "Neptune",
		SE_PLUTO:     "Pluto",
		SE_MEAN_NODE: "North Node",
		SE_CHIRON:    "Chiron",
	}

	if name, exists := planetNames[planetID]; exists {
		return name
	}
	return fmt.Sprintf("Planet_%d", planetID)
}

// GetHouseSystemCode converts house system name to swephgo code
func (e *Ephemeris) GetHouseSystemCode(system string) rune {
	const (
		SE_HOUSE_PLACIDUS      = 'P'
		SE_HOUSE_KOCH          = 'K'
		SE_HOUSE_PORPHYRIUS    = 'O'
		SE_HOUSE_REGIOMONTANUS = 'R'
		SE_HOUSE_CAMPANUS      = 'C'
		SE_HOUSE_EQUAL         = 'E'
		SE_HOUSE_WHOLE_SIGN    = 'W'
	)

	switch system {
	case "Koch":
		return SE_HOUSE_KOCH
	case "Porphyrius":
		return SE_HOUSE_PORPHYRIUS
	case "Regiomontanus":
		return SE_HOUSE_REGIOMONTANUS
	case "Campanus":
		return SE_HOUSE_CAMPANUS
	case "Equal":
		return SE_HOUSE_EQUAL
	case "Whole Sign":
		return SE_HOUSE_WHOLE_SIGN
	default:
		return SE_HOUSE_PLACIDUS // Default to Placidus
	}
}

// PlanetPosition holds calculated planet position data
type PlanetPosition struct {
	PlanetID  int     `json:"planet_id"`
	Longitude float64 `json:"longitude"`  // Longitude in degrees
	Latitude  float64 `json:"latitude"`   // Latitude in degrees
	Distance  float64 `json:"distance"`   // Distance from Earth in AU
	LongSpeed float64 `json:"long_speed"` // Longitude speed in degrees/day
	LatSpeed  float64 `json:"lat_speed"`  // Latitude speed in degrees/day
	DistSpeed float64 `json:"dist_speed"` // Distance speed in AU/day
}

// HousesData holds calculated house data
type HousesData struct {
	Cusps         []float64 `json:"cusps"`          // House cusps 1-12
	Ascendant     float64   `json:"ascendant"`      // Ascendant (1st house cusp)
	Midheaven     float64   `json:"midheaven"`      // Midheaven (10th house cusp)
	IC            float64   `json:"ic"`             // IC (4th house cusp)
	Descendant    float64   `json:"descendant"`     // Descendant (7th house cusp)
	ARMC          float64   `json:"armc"`           // Right Ascension of MC
	Vertex        float64   `json:"vertex"`         // Vertex
	EquatorialAsc float64   `json:"equatorial_asc"` // Equatorial Ascendant
	CoAscendant1  float64   `json:"co_ascendant1"`  // Co-ascendant (Koch)
	CoAscendant2  float64   `json:"co_ascendant2"`  // Co-ascendant (Munkasey)
	PolarAsc      float64   `json:"polar_asc"`      // Polar ascendant
}

// IsRetrograde returns true if the planet is moving retrograde
func (p PlanetPosition) IsRetrograde() bool {
	return p.LongSpeed < 0
}

// GetDegreeInSign returns the degree within the zodiac sign
func (p PlanetPosition) GetDegreeInSign() float64 {
	return math.Mod(p.Longitude, 30.0)
}

// GetSign returns the zodiac sign for this position
func (p PlanetPosition) GetSign() string {
	return domain.GetZodiacSign(p.Longitude)
}

// ToDomainPlanet converts ephemeris data to domain planet
func (p PlanetPosition) ToDomainPlanet(ephemeris *Ephemeris, houseNumber int) domain.Planet {
	return domain.NewPlanet(
		ephemeris.GetPlanetName(p.PlanetID),
		p.Longitude,
		p.Latitude,
		p.LongSpeed,
		houseNumber,
	)
}
