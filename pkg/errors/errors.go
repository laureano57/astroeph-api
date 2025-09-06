package errors

import (
	"fmt"
	"net/http"
)

// AstroError represents an astrological calculation error
type AstroError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

// Error implements the error interface
func (e *AstroError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeCityNotFound       = "CITY_NOT_FOUND"
	ErrCodeEphemerisError     = "EPHEMERIS_ERROR"
	ErrCodeHouseCalculation   = "HOUSE_CALCULATION_ERROR"
	ErrCodePlanetCalculation  = "PLANET_CALCULATION_ERROR"
	ErrCodeAspectCalculation  = "ASPECT_CALCULATION_ERROR"
	ErrCodeChartGeneration    = "CHART_GENERATION_ERROR"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// NewAstroError creates a new astrological error
func NewAstroError(code, message string, httpStatus int) *AstroError {
	return &AstroError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Predefined errors
var (
	ErrInvalidDateRange = NewAstroError(
		ErrCodeInvalidInput,
		"Date must be between 1800 and 2200",
		http.StatusBadRequest,
	)

	ErrInvalidTime = NewAstroError(
		ErrCodeInvalidInput,
		"Time must be in HH:MM:SS format",
		http.StatusBadRequest,
	)

	ErrCityNotFound = NewAstroError(
		ErrCodeCityNotFound,
		"City not found in database",
		http.StatusNotFound,
	)

	ErrEphemerisNotInitialized = NewAstroError(
		ErrCodeEphemerisError,
		"Swiss Ephemeris not initialized",
		http.StatusInternalServerError,
	)

	ErrInvalidHouseSystem = NewAstroError(
		ErrCodeInvalidInput,
		"Invalid house system specified",
		http.StatusBadRequest,
	)

	ErrChartGenerationFailed = NewAstroError(
		ErrCodeChartGeneration,
		"Failed to generate chart visualization",
		http.StatusInternalServerError,
	)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 1 {
		return fmt.Sprintf("Validation error: %s", ve.Errors[0].Message)
	}
	return fmt.Sprintf("Multiple validation errors (%d)", len(ve.Errors))
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsEmpty returns true if there are no validation errors
func (ve *ValidationErrors) IsEmpty() bool {
	return len(ve.Errors) == 0
}

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, code, message string) *AstroError {
	if astroErr, ok := err.(*AstroError); ok {
		return astroErr
	}

	return &AstroError{
		Code:       code,
		Message:    fmt.Sprintf("%s: %s", message, err.Error()),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// IsAstroError checks if an error is an AstroError
func IsAstroError(err error) bool {
	_, ok := err.(*AstroError)
	return ok
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if astroErr, ok := err.(*AstroError); ok {
		return astroErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
