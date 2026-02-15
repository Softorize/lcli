package linkedin

import (
	"context"
	"testing"

	"github.com/Softorize/lcli/internal/model"
)

func TestCommentCreateSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"$URN":    "urn:li:comment:456",
			"actor":   "me",
			"created": 1700000000000,
			"message": map[string]any{"text": "Nice!"},
		}},
	}}

	svc := NewCommentService(doer)
	comment, err := svc.Create(context.Background(), &model.CreateCommentRequest{
		PostURN: "urn:li:share:123",
		Text:    "Nice!",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if comment.ID != "urn:li:comment:456" {
		t.Errorf("ID = %q", comment.ID)
	}
	if comment.Text != "Nice!" {
		t.Errorf("Text = %q", comment.Text)
	}
}

func TestCommentCreateError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewCommentService(doer)
	_, err := svc.Create(context.Background(), &model.CreateCommentRequest{
		PostURN: "urn", Text: "hi",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCommentListSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{"$URN": "c1", "actor": "a1", "message": map[string]any{"text": "Hello"}},
				{"$URN": "c2", "actor": "a2", "message": map[string]any{"text": "World"}},
			},
			"paging": map[string]any{"count": 10, "start": 0, "total": 2},
		}},
	}}

	svc := NewCommentService(doer)
	list, err := svc.List(context.Background(), "urn:li:share:123", 0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list.Elements) != 2 {
		t.Errorf("got %d comments, want 2", len(list.Elements))
	}
}

func TestCommentListError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 500, body: map[string]any{"status": 500, "message": "error"}},
	}}

	svc := NewCommentService(doer)
	_, err := svc.List(context.Background(), "urn", 0, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCommentDeleteSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 204, body: nil},
	}}

	svc := NewCommentService(doer)
	if err := svc.Delete(context.Background(), "urn:li:comment:456"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestCommentDeleteError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewCommentService(doer)
	if err := svc.Delete(context.Background(), "urn"); err == nil {
		t.Fatal("expected error")
	}
}
