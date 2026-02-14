package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// mockDoer implements Doer for testing. It returns responses in order.
type mockDoer struct {
	responses []mockResponse
	callIdx   int
}

type mockResponse struct {
	status int
	body   any
}

func (m *mockDoer) Do(_ context.Context, _, _ string, _ any) (*http.Response, error) {
	if m.callIdx >= len(m.responses) {
		return nil, fmt.Errorf("no more mock responses")
	}
	r := m.responses[m.callIdx]
	m.callIdx++
	data, _ := json.Marshal(r.body)
	return &http.Response{
		StatusCode: r.status,
		Body:       io.NopCloser(bytes.NewReader(data)),
		Header:     http.Header{},
	}, nil
}

func TestMeSuccess(t *testing.T) {
	profileResp := map[string]any{
		"id":                 "abc123",
		"localizedFirstName": "John",
		"localizedLastName":  "Doe",
		"localizedHeadline":  "Engineer",
		"vanityName":         "johndoe",
	}
	emailResp := map[string]any{
		"elements": []map[string]any{
			{"handle~": map[string]any{"emailAddress": "john@example.com"}},
		},
	}

	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: profileResp},
		{status: 200, body: emailResp},
	}}

	svc := NewProfileService(doer)
	profile, err := svc.Me(context.Background())
	if err != nil {
		t.Fatalf("Me: %v", err)
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ID", profile.ID, "abc123"},
		{"FirstName", profile.FirstName, "John"},
		{"LastName", profile.LastName, "Doe"},
		{"Headline", profile.Headline, "Engineer"},
		{"Vanity", profile.Vanity, "johndoe"},
		{"Email", profile.Email, "john@example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestMeNon200(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"401 unauthorized", 401},
		{"403 forbidden", 403},
		{"500 server error", 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errBody := map[string]any{
				"status":  tt.status,
				"message": "error",
			}
			doer := &mockDoer{responses: []mockResponse{
				{status: tt.status, body: errBody},
			}}

			svc := NewProfileService(doer)
			_, err := svc.Me(context.Background())
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestMeEmailFetchFails(t *testing.T) {
	profileResp := map[string]any{
		"id":                 "xyz",
		"localizedFirstName": "Jane",
		"localizedLastName":  "Smith",
	}
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: profileResp},
		{status: 403, body: map[string]any{"status": 403, "message": "no access"}},
	}}

	svc := NewProfileService(doer)
	profile, err := svc.Me(context.Background())
	if err != nil {
		t.Fatalf("Me should succeed even if email fails: %v", err)
	}
	if profile.Email != "" {
		t.Errorf("Email = %q, want empty when fetch fails", profile.Email)
	}
}
