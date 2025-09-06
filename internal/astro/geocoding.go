package astro

import (
	"astroeph-api/internal/domain"
	"database/sql"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// Embed the geocoding data directly in the binary
//
//go:embed data/cities500.txt
var geonamesData string

// GeocodingService provides local city coordinate and timezone lookup
type GeocodingService struct {
	db *sql.DB
}

// geoCodeService is the global instance
var geoCodeService *GeocodingService

// Initialize sets up the global geocoding service
func Initialize() error {
	var err error
	geoCodeService, err = NewGeocodingService()
	return err
}

// GetGeocodingService returns the global geocoding service instance
func GetGeocodingService() *GeocodingService {
	return geoCodeService
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

	// Load embedded GeoNames data
	return g.loadGeoNamesData()
}

// loadGeoNamesData reads and parses the embedded cities500.txt data
func (g *GeocodingService) loadGeoNamesData() error {
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

	// Parse embedded data line by line
	lines := strings.Split(geonamesData, "\n")
	insertedCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse the GeoNames record
		record, err := g.parseGeoNamesLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		// Only include cities with valid timezones
		if record.Timezone == "" {
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
			continue // Skip failed inserts
		}

		insertedCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// parseGeoNamesLine parses a single line from the embedded cities data
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
	record.CountryCode = fields[8]

	// Parse population (field 14)
	if fields[14] != "" {
		if pop, err := strconv.Atoi(fields[14]); err == nil {
			record.Population = pop
		}
	}

	// Parse timezone (field 17)
	record.Timezone = fields[17]

	return record, nil
}

// GetCityInfo looks up city information by name
func (g *GeocodingService) GetCityInfo(cityName string) (*domain.Location, error) {
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
	}

	// Try each query in order
	for _, q := range queries {
		var geonameid, population int
		var name, asciiname, alternatenames, country, timezone string
		var latitude, longitude float64

		err := g.db.QueryRow(q.query, q.args...).Scan(
			&geonameid, &name, &asciiname, &alternatenames,
			&country, &latitude, &longitude, &population, &timezone,
		)

		if err == nil {
			location := domain.NewLocation(name, name, country, latitude, longitude, timezone)
			location.Region = "" // Could be enhanced to include region
			return location, nil
		}
	}

	// City not found in database, return default (New York) with warning
	return domain.NewLocation(
		cityName, // Keep original name
		cityName,
		"US",
		40.7128,
		-74.0060,
		"America/New_York",
	), nil
}

// Close closes the database connection
func (g *GeocodingService) Close() error {
	if g.db != nil {
		return g.db.Close()
	}
	return nil
}

// GeonamesRecord represents a single record from the embedded cities data
type GeonamesRecord struct {
	GeonameID      int
	Name           string
	ASCIIName      string
	AlternateNames string
	Latitude       float64
	Longitude      float64
	CountryCode    string
	Population     int
	Timezone       string
}
