package services

import (
	"fmt"
	"math"
	"time"

	"astroeph-api/models"

	"github.com/mshafiee/swephgo"
)

// AstrologyService handles all astrological calculations using swephgo
type AstrologyService struct {
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

// House system constants
const (
	SE_HOUSE_PLACIDUS      = 'P'
	SE_HOUSE_KOCH          = 'K'
	SE_HOUSE_PORPHYRIUS    = 'O'
	SE_HOUSE_REGIOMONTANUS = 'R'
	SE_HOUSE_CAMPANUS      = 'C'
	SE_HOUSE_EQUAL         = 'E'
	SE_HOUSE_WHOLE_SIGN    = 'W'
)

var planetNames = map[int]string{
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

var mainPlanets = []int{
	SE_SUN, SE_MOON, SE_MERCURY, SE_VENUS, SE_MARS,
	SE_JUPITER, SE_SATURN, SE_URANUS, SE_NEPTUNE, SE_PLUTO,
}

// NewAstrologyService creates and initializes a new astrology service
func NewAstrologyService() (*AstrologyService, error) {
	service := &AstrologyService{}

	// Initialize swephgo
	// Set ephemeris path - in production, you might want to set a specific path
	swephgo.SetEphePath([]byte("")) // Use built-in ephemeris data

	// Test Swiss Ephemeris initialization
	if AppLogger != nil {
		AppLogger.Info().Msg("🔮 Initializing Swiss Ephemeris")

		// Test calculation to verify Swiss Ephemeris is working
		testJD := swephgo.Julday(2000, 1, 1, 12.0, 1)
		xx := make([]float64, 6)
		serr := make([]byte, 256)
		result := swephgo.Calc(testJD, 0, 0, xx, serr)

		if result < 0 {
			AppLogger.Error().
				Int("result_code", int(result)).
				Str("error", string(serr)).
				Msg("Swiss Ephemeris test calculation failed")
			return nil, fmt.Errorf("Swiss Ephemeris initialization failed: %s", string(serr))
		}

		AppLogger.Info().
			Float64("test_sun_longitude", xx[0]).
			Str("ephemeris_status", string(serr)).
			Msg("✅ Swiss Ephemeris initialized successfully")
	}

	service.initialized = true
	return service, nil
}

// CalculateNatalChart calculates a complete natal chart
func (s *AstrologyService) CalculateNatalChart(req *models.NatalChartRequest) (*models.NatalChartResponse, error) {
	if !s.initialized {
		return nil, fmt.Errorf("astrology service not initialized")
	}

	// Get coordinates for the city using geocoding service
	cityInfo, err := s.getCityInformation(req.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get city information for %s: %w", req.City, err)
	}

	lat, lon, timezone := cityInfo.Latitude, cityInfo.Longitude, cityInfo.Timezone

	// Parse the local time and create UTC time
	utcTime, err := s.parseTimeToUTC(req.Year, req.Month, req.Day, req.LocalTime, timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	// Convert to Julian Day using UTC time
	julianDay := swephgo.Julday(utcTime.Year(), int(utcTime.Month()), utcTime.Day(),
		float64(utcTime.Hour())+float64(utcTime.Minute())/60.0+float64(utcTime.Second())/3600.0, 1)

	// Calculate planet positions using Swiss Ephemeris
	planets, err := s.calculatePlanetPositions(julianDay)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate planet positions: %w", err)
	}

	// Calculate houses using Swiss Ephemeris
	houseSystem := s.getHouseSystemCode(req.HouseSystem)
	houses, ascendant, midheaven, err := s.calculateHouses(julianDay, lat, lon, houseSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate houses: %w", err)
	}

	// Assign houses to planets
	for i := range planets {
		planets[i].House = s.getHouseForPlanet(planets[i].Longitude, houses)
	}

	// Calculate aspects
	aspects := s.calculateAspects(planets)

	// Build response
	response := &models.NatalChartResponse{
		ChartData: models.ChartData{
			Planets:     planets,
			Houses:      houses,
			Aspects:     aspects,
			Ascendant:   ascendant,
			Midheaven:   midheaven,
			HouseSystem: req.HouseSystem,
		},
		BirthInfo: models.BirthInfo{
			Date:      fmt.Sprintf("%04d-%02d-%02d", req.Year, req.Month, req.Day),
			Time:      req.LocalTime,
			City:      req.City,
			Latitude:  lat,
			Longitude: lon,
		},
		Timezone: timezone,
		UTCTime:  utcTime,
	}

	return response, nil
}

// calculatePlanetPositions calculates positions for all main planets using Swiss Ephemeris
func (s *AstrologyService) calculatePlanetPositions(julianDay float64) ([]models.PlanetPosition, error) {
	var planets []models.PlanetPosition

	for _, planetId := range mainPlanets {
		// Calculate planet position using swephgo
		xx := make([]float64, 6)
		serr := make([]byte, 256)
		result := swephgo.Calc(julianDay, planetId, 0, xx, serr)
		if result < 0 {
			// Only treat negative values as errors, positive values are warnings
			return nil, fmt.Errorf("failed to calculate position for planet %d: %s", planetId, string(serr))
		}

		longitude := xx[0]
		latitude := xx[1]
		distance := xx[2]
		longitudeSpeed := xx[3]

		planet := models.PlanetPosition{
			Name:           planetNames[planetId],
			Longitude:      longitude,
			Latitude:       latitude,
			Distance:       distance,
			LongitudeSpeed: longitudeSpeed,
			Sign:           models.GetZodiacSign(longitude),
			Degree:         models.GetDegreeInSign(longitude),
		}

		planets = append(planets, planet)
	}

	return planets, nil
}

// calculateHouses calculates house cusps using Swiss Ephemeris
func (s *AstrologyService) calculateHouses(julianDay, lat, lon float64, houseSystem rune) ([]models.HouseCusp, float64, float64, error) {
	// Calculate houses using swephgo
	cusps := make([]float64, 13) // 0-12, where 1-12 are the house cusps
	ascmc := make([]float64, 10) // Ascendant, MC, etc.
	result := swephgo.Houses(julianDay, lat, lon, int(houseSystem), cusps, ascmc)
	if result < 0 {
		// Only treat negative values as errors, positive values are warnings
		return nil, 0, 0, fmt.Errorf("failed to calculate houses: house system not supported or invalid parameters")
	}

	var houses []models.HouseCusp

	// Create house cusps (houses 1-12)
	for i := 1; i <= 12; i++ {
		cusp := cusps[i] // swephgo uses 1-based indexing for cusps
		house := models.HouseCusp{
			House: i,
			Cusp:  cusp,
			Sign:  models.GetZodiacSign(cusp),
		}
		houses = append(houses, house)
	}

	ascendant := cusps[1]  // 1st house cusp is the Ascendant
	midheaven := cusps[10] // 10th house cusp is the Midheaven

	return houses, ascendant, midheaven, nil
}

// getHouseForPlanet determines which house a planet is in based on its longitude
func (s *AstrologyService) getHouseForPlanet(longitude float64, houses []models.HouseCusp) int {
	// Simple algorithm to determine house placement
	// This could be more sophisticated in a production system
	for i := 0; i < 12; i++ {
		currentHouse := houses[i].Cusp
		nextHouse := houses[(i+1)%12].Cusp

		// Handle the case where we cross 0 degrees
		if nextHouse < currentHouse {
			nextHouse += 360
			if longitude < currentHouse && longitude >= 0 {
				longitude += 360
			}
		}

		if longitude >= currentHouse && longitude < nextHouse {
			return houses[i].House
		}
	}

	return 1 // Default to 1st house if calculation fails
}

// calculateAspects calculates aspects between planets
func (s *AstrologyService) calculateAspects(planets []models.PlanetPosition) []models.Aspect {
	var aspects []models.Aspect

	// Define major aspects and their orbs
	majorAspects := map[string]struct {
		angle float64
		orb   float64
	}{
		"conjunction": {0, 8},
		"opposition":  {180, 8},
		"trine":       {120, 6},
		"square":      {90, 6},
		"sextile":     {60, 4},
	}

	// Calculate aspects between all planet pairs
	for i, planet1 := range planets {
		for j, planet2 := range planets {
			if i >= j {
				continue // Avoid duplicate pairs and self-aspects
			}

			// Calculate angular difference
			angle := math.Abs(planet1.Longitude - planet2.Longitude)
			if angle > 180 {
				angle = 360 - angle
			}

			// Check for major aspects
			for aspectName, aspectInfo := range majorAspects {
				diff := math.Abs(angle - aspectInfo.angle)
				if diff <= aspectInfo.orb {
					aspect := models.Aspect{
						Planet1:    planet1.Name,
						Planet2:    planet2.Name,
						Type:       aspectName,
						Angle:      angle,
						Orb:        diff,
						IsApplying: planet1.LongitudeSpeed > planet2.LongitudeSpeed,
					}
					aspects = append(aspects, aspect)
					break // Only one aspect per planet pair
				}
			}
		}
	}

	return aspects
}

// Utility functions

// getCityInformation returns detailed city information using the geocoding service
func (s *AstrologyService) getCityInformation(city string) (*CityInfo, error) {
	if GeoService == nil {
		// Fallback to default coordinates if geocoding service is not available
		if AppLogger != nil {
			AppLogger.Warn().
				Str("city", city).
				Msg("Geocoding service not available, using default coordinates")
		}
		return &CityInfo{
			Name:      city,
			Country:   "Unknown",
			Latitude:  40.7128, // New York
			Longitude: -74.0060,
			Timezone:  "America/New_York",
		}, nil
	}

	return GeoService.GetCityInfo(city)
}

// parseTimeToUTC converts local time to UTC
func (s *AstrologyService) parseTimeToUTC(year, month, day int, localTime, timezone string) (time.Time, error) {
	// Parse the time string (HH:MM:SS)
	parsedTime, err := time.Parse("15:04:05", localTime)
	if err != nil {
		return time.Time{}, err
	}

	// Load the timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC // Fallback to UTC
	}

	// Create the local time
	localDateTime := time.Date(year, time.Month(month), day,
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc)

	return localDateTime.UTC(), nil
}

// timeToHours converts HH:MM:SS to decimal hours
func (s *AstrologyService) timeToHours(timeStr string) float64 {
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return 0
	}

	hours := float64(parsedTime.Hour())
	minutes := float64(parsedTime.Minute()) / 60.0
	seconds := float64(parsedTime.Second()) / 3600.0

	return hours + minutes + seconds
}

// getHouseSystemCode converts house system name to swephgo code
func (s *AstrologyService) getHouseSystemCode(system string) rune {
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

// Placeholder functions for other calculations (to be implemented)

func (s *AstrologyService) CalculateTransits(req *models.TransitsRequest) (*models.TransitsResponse, error) {
	return nil, fmt.Errorf("transits calculation not yet implemented")
}

func (s *AstrologyService) CalculateSynastry(req *models.SynastryRequest) (*models.SynastryResponse, error) {
	return nil, fmt.Errorf("synastry calculation not yet implemented")
}
