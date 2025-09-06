package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string // For SQLite database path if needed in the future
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables and defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Host: getEnvOrDefault("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Path: getEnvOrDefault("DB_PATH", "data/astroeph.db"),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "console"),
		},
	}
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
