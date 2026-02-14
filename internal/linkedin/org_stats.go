package linkedin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/toto/lcli/internal/model"
)

// FollowerStats retrieves follower statistics for an organization.
func (s *OrgService) FollowerStats(ctx context.Context, orgURN string) (*model.OrgFollowerStats, error) {
	path := fmt.Sprintf(
		"/organizationalEntityFollowerStatistics?q=organizationalEntity&organizationalEntity=%s",
		url.QueryEscape(orgURN),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("follower stats for %s: %w", orgURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("follower stats for %s: %w", orgURN, err)
	}

	var raw struct {
		Elements []followerStatsResponse `json:"elements"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("follower stats for %s: %w", orgURN, err)
	}

	if len(raw.Elements) == 0 {
		return &model.OrgFollowerStats{}, nil
	}

	return raw.Elements[0].toStats(), nil
}

// followerStatsResponse maps the raw follower statistics from the API.
type followerStatsResponse struct {
	OrganicCount int               `json:"organicFollowerCount"`
	PaidCount    int               `json:"paidFollowerCount"`
	ByFunction   []followerSegment `json:"followerCountsByFunction"`
	BySeniority  []followerSegment `json:"followerCountsBySeniority"`
}

// followerSegment holds a single segment of follower counts.
type followerSegment struct {
	Segment string `json:"segment"`
	Count   int    `json:"followerCounts"`
}

// toStats converts the raw response into a domain OrgFollowerStats.
func (r *followerStatsResponse) toStats() *model.OrgFollowerStats {
	stats := &model.OrgFollowerStats{
		OrganicCount: r.OrganicCount,
		PaidCount:    r.PaidCount,
		TotalCount:   r.OrganicCount + r.PaidCount,
		ByFunction:   make(map[string]int, len(r.ByFunction)),
		BySeniority:  make(map[string]int, len(r.BySeniority)),
	}

	for _, seg := range r.ByFunction {
		stats.ByFunction[seg.Segment] = seg.Count
	}
	for _, seg := range r.BySeniority {
		stats.BySeniority[seg.Segment] = seg.Count
	}

	return stats
}

// PageStats retrieves page view statistics for an organization.
func (s *OrgService) PageStats(ctx context.Context, orgURN string) (*model.OrgPageStats, error) {
	path := fmt.Sprintf(
		"/organizationPageStatistics?q=organization&organization=%s",
		url.QueryEscape(orgURN),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("page stats for %s: %w", orgURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("page stats for %s: %w", orgURN, err)
	}

	var raw struct {
		Elements []pageStatsResponse `json:"elements"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("page stats for %s: %w", orgURN, err)
	}

	if len(raw.Elements) == 0 {
		return &model.OrgPageStats{}, nil
	}

	return raw.Elements[0].toPageStats(), nil
}

// pageStatsResponse maps the raw page statistics from the API.
type pageStatsResponse struct {
	Views          int    `json:"views"`
	UniqueVisitors int    `json:"uniqueVisitors"`
	Clicks         int    `json:"clicks"`
	Period         string `json:"timeRange"`
}

// toPageStats converts the raw response into a domain OrgPageStats.
func (r *pageStatsResponse) toPageStats() *model.OrgPageStats {
	return &model.OrgPageStats{
		Views:          r.Views,
		UniqueVisitors: r.UniqueVisitors,
		Clicks:         r.Clicks,
		Period:         r.Period,
	}
}
