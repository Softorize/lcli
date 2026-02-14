package linkedin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Softorize/lcli/internal/model"
)

// ProfileService provides access to LinkedIn profile endpoints.
type ProfileService struct {
	doer Doer
}

// NewProfileService creates a ProfileService backed by the given Doer.
func NewProfileService(d Doer) *ProfileService {
	return &ProfileService{doer: d}
}

// Me returns the authenticated user's profile.
func (s *ProfileService) Me(ctx context.Context) (*model.Profile, error) {
	path := "/me?projection=(id,localizedFirstName,localizedLastName,localizedHeadline,vanityName,profilePicture(displayImage~:playableStreams))"

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get my profile: %w", err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get my profile: %w", err)
	}

	var raw model.ProfileResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get my profile: %w", err)
	}

	profile := raw.ToProfile()

	email, err := s.fetchEmail(ctx)
	if err == nil {
		profile.Email = email
	}

	return profile, nil
}

// GetByID returns a profile for the given person ID.
func (s *ProfileService) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	path := fmt.Sprintf("/people/(id:%s)?projection=(id,localizedFirstName,localizedLastName,localizedHeadline,vanityName)", id)

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

// emailResponse is the structure returned by the /emailAddress endpoint.
type emailResponse struct {
	Elements []emailElement `json:"elements"`
}

// emailElement holds a single email entry from the API.
type emailElement struct {
	Handle      emailHandle `json:"handle~"`
	HandleTilde emailHandle `json:"handle"`
}

// emailHandle contains the actual email address string.
type emailHandle struct {
	EmailAddress string `json:"emailAddress"`
}

// fetchEmail retrieves the primary email for the authenticated user.
func (s *ProfileService) fetchEmail(ctx context.Context) (string, error) {
	path := "/emailAddress?q=members&projection=(elements*(handle~))"

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("fetch email: %w", err)
	}

	if err := checkError(resp); err != nil {
		return "", fmt.Errorf("fetch email: %w", err)
	}

	var result emailResponse
	if err := decodeJSON(resp, &result); err != nil {
		return "", fmt.Errorf("fetch email: %w", err)
	}

	if len(result.Elements) == 0 {
		return "", nil
	}

	return result.Elements[0].Handle.EmailAddress, nil
}
