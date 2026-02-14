package client

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func makeResp(status int, headers map[string]string) *http.Response {
	resp := &http.Response{StatusCode: status, Header: http.Header{}}
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}

func TestParseRateLimitValid(t *testing.T) {
	h := map[string]string{
		"X-RateLimit-Limit":     "100",
		"X-RateLimit-Remaining": "42",
		"X-RateLimit-Reset":     strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10),
	}
	rl := ParseRateLimit(makeResp(200, h))
	if rl == nil {
		t.Fatal("expected non-nil RateLimit")
	}
	if rl.Limit != 100 {
		t.Errorf("Limit = %d, want 100", rl.Limit)
	}
	if rl.Remaining != 42 {
		t.Errorf("Remaining = %d, want 42", rl.Remaining)
	}
}

func TestParseRateLimitMissing(t *testing.T) {
	rl := ParseRateLimit(makeResp(200, map[string]string{}))
	if rl != nil {
		t.Errorf("expected nil, got %+v", rl)
	}
}

func TestParseRateLimitPartial(t *testing.T) {
	tests := []struct {
		name      string
		headers   map[string]string
		wantLimit int
		wantRemn  int
	}{
		{"limit only", map[string]string{"X-RateLimit-Limit": "50"}, 50, 0},
		{"remaining only", map[string]string{"X-RateLimit-Remaining": "10"}, 0, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := ParseRateLimit(makeResp(200, tt.headers))
			if rl == nil {
				t.Fatal("expected non-nil RateLimit")
			}
			if rl.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", rl.Limit, tt.wantLimit)
			}
			if rl.Remaining != tt.wantRemn {
				t.Errorf("Remaining = %d, want %d", rl.Remaining, tt.wantRemn)
			}
		})
	}
}

func TestCheckRateLimit200(t *testing.T) {
	err := CheckRateLimit(makeResp(200, nil))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheckRateLimit429WithReset(t *testing.T) {
	reset := strconv.FormatInt(time.Now().Add(30*time.Second).Unix(), 10)
	h := map[string]string{
		"X-RateLimit-Limit":     "100",
		"X-RateLimit-Remaining": "0",
		"X-RateLimit-Reset":     reset,
	}
	err := CheckRateLimit(makeResp(429, h))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "retry after") {
		t.Errorf("error %q should contain 'retry after'", err.Error())
	}
}

func TestCheckRateLimit429NoReset(t *testing.T) {
	err := CheckRateLimit(makeResp(429, map[string]string{}))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "retry later") {
		t.Errorf("error %q should contain 'retry later'", err.Error())
	}
}
