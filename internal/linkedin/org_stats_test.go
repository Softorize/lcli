package linkedin

import (
	"context"
	"testing"
)

func TestFollowerStatsSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{
					"organicFollowerCount": 800,
					"paidFollowerCount":    200,
					"followerCountsByFunction": []map[string]any{
						{"segment": "Engineering", "followerCounts": 300},
					},
					"followerCountsBySeniority": []map[string]any{
						{"segment": "Senior", "followerCounts": 150},
					},
				},
			},
		}},
	}}

	svc := NewOrgService(doer)
	stats, err := svc.FollowerStats(context.Background(), "urn:li:organization:123")
	if err != nil {
		t.Fatalf("FollowerStats: %v", err)
	}
	if stats.OrganicCount != 800 {
		t.Errorf("OrganicCount = %d", stats.OrganicCount)
	}
	if stats.PaidCount != 200 {
		t.Errorf("PaidCount = %d", stats.PaidCount)
	}
	if stats.TotalCount != 1000 {
		t.Errorf("TotalCount = %d", stats.TotalCount)
	}
	if stats.ByFunction["Engineering"] != 300 {
		t.Errorf("ByFunction[Engineering] = %d", stats.ByFunction["Engineering"])
	}
	if stats.BySeniority["Senior"] != 150 {
		t.Errorf("BySeniority[Senior] = %d", stats.BySeniority["Senior"])
	}
}

func TestFollowerStatsEmpty(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{},
		}},
	}}

	svc := NewOrgService(doer)
	stats, err := svc.FollowerStats(context.Background(), "urn")
	if err != nil {
		t.Fatalf("FollowerStats: %v", err)
	}
	if stats.TotalCount != 0 {
		t.Errorf("TotalCount = %d, want 0", stats.TotalCount)
	}
}

func TestFollowerStatsError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewOrgService(doer)
	_, err := svc.FollowerStats(context.Background(), "urn")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPageStatsSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{
					"views":          500,
					"uniqueVisitors": 300,
					"clicks":         42,
					"timeRange":      "LAST_30_DAYS",
				},
			},
		}},
	}}

	svc := NewOrgService(doer)
	stats, err := svc.PageStats(context.Background(), "urn:li:organization:123")
	if err != nil {
		t.Fatalf("PageStats: %v", err)
	}
	if stats.Views != 500 {
		t.Errorf("Views = %d", stats.Views)
	}
	if stats.UniqueVisitors != 300 {
		t.Errorf("UniqueVisitors = %d", stats.UniqueVisitors)
	}
	if stats.Clicks != 42 {
		t.Errorf("Clicks = %d", stats.Clicks)
	}
}

func TestPageStatsEmpty(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{},
		}},
	}}

	svc := NewOrgService(doer)
	stats, err := svc.PageStats(context.Background(), "urn")
	if err != nil {
		t.Fatalf("PageStats: %v", err)
	}
	if stats.Views != 0 {
		t.Errorf("Views = %d, want 0", stats.Views)
	}
}

func TestPageStatsError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 500, body: map[string]any{"status": 500, "message": "error"}},
	}}

	svc := NewOrgService(doer)
	_, err := svc.PageStats(context.Background(), "urn")
	if err == nil {
		t.Fatal("expected error")
	}
}
