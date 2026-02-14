package model

import (
	"errors"
	"testing"
)

func TestAPIErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with code",
			err:  APIError{StatusCode: 403, Code: "ACCESS_DENIED", Message: "no access"},
			want: "linkedin api 403 (ACCESS_DENIED): no access",
		},
		{
			name: "without code",
			err:  APIError{StatusCode: 500, Message: "internal error"},
			want: "linkedin api 500: internal error",
		},
		{
			name: "empty message",
			err:  APIError{StatusCode: 400},
			want: "linkedin api 400: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIErrorUnwrapSentinels(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		sentinel error
	}{
		{"401 unauthorized", 401, ErrUnauthorized},
		{"403 forbidden", 403, ErrForbidden},
		{"404 not found", 404, ErrNotFound},
		{"429 rate limited", 429, ErrRateLimited},
		{"500 server error", 500, ErrServer},
		{"502 server error", 502, ErrServer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{StatusCode: tt.status, Message: "test"}
			if !errors.Is(apiErr, tt.sentinel) {
				t.Errorf("errors.Is(%v, %v) = false, want true", apiErr, tt.sentinel)
			}
		})
	}
}

func TestAPIErrorUnwrapNilForUnknownStatus(t *testing.T) {
	apiErr := &APIError{StatusCode: 418, Message: "teapot"}
	// Should not match any sentinel
	sentinels := []error{ErrUnauthorized, ErrForbidden, ErrNotFound, ErrRateLimited, ErrServer}
	for _, s := range sentinels {
		if errors.Is(apiErr, s) {
			t.Errorf("errors.Is(418, %v) = true, want false", s)
		}
	}
}
