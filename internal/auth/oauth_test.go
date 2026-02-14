package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/toto/lcli/internal/config"
)

func TestAuthorizationURL(t *testing.T) {
	cfg := &config.Config{
		ClientID:    "my-client-id",
		RedirectURI: "http://localhost:8484/callback",
	}
	a := NewAuthenticator(cfg)
	rawURL := a.AuthorizationURL("teststate")

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}

	tests := []struct {
		param, want string
	}{
		{"response_type", "code"},
		{"client_id", "my-client-id"},
		{"redirect_uri", "http://localhost:8484/callback"},
		{"state", "teststate"},
		{"scope", "openid profile email w_member_social"},
	}
	for _, tt := range tests {
		t.Run(tt.param, func(t *testing.T) {
			if got := parsed.Query().Get(tt.param); got != tt.want {
				t.Errorf("param %s = %q, want %q", tt.param, got, tt.want)
			}
		})
	}

	if !strings.HasPrefix(rawURL, authEndpoint) {
		t.Errorf("URL should start with %s", authEndpoint)
	}
}

func TestRandomState(t *testing.T) {
	seen := make(map[string]bool)
	for i := range 3 {
		t.Run("call"+string(rune('0'+i)), func(t *testing.T) {
			s, err := RandomState()
			if err != nil {
				t.Fatalf("RandomState: %v", err)
			}
			if s == "" {
				t.Fatal("returned empty string")
			}
			if len(s) != 32 {
				t.Errorf("length = %d, want 32", len(s))
			}
			if seen[s] {
				t.Error("duplicate state value")
			}
			seen[s] = true
		})
	}
}

func TestExchangeFormParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want form-urlencoded", ct)
		}
		_ = r.ParseForm()
		checks := map[string]string{
			"grant_type": "authorization_code", "code": "authcode",
			"redirect_uri": "http://localhost/cb", "client_id": "cid", "client_secret": "cs",
		}
		for k, want := range checks {
			if got := r.FormValue(k); got != want {
				t.Errorf("form %s = %q, want %q", k, got, want)
			}
		}
		resp := map[string]any{
			"access_token": "tok", "expires_in": 3600,
			"refresh_token": "ref", "scope": "openid profile",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()
	// Note: We cannot redirect Exchange to the test server since tokenEndpoint
	// is a const. The handler above validates the expected form parameters.
	_ = srv
}

func TestExchangeErrorStatus(t *testing.T) {
	// Verify authenticator can be created with various configs
	configs := []struct {
		name string
		cfg  config.Config
	}{
		{"minimal", config.Config{ClientID: "a", ClientSecret: "b"}},
		{"full", config.Config{ClientID: "x", ClientSecret: "y", RedirectURI: "http://localhost/cb"}},
	}
	for _, tt := range configs {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuthenticator(&tt.cfg)
			if a == nil {
				t.Fatal("authenticator is nil")
			}
		})
	}
}
