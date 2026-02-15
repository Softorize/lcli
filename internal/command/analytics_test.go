package command

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestAnalyticsPostSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Analytics = &mockAnalyticsReader{
		postAnalyticsFunc: func(_ context.Context, urn string) (map[string]any, error) {
			return map[string]any{
				"impressionCount": 1000,
				"likeCount":       50,
			}, nil
		},
	}

	if err := runAnalyticsPost([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runAnalyticsPost: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "1000") {
		t.Errorf("output missing impression count:\n%s", out)
	}
}

func TestAnalyticsPostMissingURN(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Analytics = &mockAnalyticsReader{}

	err := runAnalyticsPost(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing URN")
	}
	if !strings.Contains(err.Error(), "post URN argument is required") {
		t.Errorf("error = %q", err)
	}
}

func TestAnalyticsPostNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runAnalyticsPost([]string{"urn"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestAnalyticsPostServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Analytics = &mockAnalyticsReader{
		postAnalyticsFunc: func(_ context.Context, _ string) (map[string]any, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	err := runAnalyticsPost([]string{"urn"}, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalyticsViewsSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Analytics = &mockAnalyticsReader{
		profileViewsFunc: func(_ context.Context) (int, error) {
			return 42, nil
		},
	}

	if err := runAnalyticsViews(nil, deps); err != nil {
		t.Fatalf("runAnalyticsViews: %v", err)
	}

	if !strings.Contains(stdout.String(), "42") {
		t.Errorf("output missing view count:\n%s", stdout.String())
	}
}

func TestAnalyticsViewsNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runAnalyticsViews(nil, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestAnalyticsViewsServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Analytics = &mockAnalyticsReader{
		profileViewsFunc: func(_ context.Context) (int, error) {
			return 0, fmt.Errorf("api error")
		},
	}

	err := runAnalyticsViews(nil, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalyticsDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runAnalytics(nil, deps); err != nil {
		t.Fatalf("runAnalytics help: %v", err)
	}
	if !strings.Contains(stdout.String(), "post") {
		t.Error("help output missing 'post' subcommand")
	}
}

func TestAnalyticsDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runAnalytics([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
