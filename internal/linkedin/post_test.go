package linkedin

import (
	"context"
	"testing"

	"github.com/Softorize/lcli/internal/model"
)

func TestPostCreateSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id":          "urn:li:share:123",
			"author":      "me",
			"commentary":  "Hello",
			"visibility":  "PUBLIC",
			"createdAt":   1700000000000,
		}},
	}}

	svc := NewPostService(doer)
	post, err := svc.Create(context.Background(), &model.CreatePostRequest{
		Text:       "Hello",
		Visibility: "PUBLIC",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if post.ID != "urn:li:share:123" {
		t.Errorf("ID = %q", post.ID)
	}
	if post.Text != "Hello" {
		t.Errorf("Text = %q", post.Text)
	}
}

func TestPostCreateWithMedia(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id": "urn:li:share:456",
			"content": map[string]any{
				"media": map[string]any{"id": "urn:li:image:789"},
			},
		}},
	}}

	svc := NewPostService(doer)
	post, err := svc.Create(context.Background(), &model.CreatePostRequest{
		Text:     "With media",
		MediaURN: "urn:li:image:789",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if post.MediaCategory != "IMAGE" {
		t.Errorf("MediaCategory = %q, want IMAGE", post.MediaCategory)
	}
}

func TestPostCreateError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewPostService(doer)
	_, err := svc.Create(context.Background(), &model.CreatePostRequest{Text: "hi"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPostGetSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id":          "urn:li:share:123",
			"commentary":  "Hello",
			"visibility":  "PUBLIC",
		}},
	}}

	svc := NewPostService(doer)
	post, err := svc.Get(context.Background(), "urn:li:share:123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if post.ID != "urn:li:share:123" {
		t.Errorf("ID = %q", post.ID)
	}
}

func TestPostGetError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewPostService(doer)
	_, err := svc.Get(context.Background(), "urn:li:share:123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPostDeleteSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 204, body: nil},
	}}

	svc := NewPostService(doer)
	if err := svc.Delete(context.Background(), "urn:li:share:123"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestPostDeleteError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewPostService(doer)
	if err := svc.Delete(context.Background(), "urn"); err == nil {
		t.Fatal("expected error")
	}
}

func TestPostListByAuthorSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{
				{"id": "post1", "commentary": "Hello"},
				{"id": "post2", "commentary": "World"},
			},
			"paging": map[string]any{"count": 10, "start": 0, "total": 2},
		}},
	}}

	svc := NewPostService(doer)
	list, err := svc.ListByAuthor(context.Background(), "me", 0, 10)
	if err != nil {
		t.Fatalf("ListByAuthor: %v", err)
	}
	if len(list.Elements) != 2 {
		t.Errorf("got %d posts, want 2", len(list.Elements))
	}
	if list.Paging == nil {
		t.Fatal("paging is nil")
	}
	if list.Paging.Total != 2 {
		t.Errorf("total = %d, want 2", list.Paging.Total)
	}
}

func TestPostListByAuthorEmpty(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"elements": []map[string]any{},
		}},
	}}

	svc := NewPostService(doer)
	list, err := svc.ListByAuthor(context.Background(), "me", 0, 10)
	if err != nil {
		t.Fatalf("ListByAuthor: %v", err)
	}
	if len(list.Elements) != 0 {
		t.Errorf("got %d posts, want 0", len(list.Elements))
	}
}
