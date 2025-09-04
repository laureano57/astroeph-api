package services

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// GeocodingService provides local city coordinate and timezone lookup
type GeocodingService struct {
	db *sql.DB
}

// CityInfo represents geographic and timezone information for a city
type CityInfo struct {
	GeonameID      int     `json:"geoname_id"`
	Name           string  `json:"name"`
	ASCIIName      string  `json:"ascii_name"`
	AlternateNames string  `json:"alternate_names"`
	Country        string  `json:"country"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Population     int     `json:"population"`
	Timezone       string  `json:"timezone"`
}

// GeonamesRecord represents a single record from the cities500.txt file
type GeonamesRecord struct {
	GeonameID        int
	Name             string
	ASCIIName        string
	AlternateNames   string
	Latitude         float64
	Longitude        float64
	FeatureClass     string
	FeatureCode      string
	CountryCode      string
	CC2              string
	Admin1Code       string
	Admin2Code       string
	Admin3Code       string
	Admin4Code       string
	Population       int
	Elevation        int
	DEM              int
	Timezone         string
	ModificationDate string
}

// NewGeocodingService creates a new geocoding service with embedded SQLite database
func NewGeocodingService() (*GeocodingService, error) {
	// Create in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding database: %w", err)
	}

	service := &GeocodingService{db: db}

	// Initialize the database with city data
	if err := service.initializeDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize geocoding database: %w", err)
	}

	if AppLogger != nil {
		AppLogger.Info().Msg("🌍 Geocoding service initialized with local database")
	}

	return service, nil
}

// initializeDatabase creates tables and populates with GeoNames data
func (g *GeocodingService) initializeDatabase() error {
	// Create cities table based on GeoNames structure
	createTableSQL := `
		CREATE TABLE cities (
			geonameid INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			asciiname TEXT NOT NULL,
			alternatenames TEXT,
			country TEXT NOT NULL,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			population INTEGER,
			timezone TEXT NOT NULL
		);
		
		CREATE INDEX idx_city_name ON cities(name);
		CREATE INDEX idx_city_asciiname ON cities(asciiname);
		CREATE INDEX idx_city_alternatenames ON cities(alternatenames);
		CREATE INDEX idx_city_country ON cities(country);
		CREATE INDEX idx_city_population ON cities(population);
	`

	if _, err := g.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create cities table: %w", err)
	}

	// Load GeoNames data from cities500.txt
	return g.loadGeoNamesData()
}

// loadGeoNamesData reads and parses the cities500.txt file
func (g *GeocodingService) loadGeoNamesData() error {
	// Find the cities500.txt file
	dataPath := filepath.Join("data", "geocoding", "cities500.txt")

	if AppLogger != nil {
		AppLogger.Info().
			Str("file_path", dataPath).
			Msg("📂 Loading GeoNames data from cities500.txt")
	}

	file, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("failed to open cities500.txt: %w", err)
	}
	defer file.Close()

	// Begin transaction for better performance
	tx, err := g.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Prepare insert statement within the transaction
	insertSQL := `INSERT INTO cities (geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone) 
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	insertedCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		// Parse the GeoNames record
		record, err := g.parseGeoNamesLine(line)
		if err != nil {
			if AppLogger != nil {
				AppLogger.Debug().
					Err(err).
					Int("line_number", lineCount).
					Msg("Skipping invalid line in cities500.txt")
			}
			continue
		}

		// Only include cities with valid timezones and reasonable population
		if record.Timezone == "" || len(record.Timezone) == 0 {
			continue
		}

		// Insert the record
		_, err = stmt.Exec(
			record.GeonameID,
			record.Name,
			record.ASCIIName,
			record.AlternateNames,
			record.CountryCode,
			record.Latitude,
			record.Longitude,
			record.Population,
			record.Timezone,
		)

		if err != nil {
			if AppLogger != nil {
				AppLogger.Debug().
					Err(err).
					Str("city_name", record.Name).
					Int("geoname_id", record.GeonameID).
					Msg("Failed to insert city record")
			}
			continue
		}

		insertedCount++

		// Log progress every 10,000 cities
		if insertedCount%10000 == 0 {
			if AppLogger != nil {
				AppLogger.Info().
					Int("cities_loaded", insertedCount).
					Int("lines_processed", lineCount).
					Msg("🔄 GeoNames loading progress...")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading cities500.txt: %w", err)
	}

	if AppLogger != nil {
		AppLogger.Info().
			Int("total_cities_loaded", insertedCount).
			Int("total_lines_processed", lineCount).
			Msg("✅ GeoNames database populated successfully")
	}

	return nil
}

// parseGeoNamesLine parses a single line from the cities500.txt file
func (g *GeocodingService) parseGeoNamesLine(line string) (*GeonamesRecord, error) {
	fields := strings.Split(line, "\t")

	// GeoNames format has 19 fields
	if len(fields) < 19 {
		return nil, fmt.Errorf("invalid line format: expected 19 fields, got %d", len(fields))
	}

	record := &GeonamesRecord{}

	// Parse geonameid (field 0)
	if id, err := strconv.Atoi(fields[0]); err != nil {
		return nil, fmt.Errorf("invalid geoname ID: %w", err)
	} else {
		record.GeonameID = id
	}

	// Parse strings (fields 1-3)
	record.Name = fields[1]
	record.ASCIIName = fields[2]
	record.AlternateNames = fields[3]

	// Parse latitude (field 4)
	if lat, err := strconv.ParseFloat(fields[4], 64); err != nil {
		return nil, fmt.Errorf("invalid latitude: %w", err)
	} else {
		record.Latitude = lat
	}

	// Parse longitude (field 5)
	if lon, err := strconv.ParseFloat(fields[5], 64); err != nil {
		return nil, fmt.Errorf("invalid longitude: %w", err)
	} else {
		record.Longitude = lon
	}

	// Parse other fields
	record.FeatureClass = fields[6]
	record.FeatureCode = fields[7]
	record.CountryCode = fields[8]
	record.CC2 = fields[9]
	record.Admin1Code = fields[10]
	record.Admin2Code = fields[11]
	record.Admin3Code = fields[12]
	record.Admin4Code = fields[13]

	// Parse population (field 14)
	if fields[14] != "" {
		if pop, err := strconv.Atoi(fields[14]); err == nil {
			record.Population = pop
		}
	}

	// Parse elevation (field 15)
	if fields[15] != "" {
		if elev, err := strconv.Atoi(fields[15]); err == nil {
			record.Elevation = elev
		}
	}

	// Parse DEM (field 16)
	if fields[16] != "" {
		if dem, err := strconv.Atoi(fields[16]); err == nil {
			record.DEM = dem
		}
	}

	// Parse timezone (field 17)
	record.Timezone = fields[17]

	// Parse modification date (field 18)
	record.ModificationDate = fields[18]

	return record, nil
}

// GetCityInfo looks up city information by name (case-insensitive, fuzzy matching across multiple fields)
func (g *GeocodingService) GetCityInfo(cityName string) (*CityInfo, error) {
	if AppLogger != nil {
		AppLogger.Debug().
			Str("city_query", cityName).
			Msg("🔍 Looking up city coordinates in GeoNames database")
	}

	// Define the search queries in order of preference
	queries := []struct {
		name  string
		query string
		args  []interface{}
	}{
		{
			"exact name match",
			`SELECT geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone 
			 FROM cities WHERE LOWER(name) = LOWER(?) ORDER BY population DESC LIMIT 1`,
			[]interface{}{cityName},
		},
		{
			"exact ascii name match",
			`SELECT geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone 
			 FROM cities WHERE LOWER(asciiname) = LOWER(?) ORDER BY population DESC LIMIT 1`,
			[]interface{}{cityName},
		},
		{
			"alternate names exact match",
			`SELECT geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone 
			 FROM cities WHERE LOWER(alternatenames) LIKE LOWER(?) ORDER BY population DESC LIMIT 1`,
			[]interface{}{"%" + cityName + "%"},
		},
		{
			"fuzzy name match",
			`SELECT geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone 
			 FROM cities WHERE LOWER(name) LIKE LOWER(?) ORDER BY population DESC LIMIT 1`,
			[]interface{}{"%" + cityName + "%"},
		},
		{
			"fuzzy ascii name match",
			`SELECT geonameid, name, asciiname, alternatenames, country, latitude, longitude, population, timezone 
			 FROM cities WHERE LOWER(asciiname) LIKE LOWER(?) ORDER BY population DESC LIMIT 1`,
			[]interface{}{"%" + cityName + "%"},
		},
	}

	// Try each query in order
	for _, q := range queries {
		var city CityInfo

		err := g.db.QueryRow(q.query, q.args...).Scan(
			&city.GeonameID,
			&city.Name,
			&city.ASCIIName,
			&city.AlternateNames,
			&city.Country,
			&city.Latitude,
			&city.Longitude,
			&city.Population,
			&city.Timezone,
		)

		if err == nil {
			if AppLogger != nil {
				AppLogger.Info().
					Str("query_type", q.name).
					Str("query", cityName).
					Str("matched_city", city.Name).
					Str("ascii_name", city.ASCIIName).
					Str("country", city.Country).
					Int("population", city.Population).
					Float64("lat", city.Latitude).
					Float64("lon", city.Longitude).
					Str("timezone", city.Timezone).
					Int("geoname_id", city.GeonameID).
					Msg("✅ City found in GeoNames database")
			}
			return &city, nil
		}
	}

	// City not found in database, return default (New York) with warning
	if AppLogger != nil {
		AppLogger.Warn().
			Str("city_query", cityName).
			Msg("🌍 City not found in GeoNames database, using default coordinates (New York)")
	}

	return &CityInfo{
		GeonameID:      5128581,  // New York City GeoName ID
		Name:           cityName, // Keep original name
		ASCIIName:      cityName,
		AlternateNames: "",
		Country:        "US",
		Latitude:       40.7128,
		Longitude:      -74.0060,
		Population:     8175133,
		Timezone:       "America/New_York",
	}, nil
}

// Close closes the database connection
func (g *GeocodingService) Close() error {
	if g.db != nil {
		return g.db.Close()
	}
	return nil
}

// Global geocoding service instance
var GeoService *GeocodingService

// InitializeGeocodingService sets up the global geocoding service
func InitializeGeocodingService() error {
	var err error
	GeoService, err = NewGeocodingService()
	return err
}
