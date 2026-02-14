package linkedin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AnalyticsService provides access to LinkedIn analytics endpoints.
type AnalyticsService struct {
	doer Doer
}

// NewAnalyticsService creates an AnalyticsService backed by the given Doer.
func NewAnalyticsService(d Doer) *AnalyticsService {
	return &AnalyticsService{doer: d}
}

// PostAnalytics retrieves engagement metrics for a single post.
// The returned map contains keys like "impressionCount", "clickCount",
// "likeCount", "commentCount", "shareCount", and "engagementRate".
func (s *AnalyticsService) PostAnalytics(ctx context.Context, postURN string) (map[string]any, error) {
	path := fmt.Sprintf(
		"/organizationalEntityShareStatistics?q=organizationalEntity&shares[0]=%s",
		url.QueryEscape(postURN),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("post analytics for %s: %w", postURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("post analytics for %s: %w", postURN, err)
	}

	var raw struct {
		Elements []struct {
			TotalShareStatistics map[string]any `json:"totalShareStatistics"`
		} `json:"elements"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("post analytics for %s: %w", postURN, err)
	}

	if len(raw.Elements) == 0 {
		return map[string]any{}, nil
	}

	return raw.Elements[0].TotalShareStatistics, nil
}

// ProfileViews retrieves the number of profile views for the
// authenticated user.
func (s *AnalyticsService) ProfileViews(ctx context.Context) (int, error) {
	resp, err := s.doer.Do(ctx, http.MethodGet, "/networkSizes/me?edgeType=CompanyFollowedByMember", nil)
	if err != nil {
		return 0, fmt.Errorf("profile views: %w", err)
	}

	if err := checkError(resp); err != nil {
		return 0, fmt.Errorf("profile views: %w", err)
	}

	var raw struct {
		FirstDegreeSize int `json:"firstDegreeSize"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return 0, fmt.Errorf("profile views: %w", err)
	}

	return raw.FirstDegreeSize, nil
}
