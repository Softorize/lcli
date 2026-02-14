package model

import (
	"errors"
	"fmt"
)

// Sentinel errors for common LinkedIn API failure modes.
var (
	// ErrUnauthorized indicates a 401 response (invalid or expired token).
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates a 403 response (insufficient permissions).
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound indicates a 404 response (resource does not exist).
	ErrNotFound = errors.New("not found")
	// ErrRateLimited indicates a 429 response (too many requests).
	ErrRateLimited = errors.New("rate limited")
	// ErrServer indicates a 5xx response (LinkedIn server error).
	ErrServer = errors.New("server error")
)

// APIError represents a structured error returned by the LinkedIn API.
type APIError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	Code       string `json:"serviceErrorCode"`
	TraceID    string `json:"traceId"`
}

// Error returns a human-readable representation of the API error.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("linkedin api %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("linkedin api %d: %s", e.StatusCode, e.Message)
}

// Unwrap returns the corresponding sentinel error for the status code,
// allowing callers to use errors.Is for classification.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 429:
		return ErrRateLimited
	default:
		if e.StatusCode >= 500 {
			return ErrServer
		}
		return nil
	}
}
