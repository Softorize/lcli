package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigDefaults(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		{"RedirectURI", "redirect", defaultRedirect},
		{"APIVersion", "version", defaultVersion},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{RedirectURI: defaultRedirect, APIVersion: defaultVersion}
			var got string
			if tt.field == "redirect" {
				got = cfg.RedirectURI
			} else {
				got = cfg.APIVersion
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTokenValid(t *testing.T) {
	tests := []struct {
		name string
		tok  Token
		want bool
	}{
		{"valid token", Token{AccessToken: "abc", ExpiresAt: time.Now().Add(time.Hour)}, true},
		{"expired token", Token{AccessToken: "abc", ExpiresAt: time.Now().Add(-time.Hour)}, false},
		{"empty access token", Token{AccessToken: "", ExpiresAt: time.Now().Add(time.Hour)}, false},
		{"expires within buffer", Token{AccessToken: "abc", ExpiresAt: time.Now().Add(2 * time.Minute)}, false},
		{"expires exactly at buffer", Token{AccessToken: "abc", ExpiresAt: time.Now().Add(expiryBuffer)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tok.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveLoadConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{
		ClientID:     "id123",
		ClientSecret: "secret456",
		RedirectURI:  "http://localhost:9999/cb",
		APIVersion:   "202501",
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var loaded Config
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if loaded.ClientID != cfg.ClientID {
		t.Errorf("ClientID = %q, want %q", loaded.ClientID, cfg.ClientID)
	}
	if loaded.ClientSecret != cfg.ClientSecret {
		t.Errorf("ClientSecret = %q, want %q", loaded.ClientSecret, cfg.ClientSecret)
	}
}

func TestSaveLoadTokenRoundTrip(t *testing.T) {
	dir := t.TempDir()
	tok := &Token{
		AccessToken:  "access_abc",
		RefreshToken: "refresh_xyz",
		ExpiresAt:    time.Now().Add(time.Hour).Truncate(time.Second),
		Scopes:       []string{"openid", "profile"},
	}

	data, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(dir, "tokens.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var loaded Token
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if loaded.AccessToken != tok.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, tok.AccessToken)
	}
	if !loaded.ExpiresAt.Equal(tok.ExpiresAt) {
		t.Errorf("ExpiresAt = %v, want %v", loaded.ExpiresAt, tok.ExpiresAt)
	}
}

func TestLoadNonExistentFileReturnsDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nope.json")
	_, err := os.ReadFile(path)
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}

func TestLoadTokenNonExistentReturnsNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tokens.json")
	_, err := os.ReadFile(path)
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}

func TestConfigDirExists(t *testing.T) {
	dir := ConfigDir()
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if !info.IsDir() {
		t.Error("config dir is not a directory")
	}
}
