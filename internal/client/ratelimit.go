// ratelimit.go provides rate limit aware request execution.
package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RateLimit holds the rate limit information parsed from LinkedIn API response
// headers.
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// ParseRateLimit extracts rate limit information from the response headers.
// Returns nil if the headers are not present.
func ParseRateLimit(resp *http.Response) *RateLimit {
	limitStr := resp.Header.Get("X-RateLimit-Limit")
	remainStr := resp.Header.Get("X-RateLimit-Remaining")
	resetStr := resp.Header.Get("X-RateLimit-Reset")

	if limitStr == "" && remainStr == "" && resetStr == "" {
		return nil
	}

	rl := &RateLimit{}

	if v, err := strconv.Atoi(limitStr); err == nil {
		rl.Limit = v
	}
	if v, err := strconv.Atoi(remainStr); err == nil {
		rl.Remaining = v
	}
	if v, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
		rl.Reset = time.Unix(v, 0)
	}

	return rl
}

// CheckRateLimit inspects the response for a 429 status and returns a
// descriptive error that includes how long the caller should wait before
// retrying. For non-429 responses it returns nil.
func CheckRateLimit(resp *http.Response) error {
	if resp.StatusCode != http.StatusTooManyRequests {
		return nil
	}

	rl := ParseRateLimit(resp)
	if rl != nil && !rl.Reset.IsZero() {
		wait := time.Until(rl.Reset).Truncate(time.Second)
		if wait < 0 {
			wait = 0
		}
		return fmt.Errorf("rate limited (429): retry after %s", wait)
	}

	return fmt.Errorf("rate limited (429): retry later")
}
