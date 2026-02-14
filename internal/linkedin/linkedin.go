// Package linkedin implements service-layer calls to the LinkedIn REST API.
package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/toto/lcli/internal/model"
)

// Doer executes HTTP requests against the LinkedIn API.
type Doer interface {
	Do(ctx context.Context, method, path string, body any) (*http.Response, error)
}

// decodeJSON reads the response body and unmarshals it into dst.
func decodeJSON(resp *http.Response, dst any) error {
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

// checkError inspects the HTTP status code and returns a structured
// APIError when the response indicates failure.
func checkError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &model.APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to read error response",
		}
	}

	apiErr := &model.APIError{StatusCode: resp.StatusCode}
	if err := json.Unmarshal(data, apiErr); err != nil {
		apiErr.Message = string(data)
	}

	return apiErr
}

// drainBody reads and closes the response body to allow connection reuse.
func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
