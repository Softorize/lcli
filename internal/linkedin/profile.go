package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Softorize/lcli/internal/model"
)

const userinfoURL = "https://api.linkedin.com/v2/userinfo"

// userinfoURLOverride allows tests to point at a local httptest server.
var userinfoURLOverride string

func getUserinfoURL() string {
	if userinfoURLOverride != "" {
		return userinfoURLOverride
	}
	return userinfoURL
}

// ProfileService provides access to LinkedIn profile endpoints.
type ProfileService struct {
	doer        Doer
	accessToken string
}

// NewProfileService creates a ProfileService backed by the given Doer.
// The accessToken is used for the OpenID Connect /userinfo endpoint.
func NewProfileService(d Doer, accessToken string) *ProfileService {
	return &ProfileService{doer: d, accessToken: accessToken}
}

// userInfoResponse maps the OpenID Connect /userinfo endpoint response.
type userInfoResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// Me returns the authenticated user's profile via OpenID Connect /userinfo.
func (s *ProfileService) Me(ctx context.Context) (*model.Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getUserinfoURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("get my profile: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get my profile: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get my profile: read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get my profile: linkedin api %d: %s", resp.StatusCode, data)
	}

	var info userInfoResponse
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("get my profile: decode: %w", err)
	}

	profile := &model.Profile{
		ID:             info.Sub,
		FirstName:      info.GivenName,
		LastName:       info.FamilyName,
		Headline:       info.Name,
		ProfilePicture: info.Picture,
		Email:          info.Email,
	}

	return profile, nil
}

// GetByID returns a profile for the given person ID.
func (s *ProfileService) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	path := fmt.Sprintf("/people/(id:%s)", id)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get profile %s: %w", id, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get profile %s: %w", id, err)
	}

	var raw model.ProfileResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get profile %s: %w", id, err)
	}

	return raw.ToProfile(), nil
}
