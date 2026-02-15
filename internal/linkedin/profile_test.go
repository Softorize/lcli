package linkedin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMeSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"sub":            "abc123",
			"given_name":     "John",
			"family_name":    "Doe",
			"name":           "John Doe",
			"picture":        "https://example.com/photo.jpg",
			"email":          "john@example.com",
			"email_verified": true,
		})
	}))
	defer srv.Close()

	// Override the userinfo URL for testing.
	origURL := userinfoURL
	defer func() { userinfoURLOverride = "" }()
	userinfoURLOverride = srv.URL
	_ = origURL

	doer := &mockDoer{}
	svc := NewProfileService(doer, "test-token")
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"error": "forbidden"}`))
	}))
	defer srv.Close()

	defer func() { userinfoURLOverride = "" }()
	userinfoURLOverride = srv.URL

	doer := &mockDoer{}
	svc := NewProfileService(doer, "bad-token")
	_, err := svc.Me(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetByIDSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id":                 "person123",
			"localizedFirstName": "Alice",
			"localizedLastName":  "Jones",
			"localizedHeadline":  "Designer",
			"vanityName":         "alicejones",
		}},
	}}

	svc := NewProfileService(doer, "")
	profile, err := svc.GetByID(context.Background(), "person123")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if profile.ID != "person123" {
		t.Errorf("ID = %q", profile.ID)
	}
	if profile.FirstName != "Alice" {
		t.Errorf("FirstName = %q", profile.FirstName)
	}
}

func TestGetByIDError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewProfileService(doer, "")
	_, err := svc.GetByID(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}
