// Package client provides an HTTP client for the LinkedIn REST API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.linkedin.com/rest"

// Client is an authenticated HTTP client for the LinkedIn REST API.
type Client struct {
	http        *http.Client
	accessToken string
	apiVersion  string
}

// New creates a Client with the given access token and API version.
func New(accessToken, apiVersion string) *Client {
	return &Client{
		http:        &http.Client{},
		accessToken: accessToken,
		apiVersion:  apiVersion,
	}
}

// Do executes an authenticated request against the LinkedIn API.
// If body is non-nil it is JSON-marshalled and sent as the request body.
func (c *Client) Do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("LinkedIn-Version", c.apiVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	return resp, nil
}

// Get performs an authenticated GET request.
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, path, nil)
}

// Post performs an authenticated POST request with a JSON body.
func (c *Client) Post(ctx context.Context, path string, body any) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, path, body)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, http.MethodDelete, path, nil)
}

// APIError is returned when the LinkedIn API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("linkedin api %d: %s", e.StatusCode, e.Body)
}

// DecodeResponse reads the response body and JSON-unmarshals it into v.
// If the response status is not 2xx, an *APIError is returned.
func DecodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}
	return nil
}
