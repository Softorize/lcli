package command

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Softorize/lcli/internal/model"
)

func TestCommentCreateSuccess(t *testing.T) {
	deps, _, stderr := testDeps()
	deps.Comments = &mockCommenter{
		createFunc: func(_ context.Context, req *model.CreateCommentRequest) (*model.Comment, error) {
			if req.PostURN != "urn:li:share:123" {
				t.Errorf("postURN = %q", req.PostURN)
			}
			if req.Text != "Nice post!" {
				t.Errorf("text = %q", req.Text)
			}
			return &model.Comment{ID: "urn:li:comment:456"}, nil
		},
	}

	if err := runCommentCreate([]string{"--post", "urn:li:share:123", "--text", "Nice post!"}, deps); err != nil {
		t.Fatalf("runCommentCreate: %v", err)
	}

	if !strings.Contains(stderr.String(), "urn:li:comment:456") {
		t.Errorf("stderr missing comment ID:\n%s", stderr.String())
	}
}

func TestCommentCreateMissingPost(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Comments = &mockCommenter{}

	err := runCommentCreate([]string{"--text", "hi"}, deps)
	if err == nil {
		t.Fatal("expected error for missing --post")
	}
	if !strings.Contains(err.Error(), "--post is required") {
		t.Errorf("error = %q", err)
	}
}

func TestCommentCreateMissingText(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Comments = &mockCommenter{}

	err := runCommentCreate([]string{"--post", "urn"}, deps)
	if err == nil {
		t.Fatal("expected error for missing --text")
	}
}

func TestCommentCreateNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runCommentCreate([]string{"--post", "urn", "--text", "hi"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestCommentListSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Comments = &mockCommenter{
		listFunc: func(_ context.Context, _ string, _, _ int) (*model.CommentList, error) {
			return &model.CommentList{
				Elements: []model.Comment{
					{ID: "c1", Author: "actor1", Text: "Hello", CreatedAt: time.Now()},
				},
			}, nil
		},
	}

	if err := runCommentList([]string{"--post", "urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runCommentList: %v", err)
	}

	if !strings.Contains(stdout.String(), "c1") {
		t.Errorf("output missing comment:\n%s", stdout.String())
	}
}

func TestCommentListMissingPost(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Comments = &mockCommenter{}

	err := runCommentList(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing --post")
	}
}

func TestCommentDeleteWithConfirm(t *testing.T) {
	deps, _, stderr := testDeps()
	deleted := false
	deps.Comments = &mockCommenter{
		deleteFunc: func(_ context.Context, _ string) error {
			deleted = true
			return nil
		},
	}

	if err := runCommentDelete([]string{"--confirm", "urn:li:comment:456"}, deps); err != nil {
		t.Fatalf("runCommentDelete: %v", err)
	}

	if !deleted {
		t.Error("comment was not deleted")
	}
	if !strings.Contains(stderr.String(), "deleted") {
		t.Error("stderr missing delete confirmation")
	}
}

func TestCommentDeleteWithoutConfirm(t *testing.T) {
	deps, _, stderr := testDeps()

	if err := runCommentDelete([]string{"urn:li:comment:456"}, deps); err != nil {
		t.Fatalf("runCommentDelete: %v", err)
	}

	if !strings.Contains(stderr.String(), "--confirm") {
		t.Error("stderr missing --confirm hint")
	}
}

func TestCommentDeleteServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Comments = &mockCommenter{
		deleteFunc: func(_ context.Context, _ string) error {
			return fmt.Errorf("not found")
		},
	}

	err := runCommentDelete([]string{"--confirm", "urn"}, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCommentDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runComment(nil, deps); err != nil {
		t.Fatalf("runComment help: %v", err)
	}
	if !strings.Contains(stdout.String(), "create") {
		t.Error("help output missing 'create' subcommand")
	}
}

func TestCommentDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runComment([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
