package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger provides structured logging for the application
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger instance with proper configuration
func NewLogger() *Logger {
	// Configure zerolog for beautiful console output in development
	zerolog.TimeFieldFormat = time.RFC3339

	// Use console writer for better readability in development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		FormatLevel: func(i interface{}) string {
			switch i {
			case "info":
				return "📘 INFO"
			case "warn":
				return "⚠️  WARN"
			case "error":
				return "❌ ERROR"
			case "debug":
				return "🔍 DEBUG"
			default:
				return "📝 " + i.(string)
			}
		},
		FormatCaller: func(i interface{}) string {
			return "📍 " + i.(string)
		},
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", "astroeph-api").
		Logger()

	return &Logger{logger: logger}
}

// Info logs an info message
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn logs a warning message
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error logs an error message
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Debug logs a debug message
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// With creates a new logger with additional fields
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// RequestLogger logs HTTP request details
func (l *Logger) RequestLogger() *zerolog.Event {
	return l.logger.Info().Str("type", "request")
}

// CalculationLogger logs astrological calculation details
func (l *Logger) CalculationLogger() *zerolog.Event {
	return l.logger.Info().Str("type", "calculation")
}
