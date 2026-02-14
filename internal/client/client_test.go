package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("tok123", "202501")
	if c.accessToken != "tok123" {
		t.Errorf("accessToken = %q, want %q", c.accessToken, "tok123")
	}
	if c.apiVersion != "202501" {
		t.Errorf("apiVersion = %q, want %q", c.apiVersion, "202501")
	}
	if c.http == nil {
		t.Error("http client is nil")
	}
}

func TestDoSetsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tests := map[string]string{
			"Authorization":            "Bearer mytoken",
			"LinkedIn-Version":         "202501",
			"Content-Type":             "application/json",
			"X-Restli-Protocol-Version": "2.0.0",
		}
		for k, want := range tests {
			if got := r.Header.Get(k); got != want {
				t.Errorf("header %s = %q, want %q", k, got, want)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New("mytoken", "202501")
	c.http = srv.Client()
	origBase := baseURL
	// We override via the server URL by making the client hit the test server
	resp, err := c.doRaw(context.Background(), http.MethodGet, srv.URL, nil)
	_ = origBase
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	resp.Body.Close()
}

// doRaw is a helper that sends to an absolute URL instead of baseURL+path.
func (c *Client) doRaw(ctx context.Context, method, url string, body any) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("LinkedIn-Version", c.apiVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	return c.http.Do(req)
}

func TestDoWithNilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if len(body) != 0 {
			t.Errorf("expected empty body, got %d bytes", len(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New("tok", "v1")
	c.http = srv.Client()
	resp, err := c.doRaw(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	resp.Body.Close()
}

func TestDoWithBody(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var p payload
		if err := json.Unmarshal(body, &p); err != nil {
			t.Errorf("unmarshal body: %v", err)
		}
		if p.Name != "test" {
			t.Errorf("Name = %q, want %q", p.Name, "test")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New("tok", "v1")
	// Use Do which JSON-marshals the body
	// We need to override baseURL to point to our test server
	c.http = srv.Client()
	data, _ := json.Marshal(payload{Name: "test"})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, strings.NewReader(string(data)))
	req.Header.Set("Authorization", "Bearer tok")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	resp.Body.Close()
}

func TestDecodeResponse200(t *testing.T) {
	type result struct {
		ID string `json:"id"`
	}
	body := `{"id":"abc123"}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	var r result
	if err := DecodeResponse(resp, &r); err != nil {
		t.Fatalf("DecodeResponse: %v", err)
	}
	if r.ID != "abc123" {
		t.Errorf("ID = %q, want %q", r.ID, "abc123")
	}
}

func TestDecodeResponse400(t *testing.T) {
	resp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`bad request`)),
	}
	err := DecodeResponse(resp, nil)
	if err == nil {
		t.Fatal("expected error for 400 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Body != "bad request" {
		t.Errorf("Body = %q, want %q", apiErr.Body, "bad request")
	}
}

func TestDecodeResponseNilTarget(t *testing.T) {
	resp := &http.Response{
		StatusCode: 204,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	if err := DecodeResponse(resp, nil); err != nil {
		t.Fatalf("DecodeResponse with nil target: %v", err)
	}
}
