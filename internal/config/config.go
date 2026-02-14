// Package config manages lcli configuration and credential storage.
// Config is stored at ~/.config/lcli/config.json
// Tokens are stored at ~/.config/lcli/tokens.json
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	dirName         = "lcli"
	configFileName  = "config.json"
	tokensFileName  = "tokens.json"
	defaultRedirect = "http://localhost:8484/callback"
	defaultVersion  = "202501"
	expiryBuffer    = 5 * time.Minute
)

// Config holds the LinkedIn application credentials and API settings.
type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	APIVersion   string `json:"api_version"`
}

// Token holds OAuth 2.0 credentials obtained from LinkedIn.
type Token struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       []string `json:"scopes"`
}

// Valid reports whether the token is present and not expired.
// A 5-minute buffer is applied so tokens about to expire are treated as invalid.
func (t *Token) Valid() bool {
	if t.AccessToken == "" {
		return false
	}
	return time.Now().Add(expiryBuffer).Before(t.ExpiresAt)
}

// ConfigDir returns the path to the lcli configuration directory
// (~/.config/lcli) and creates it if it does not exist.
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".config", dirName)
	_ = os.MkdirAll(dir, 0o700)
	return dir
}

// Load reads the configuration from config.json.
// If the file does not exist, default values are returned.
func Load() (*Config, error) {
	cfg := &Config{
		RedirectURI: defaultRedirect,
		APIVersion:  defaultVersion,
	}

	data, err := os.ReadFile(filepath.Join(ConfigDir(), configFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// Save writes the configuration to config.json.
func Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	path := filepath.Join(ConfigDir(), configFileName)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// LoadToken reads the OAuth token from tokens.json.
// If the file does not exist, nil is returned without error.
func LoadToken() (*Token, error) {
	data, err := os.ReadFile(filepath.Join(ConfigDir(), tokensFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read token: %w", err)
	}

	var tok Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	return &tok, nil
}

// SaveToken writes the OAuth token to tokens.json.
func SaveToken(tok *Token) error {
	data, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	path := filepath.Join(ConfigDir(), tokensFileName)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write token: %w", err)
	}
	return nil
}
