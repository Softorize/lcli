package command

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Softorize/lcli/internal/model"
)

func TestProfileMeSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Profile = &mockProfiler{
		meFunc: func(_ context.Context) (*model.Profile, error) {
			return &model.Profile{
				ID:        "abc123",
				FirstName: "John",
				LastName:  "Doe",
				Headline:  "Engineer",
				Vanity:    "johndoe",
				Email:     "john@example.com",
			}, nil
		},
	}

	if err := runProfileMe(nil, deps); err != nil {
		t.Fatalf("runProfileMe: %v", err)
	}

	out := stdout.String()
	for _, want := range []string{"abc123", "John Doe", "Engineer", "johndoe", "john@example.com"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestProfileMeNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runProfileMe(nil, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("error = %q, want auth error", err)
	}
}

func TestProfileMeServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Profile = &mockProfiler{
		meFunc: func(_ context.Context) (*model.Profile, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	err := runProfileMe(nil, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestProfileMeJSONOutput(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Profile = &mockProfiler{
		meFunc: func(_ context.Context) (*model.Profile, error) {
			return &model.Profile{ID: "x"}, nil
		},
	}

	if err := runProfileMe([]string{"--output", "json"}, deps); err != nil {
		t.Fatalf("runProfileMe json: %v", err)
	}

	if !strings.Contains(stdout.String(), `"id"`) {
		t.Errorf("JSON output missing 'id' key:\n%s", stdout.String())
	}
}

func TestProfileViewSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Profile = &mockProfiler{
		getByIDFunc: func(_ context.Context, id string) (*model.Profile, error) {
			return &model.Profile{
				ID:        id,
				FirstName: "Jane",
				LastName:  "Smith",
			}, nil
		},
	}

	if err := runProfileView([]string{"--id", "person123"}, deps); err != nil {
		t.Fatalf("runProfileView: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "person123") {
		t.Errorf("output missing profile ID:\n%s", out)
	}
}

func TestProfileViewMissingID(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Profile = &mockProfiler{}

	err := runProfileView(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --id")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("error = %q, want --id required", err)
	}
}

func TestProfileViewNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runProfileView([]string{"--id", "x"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestProfileDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runProfile(nil, deps); err != nil {
		t.Fatalf("runProfile help: %v", err)
	}
	if !strings.Contains(stdout.String(), "me") {
		t.Error("help output missing 'me' subcommand")
	}
}

func TestProfileDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runProfile([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
