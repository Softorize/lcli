package command

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Softorize/lcli/internal/model"
)

func TestReactionLikeSuccess(t *testing.T) {
	deps, _, stderr := testDeps()
	deps.Reactions = &mockReacter{
		reactFunc: func(_ context.Context, actor, entity string, rt model.ReactionType) error {
			if rt != model.ReactionLike {
				t.Errorf("reaction type = %q, want LIKE", rt)
			}
			return nil
		},
	}

	if err := runReactionLike([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runReactionLike: %v", err)
	}

	if !strings.Contains(stderr.String(), "Reacted") {
		t.Error("stderr missing confirmation")
	}
}

func TestReactionLikeCustomType(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{
		reactFunc: func(_ context.Context, _, _ string, rt model.ReactionType) error {
			if rt != model.ReactionCelebrate {
				t.Errorf("reaction type = %q, want CELEBRATE", rt)
			}
			return nil
		},
	}

	if err := runReactionLike([]string{"--type", "CELEBRATE", "urn"}, deps); err != nil {
		t.Fatalf("runReactionLike: %v", err)
	}
}

func TestReactionLikeMissingURN(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{}

	err := runReactionLike(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing URN")
	}
}

func TestReactionLikeInvalidType(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{}

	err := runReactionLike([]string{"--type", "INVALID", "urn"}, deps)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestReactionLikeNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runReactionLike([]string{"urn"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestReactionUnlikeSuccess(t *testing.T) {
	deps, _, stderr := testDeps()
	deps.Reactions = &mockReacter{
		unreactFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}

	if err := runReactionUnlike([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runReactionUnlike: %v", err)
	}

	if !strings.Contains(stderr.String(), "removed") {
		t.Error("stderr missing confirmation")
	}
}

func TestReactionUnlikeMissingURN(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{}

	err := runReactionUnlike(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing URN")
	}
}

func TestReactionUnlikeServiceError(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{
		unreactFunc: func(_ context.Context, _, _ string) error {
			return fmt.Errorf("not found")
		},
	}

	err := runReactionUnlike([]string{"urn"}, deps)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReactionListSuccess(t *testing.T) {
	deps, stdout, _ := testDeps()
	deps.Reactions = &mockReacter{
		listFunc: func(_ context.Context, _ string, _, _ int) (*model.ReactionList, error) {
			return &model.ReactionList{
				Elements: []model.Reaction{
					{Actor: "actor1", Type: model.ReactionLike, CreatedAt: time.Now()},
				},
			}, nil
		},
	}

	if err := runReactionList([]string{"urn:li:share:123"}, deps); err != nil {
		t.Fatalf("runReactionList: %v", err)
	}

	if !strings.Contains(stdout.String(), "actor1") {
		t.Errorf("output missing reaction:\n%s", stdout.String())
	}
}

func TestReactionListMissingURN(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Reactions = &mockReacter{}

	err := runReactionList(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing URN")
	}
}

func TestReactionDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runReaction(nil, deps); err != nil {
		t.Fatalf("runReaction help: %v", err)
	}
	if !strings.Contains(stdout.String(), "like") {
		t.Error("help output missing 'like' subcommand")
	}
}

func TestReactionDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runReaction([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
