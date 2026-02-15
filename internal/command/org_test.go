package command

import (
	"context"
	"strings"
	"testing"

	"github.com/Softorize/lcli/internal/model"
)

func TestOrgInfoByID(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Orgs = &mockOrgReader{
		getFunc: func(_ context.Context, id int64) (*model.Organization, error) {
			return &model.Organization{
				ID:   id,
				Name: "TestCorp",
			}, nil
		},
	}

	if err := runOrgInfo([]string{"--id", "12345"}, deps); err != nil {
		t.Fatalf("runOrgInfo: %v", err)
	}

	if !strings.Contains(stdout.String(), "TestCorp") {
		t.Errorf("output missing org name:\n%s", stdout.String())
	}
}

func TestOrgInfoByVanity(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Orgs = &mockOrgReader{
		getByVanityFunc: func(_ context.Context, name string) (*model.Organization, error) {
			return &model.Organization{
				Name:       "VanityCorp",
				VanityName: name,
			}, nil
		},
	}

	if err := runOrgInfo([]string{"--vanity", "vanity-corp"}, deps); err != nil {
		t.Fatalf("runOrgInfo: %v", err)
	}

	if !strings.Contains(stdout.String(), "VanityCorp") {
		t.Errorf("output missing org name:\n%s", stdout.String())
	}
}

func TestOrgInfoMissingIDAndVanity(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Orgs = &mockOrgReader{}

	err := runOrgInfo(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --id/--vanity")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Errorf("error = %q", err)
	}
}

func TestOrgInfoInvalidID(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Orgs = &mockOrgReader{}

	err := runOrgInfo([]string{"--id", "notanumber"}, deps)
	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}

func TestOrgInfoNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runOrgInfo([]string{"--id", "123"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestOrgFollowersSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Orgs = &mockOrgReader{
		followerStatsFunc: func(_ context.Context, _ string) (*model.OrgFollowerStats, error) {
			return &model.OrgFollowerStats{
				OrganicCount: 100,
				PaidCount:    50,
				TotalCount:   150,
			}, nil
		},
	}

	if err := runOrgFollowers([]string{"--org", "urn:li:organization:123"}, deps); err != nil {
		t.Fatalf("runOrgFollowers: %v", err)
	}

	if !strings.Contains(stdout.String(), "150") {
		t.Errorf("output missing total count:\n%s", stdout.String())
	}
}

func TestOrgFollowersMissingOrg(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Orgs = &mockOrgReader{}

	err := runOrgFollowers(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --org")
	}
}

func TestOrgStatsSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Orgs = &mockOrgReader{
		pageStatsFunc: func(_ context.Context, _ string) (*model.OrgPageStats, error) {
			return &model.OrgPageStats{
				Views:          500,
				UniqueVisitors: 300,
				Clicks:         42,
			}, nil
		},
	}

	if err := runOrgStats([]string{"--org", "urn:li:organization:123"}, deps); err != nil {
		t.Fatalf("runOrgStats: %v", err)
	}

	if !strings.Contains(stdout.String(), "500") {
		t.Errorf("output missing views:\n%s", stdout.String())
	}
}

func TestOrgStatsMissingOrg(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Orgs = &mockOrgReader{}

	err := runOrgStats(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --org")
	}
}

func TestOrgDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runOrg(nil, deps); err != nil {
		t.Fatalf("runOrg help: %v", err)
	}
	if !strings.Contains(stdout.String(), "info") {
		t.Error("help output missing 'info' subcommand")
	}
}

func TestOrgDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runOrg([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
