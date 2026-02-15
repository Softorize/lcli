package linkedin

import (
	"context"
	"testing"
)

func TestOrgGetSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id":                   12345,
			"localizedName":       "TestCorp",
			"vanityName":          "testcorp",
			"localizedDescription": "A test company",
			"followerCount":       1000,
		}},
	}}

	svc := NewOrgService(doer)
	org, err := svc.Get(context.Background(), 12345)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if org.Name != "TestCorp" {
		t.Errorf("Name = %q", org.Name)
	}
	if org.FollowerCount != 1000 {
		t.Errorf("FollowerCount = %d", org.FollowerCount)
	}
}

func TestOrgGetError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewOrgService(doer)
	_, err := svc.Get(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOrgGetByVanitySuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{
					"id":             67890,
					"localizedName": "VanityCorp",
					"vanityName":    "vanity-corp",
				},
			},
		}},
	}}

	svc := NewOrgService(doer)
	org, err := svc.GetByVanity(context.Background(), "vanity-corp")
	if err != nil {
		t.Fatalf("GetByVanity: %v", err)
	}
	if org.Name != "VanityCorp" {
		t.Errorf("Name = %q", org.Name)
	}
}

func TestOrgGetByVanityNotFound(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{},
		}},
	}}

	svc := NewOrgService(doer)
	_, err := svc.GetByVanity(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}

func TestOrgGetByVanityAPIError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 500, body: map[string]any{"status": 500, "message": "error"}},
	}}

	svc := NewOrgService(doer)
	_, err := svc.GetByVanity(context.Background(), "corp")
	if err == nil {
		t.Fatal("expected error")
	}
}
