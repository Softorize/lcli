package command

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Softorize/lcli/internal/model"
)

func TestPostCreateSuccess(t *testing.T) {
	deps, _, stderr := testDeps()
	deps.Posts = &mockPoster{
		createFunc: func(_ context.Context, req *model.CreatePostRequest) (*model.Post, error) {
			if req.Text != "Hello world" {
				t.Errorf("text = %q, want 'Hello world'", req.Text)
			}
			if req.Visibility != "PUBLIC" {
				t.Errorf("visibility = %q, want 'PUBLIC'", req.Visibility)
			}
			return &model.Post{ID: "urn:li:share:123"}, nil
		},
	}

	if err := runPostCreate([]string{"--text", "Hello world"}, deps); err != nil {
		t.Fatalf("runPostCreate: %v", err)
	}

	if !strings.Contains(stderr.String(), "urn:li:share:123") {
		t.Errorf("stderr missing post ID:\n%s", stderr.String())
	}
}

func TestPostCreateMissingText(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Posts = &mockPoster{}

	err := runPostCreate(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --text")
	}
	if !strings.Contains(err.Error(), "--text is required") {
		t.Errorf("error = %q", err)
	}
}

func TestPostCreateInvalidVisibility(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Posts = &mockPoster{}

	err := runPostCreate([]string{"--text", "hi", "--visibility", "INVALID"}, deps)
	if err == nil {
		t.Fatal("expected error for invalid visibility")
	}
	if !strings.Contains(err.Error(), "invalid visibility") {
		t.Errorf("error = %q", err)
	}
}

func TestPostCreateNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runPostCreate([]string{"--text", "hi"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestPostListSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Posts = &mockPoster{
		listByAuthorFunc: func(_ context.Context, _ string, _, _ int) (*model.PostList, error) {
			return &model.PostList{
				Elements: []model.Post{
					{ID: "post1", Text: "Hello", Visibility: "PUBLIC", CreatedAt: time.Now()},
					{ID: "post2", Text: "World", Visibility: "CONNECTIONS", CreatedAt: time.Now()},
				},
			}, nil
		},
	}

	if err := runPostList(nil, deps); err != nil {
		t.Fatalf("runPostList: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "post1") || !strings.Contains(out, "post2") {
		t.Errorf("output missing posts:\n%s", out)
	}
}

func TestPostListNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runPostList(nil, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestPostGetSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Posts = &mockPoster{
		getFunc: func(_ context.Context, urn string) (*model.Post, error) {
			return &model.Post{
				ID:     urn,
				Author: "me",
				Text:   "Test post",
			}, nil
		},
	}

	if err := runPostGet([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runPostGet: %v", err)
	}

	if !strings.Contains(stdout.String(), "urn:li:share:123") {
		t.Errorf("output missing URN:\n%s", stdout.String())
	}
}

func TestPostGetMissingURN(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Posts = &mockPoster{}

	err := runPostGet(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing URN")
	}
}

func TestPostDeleteWithConfirm(t *testing.T) {
	deps, _, stderr := testDeps()
	deleted := false
	deps.Posts = &mockPoster{
		deleteFunc: func(_ context.Context, urn string) error {
			deleted = true
			return nil
		},
	}

	if err := runPostDelete([]string{"--confirm", "urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runPostDelete: %v", err)
	}

	if !deleted {
		t.Error("post was not deleted")
	}
	if !strings.Contains(stderr.String(), "deleted") {
		t.Error("stderr missing delete confirmation")
	}
}

func TestPostDeleteWithoutConfirm(t *testing.T) {
	deps, _, stderr := testDeps()
	deps.Posts = &mockPoster{}

	if err := runPostDelete([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runPostDelete: %v", err)
	}

	if !strings.Contains(stderr.String(), "--confirm") {
		t.Error("stderr missing --confirm hint")
	}
}

func TestPostDeleteServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Posts = &mockPoster{
		deleteFunc: func(_ context.Context, _ string) error {
			return fmt.Errorf("forbidden")
		},
	}

	err := runPostDelete([]string{"--confirm", "urn:li:share:123"}, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPostDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runPost(nil, deps); err != nil {
		t.Fatalf("runPost help: %v", err)
	}
	if !strings.Contains(stdout.String(), "create") {
		t.Error("help output missing 'create' subcommand")
	}
}

func TestPostDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runPost([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestValidateVisibility(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"PUBLIC", true},
		{"CONNECTIONS", true},
		{"INVALID", false},
		{"", false},
	}
	for _, tt := range tests {
		err := validateVisibility(tt.input)
		if tt.ok && err != nil {
			t.Errorf("validateVisibility(%q) = %v, want nil", tt.input, err)
		}
		if !tt.ok && err == nil {
			t.Errorf("validateVisibility(%q) = nil, want error", tt.input)
		}
	}
}
