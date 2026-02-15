package linkedin

import (
	"context"
	"testing"
)

func TestPostAnalyticsSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{
					"totalShareStatistics": map[string]any{
						"impressionCount": 1000,
						"clickCount":      50,
						"likeCount":       30,
						"commentCount":    10,
						"shareCount":      5,
					},
				},
			},
		}},
	}}

	svc := NewAnalyticsService(doer)
	stats, err := svc.PostAnalytics(context.Background(), "urn:li:share:123")
	if err != nil {
		t.Fatalf("PostAnalytics: %v", err)
	}

	if stats["impressionCount"] != float64(1000) {
		t.Errorf("impressionCount = %v", stats["impressionCount"])
	}
	if stats["clickCount"] != float64(50) {
		t.Errorf("clickCount = %v", stats["clickCount"])
	}
}

func TestPostAnalyticsEmpty(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{},
		}},
	}}

	svc := NewAnalyticsService(doer)
	stats, err := svc.PostAnalytics(context.Background(), "urn")
	if err != nil {
		t.Fatalf("PostAnalytics: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("expected empty stats, got %v", stats)
	}
}

func TestPostAnalyticsError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewAnalyticsService(doer)
	_, err := svc.PostAnalytics(context.Background(), "urn")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProfileViewsSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"firstDegreeSize": 42,
		}},
	}}

	svc := NewAnalyticsService(doer)
	count, err := svc.ProfileViews(context.Background())
	if err != nil {
		t.Fatalf("ProfileViews: %v", err)
	}
	if count != 42 {
		t.Errorf("count = %d, want 42", count)
	}
}

func TestProfileViewsError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 401, body: map[string]any{"status": 401, "message": "unauthorized"}},
	}}

	svc := NewAnalyticsService(doer)
	_, err := svc.ProfileViews(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
