package linkedin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/toto/lcli/internal/model"
)

// OrgService provides access to LinkedIn organization endpoints.
type OrgService struct {
	doer Doer
}

// NewOrgService creates an OrgService backed by the given Doer.
func NewOrgService(d Doer) *OrgService {
	return &OrgService{doer: d}
}

// orgResponse is the raw API response for an organization.
type orgResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"localizedName"`
	VanityName    string `json:"vanityName"`
	Description   string `json:"localizedDescription"`
	LogoURL       string `json:"logoV2"`
	WebsiteURL    string `json:"localizedWebsite"`
	FollowerCount int    `json:"followerCount"`
}

// toOrg converts a raw API response into a domain Organization.
func (r *orgResponse) toOrg() *model.Organization {
	return &model.Organization{
		ID:            r.ID,
		Name:          r.Name,
		VanityName:    r.VanityName,
		Description:   r.Description,
		LogoURL:       r.LogoURL,
		Website:       r.WebsiteURL,
		FollowerCount: r.FollowerCount,
	}
}

// Get retrieves an organization by its numeric ID.
func (s *OrgService) Get(ctx context.Context, id int64) (*model.Organization, error) {
	path := fmt.Sprintf("/organizations/%s", strconv.FormatInt(id, 10))

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get org %d: %w", id, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get org %d: %w", id, err)
	}

	var raw orgResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get org %d: %w", id, err)
	}

	return raw.toOrg(), nil
}

// GetByVanity retrieves an organization by its vanity name (URL slug).
func (s *OrgService) GetByVanity(ctx context.Context, vanityName string) (*model.Organization, error) {
	path := fmt.Sprintf("/organizations?q=vanityName&vanityName=%s", url.QueryEscape(vanityName))

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get org by vanity %s: %w", vanityName, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get org by vanity %s: %w", vanityName, err)
	}

	var raw struct {
		Elements []orgResponse `json:"elements"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get org by vanity %s: %w", vanityName, err)
	}

	if len(raw.Elements) == 0 {
		return nil, fmt.Errorf("get org by vanity %s: %w", vanityName, model.ErrNotFound)
	}

	return raw.Elements[0].toOrg(), nil
}
