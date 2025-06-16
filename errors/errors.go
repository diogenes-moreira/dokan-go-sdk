package errors

import (
	"fmt"
	"net/http"
)

// DokanError represents a Dokan API error
type DokanError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
}

// Error implements the error interface
func (e *DokanError) Error() string {
	return fmt.Sprintf("dokan api error: %s - %s", e.Code, e.Message)
}

// IsDokanError checks if an error is a DokanError
func IsDokanError(err error) bool {
	_, ok := err.(*DokanError)
	return ok
}

// NetworkError represents a network-related error
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication error
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s with ID %v", e.Resource, e.ID)
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %d seconds", e.RetryAfter)
}

// NewDokanError creates a new DokanError
func NewDokanError(code, message string, statusCode int) *DokanError {
	return &DokanError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(err error) *NetworkError {
	return &NetworkError{Err: err}
}

// NewAuthenticationError creates a new AuthenticationError
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{Message: message}
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, code, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// NewRateLimitError creates a new RateLimitError
func NewRateLimitError(retryAfter int) *RateLimitError {
	return &RateLimitError{RetryAfter: retryAfter}
}

// HandleHTTPError converts HTTP status codes to appropriate errors
func HandleHTTPError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return NewAuthenticationError("unauthorized access")
	case http.StatusForbidden:
		return NewAuthenticationError("forbidden access")
	case http.StatusNotFound:
		return NewNotFoundError("resource", "unknown")
	case http.StatusTooManyRequests:
		return NewRateLimitError(60) // Default retry after 60 seconds
	case http.StatusBadRequest:
		return NewDokanError("bad_request", "bad request", statusCode)
	case http.StatusInternalServerError:
		return NewDokanError("internal_error", "internal server error", statusCode)
	default:
		if statusCode >= 400 {
			return NewDokanError("http_error", fmt.Sprintf("HTTP %d error", statusCode), statusCode)
		}
		return nil
	}
}

