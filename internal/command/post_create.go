package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/toto/lcli/internal/model"
)

// runPostCreate handles the post create subcommand.
func runPostCreate(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("post create", flag.ContinueOnError)
	text := fs.String("text", "", "Post text content (required)")
	visibility := fs.String("visibility", "PUBLIC", "Visibility: PUBLIC or CONNECTIONS")
	image := fs.String("image", "", "Path to image file to attach")
	video := fs.String("video", "", "Path to video file to attach")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *text == "" {
		return fmt.Errorf("post create: --text is required")
	}
	if err := validateVisibility(*visibility); err != nil {
		return err
	}
	if err := requireAuth(deps.Posts); err != nil {
		return err
	}

	ctx := context.Background()
	req := &model.CreatePostRequest{
		Text:       *text,
		Visibility: *visibility,
	}

	if err := attachMedia(ctx, deps, req, *image, *video); err != nil {
		return err
	}

	post, err := deps.Posts.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("post create: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Post created: %s\n", post.ID)
	return nil
}

// attachMedia uploads the given image or video and sets the MediaURN on req.
func attachMedia(ctx context.Context, deps *Deps, req *model.CreatePostRequest, imagePath, videoPath string) error {
	if imagePath == "" && videoPath == "" {
		return nil
	}
	if err := requireAuth(deps.Media); err != nil {
		return err
	}

	filePath := imagePath
	mediaType := "IMAGE"
	if videoPath != "" {
		filePath = videoPath
		mediaType = "VIDEO"
	}

	upload, err := deps.Media.InitUpload(ctx, "me", mediaType)
	if err != nil {
		return fmt.Errorf("post create: init upload: %w", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("post create: open file: %w", err)
	}
	defer f.Close()

	if err := deps.Media.Upload(ctx, upload.UploadURL, f); err != nil {
		return fmt.Errorf("post create: upload: %w", err)
	}

	req.MediaURN = upload.MediaURN
	return nil
}

// validateVisibility checks that the visibility value is valid.
func validateVisibility(v string) error {
	switch v {
	case "PUBLIC", "CONNECTIONS":
		return nil
	default:
		return fmt.Errorf("invalid visibility %q: use PUBLIC or CONNECTIONS", v)
	}
}
