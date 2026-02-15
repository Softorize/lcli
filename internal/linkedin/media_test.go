package linkedin

import (
	"context"
	"testing"
)

func TestInitUploadImage(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"value": map[string]any{
				"uploadUrl": "https://upload.example.com/123",
				"image":     "urn:li:image:abc",
			},
		}},
	}}

	svc := NewMediaService(doer)
	upload, err := svc.InitUpload(context.Background(), "me", "IMAGE")
	if err != nil {
		t.Fatalf("InitUpload: %v", err)
	}

	if upload.UploadURL != "https://upload.example.com/123" {
		t.Errorf("UploadURL = %q", upload.UploadURL)
	}
	if upload.MediaURN != "urn:li:image:abc" {
		t.Errorf("MediaURN = %q", upload.MediaURN)
	}
}

func TestInitUploadVideo(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"value": map[string]any{
				"uploadUrl":   "https://upload.example.com/456",
				"video":       "urn:li:video:def",
				"uploadToken": "tok123",
			},
		}},
	}}

	svc := NewMediaService(doer)
	upload, err := svc.InitUpload(context.Background(), "me", "VIDEO")
	if err != nil {
		t.Fatalf("InitUpload: %v", err)
	}

	if upload.MediaURN != "urn:li:video:def" {
		t.Errorf("MediaURN = %q", upload.MediaURN)
	}
	if upload.UploadToken != "tok123" {
		t.Errorf("UploadToken = %q", upload.UploadToken)
	}
}

func TestInitUploadUnsupportedType(t *testing.T) {
	doer := &mockDoer{}

	svc := NewMediaService(doer)
	_, err := svc.InitUpload(context.Background(), "me", "AUDIO")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestInitUploadAPIError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 403, body: map[string]any{"status": 403, "message": "forbidden"}},
	}}

	svc := NewMediaService(doer)
	_, err := svc.InitUpload(context.Background(), "me", "IMAGE")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetStatusSuccess(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 200, body: map[string]any{
			"id":     "urn:li:image:abc",
			"status": "AVAILABLE",
		}},
	}}

	svc := NewMediaService(doer)
	status, err := svc.GetStatus(context.Background(), "urn:li:image:abc")
	if err != nil {
		t.Fatalf("GetStatus: %v", err)
	}
	if status.Status != "AVAILABLE" {
		t.Errorf("Status = %q", status.Status)
	}
}

func TestGetStatusError(t *testing.T) {
	doer := &mockDoer{responses: []mockResponse{
		{status: 404, body: map[string]any{"status": 404, "message": "not found"}},
	}}

	svc := NewMediaService(doer)
	_, err := svc.GetStatus(context.Background(), "urn")
	if err == nil {
		t.Fatal("expected error")
	}
}
