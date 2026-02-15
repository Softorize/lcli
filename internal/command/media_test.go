package command

import (
	"strings"
	"testing"
)

func TestDetectMediaType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"photo.jpg", "image"},
		{"photo.jpeg", "image"},
		{"photo.png", "image"},
		{"photo.gif", "image"},
		{"photo.webp", "image"},
		{"photo.JPG", "image"},
		{"clip.mp4", "video"},
		{"clip.mov", "video"},
		{"clip.avi", "video"},
		{"clip.wmv", "video"},
		{"clip.webm", "video"},
		{"file.txt", ""},
		{"noext", ""},
	}
	for _, tt := range tests {
		got := detectMediaType(tt.path)
		if got != tt.want {
			t.Errorf("detectMediaType(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestMediaUploadMissingFile(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Media = &mockMediaUploader{}

	err := runMediaUpload(nil, deps)
	if err == nil {
		t.Fatal("expected error for missing file arg")
	}
	if !strings.Contains(err.Error(), "file path argument is required") {
		t.Errorf("error = %q", err)
	}
}

func TestMediaUploadUnknownType(t *testing.T) {
	deps, _, _ := testDeps()
	deps.Media = &mockMediaUploader{}

	err := runMediaUpload([]string{"file.txt"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown media type")
	}
	if !strings.Contains(err.Error(), "unable to detect type") {
		t.Errorf("error = %q", err)
	}
}

func TestMediaUploadNotAuthenticated(t *testing.T) {
	deps, _, _ := testDeps()

	err := runMediaUpload([]string{"photo.jpg"}, deps)
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestMediaDispatchHelp(t *testing.T) {
	deps, stdout, _ := testDeps()

	if err := runMedia(nil, deps); err != nil {
		t.Fatalf("runMedia help: %v", err)
	}
	if !strings.Contains(stdout.String(), "upload") {
		t.Error("help output missing 'upload' subcommand")
	}
}

func TestMediaDispatchUnknown(t *testing.T) {
	deps, _, _ := testDeps()

	err := runMedia([]string{"unknown"}, deps)
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}
