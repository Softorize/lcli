// Package auth implements LinkedIn OAuth 2.0 three-legged authorization.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/toto/lcli/internal/config"
)

const (
	authEndpoint  = "https://www.linkedin.com/oauth/v2/authorization"
	tokenEndpoint = "https://www.linkedin.com/oauth/v2/accessToken"
)

// defaultScopes are the OAuth scopes requested during authorization.
var defaultScopes = []string{"openid", "profile", "email", "w_member_social"}

// Authenticator handles LinkedIn OAuth 2.0 flows.
type Authenticator struct {
	cfg *config.Config
}

// NewAuthenticator creates an Authenticator using the provided configuration.
func NewAuthenticator(cfg *config.Config) *Authenticator {
	return &Authenticator{cfg: cfg}
}

// AuthorizationURL builds the LinkedIn authorization URL with the given state
// parameter for CSRF protection.
func (a *Authenticator) AuthorizationURL(state string) string {
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {a.cfg.ClientID},
		"redirect_uri":  {a.cfg.RedirectURI},
		"state":         {state},
		"scope":         {strings.Join(defaultScopes, " ")},
	}
	return authEndpoint + "?" + params.Encode()
}

// tokenResponse is the JSON body returned by the token endpoint.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// Exchange trades an authorization code for an access token.
func (a *Authenticator) Exchange(ctx context.Context, code string) (*config.Token, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {a.cfg.RedirectURI},
		"client_id":     {a.cfg.ClientID},
		"client_secret": {a.cfg.ClientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("parse token response: %w", err)
	}

	tok := &config.Token{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
		Scopes:       strings.Fields(tr.Scope),
	}
	return tok, nil
}

// RandomState generates a cryptographically random hex string for use as an
// OAuth state parameter.
func RandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random state: %w", err)
	}
	return hex.EncodeToString(b), nil
}
