package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField    ErrorCode = "MISSING_FIELD"
	
	// Business logic errors
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	
	// System errors
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout         ErrorCode = "TIMEOUT"
	
	// Cache errors
	ErrCodeCacheMiss       ErrorCode = "CACHE_MISS"
	ErrCodeCacheError      ErrorCode = "CACHE_ERROR"
	
	// A/B testing errors
	ErrCodeExperimentNotFound ErrorCode = "EXPERIMENT_NOT_FOUND"
	ErrCodeVariantNotFound    ErrorCode = "VARIANT_NOT_FOUND"
)

// AppError represents application-specific errors
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

// NewServiceUnavailableError creates a service unavailable error
func NewServiceUnavailableError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeServiceUnavailable,
		Message:    message,
		HTTPStatus: http.StatusServiceUnavailable,
	}
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeTimeout,
		Message:    message,
		HTTPStatus: http.StatusRequestTimeout,
	}
}

// NewCacheError creates a cache error
func NewCacheError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeCacheError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

// NewExperimentNotFoundError creates an experiment not found error
func NewExperimentNotFoundError(id string) *AppError {
	return &AppError{
		Code:       ErrCodeExperimentNotFound,
		Message:    fmt.Sprintf("Experiment %s not found", id),
		HTTPStatus: http.StatusNotFound,
	}
}

// NewVariantNotFoundError creates a variant not found error
func NewVariantNotFoundError(id string) *AppError {
	return &AppError{
		Code:       ErrCodeVariantNotFound,
		Message:    fmt.Sprintf("Variant %s not found", id),
		HTTPStatus: http.StatusNotFound,
	}
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithCause adds a cause to an error
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}
