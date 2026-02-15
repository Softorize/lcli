package linkedin

import (
	"context"
	"testing"

	"github.com/Softorize/lcli/internal/model"
)

func TestReactSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 201, body: nil},
	}}

	svc := NewReactionService(doer)
	err := svc.React(context.Background(), "me", "urn:li:share:123", model.ReactionLike)
	if err != nil {
		t.Fatalf("React: %v", err)
	}

	if len(doer.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(doer.calls))
	}
	if doer.calls[0].method != "POST" {
		t.Errorf("method = %q, want POST", doer.calls[0].method)
	}
}

func TestReactError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 409, body: map[string]any{"status": 409, "message": "conflict"}},
	}}

	svc := NewReactionService(doer)
	err := svc.React(context.Background(), "me", "urn", model.ReactionLike)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnreactSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 204, body: nil},
	}}

	svc := NewReactionService(doer)
	err := svc.Unreact(context.Background(), "me", "urn:li:share:123")
	if err != nil {
		t.Fatalf("Unreact: %v", err)
	}

	if doer.calls[0].method != "DELETE" {
		t.Errorf("method = %q, want DELETE", doer.calls[0].method)
	}
}

func TestUnreactError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewReactionService(doer)
	err := svc.Unreact(context.Background(), "me", "urn")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReactionListSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{"actor": "a1", "reactionType": "LIKE", "created": 1700000000000},
				{"actor": "a2", "reactionType": "CELEBRATE", "created": 1700000000000},
			},
			"paging": map[string]any{"count": 10, "start": 0, "total": 2},
		}},
	}}

	svc := NewReactionService(doer)
	list, err := svc.List(context.Background(), "urn:li:share:123", 0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list.Elements) != 2 {
		t.Errorf("got %d reactions, want 2", len(list.Elements))
	}
	if list.Elements[0].Type != model.ReactionLike {
		t.Errorf("type = %q, want LIKE", list.Elements[0].Type)
	}
}

func TestReactionListError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 500, body: map[string]any{"status": 500, "message": "error"}},
	}}

	svc := NewReactionService(doer)
	_, err := svc.List(context.Background(), "urn", 0, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}
