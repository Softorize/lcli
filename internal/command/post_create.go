package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Softorize/lcli/internal/model"
)

// runPostCreate handles the post create subcommand.
func runPostCreate(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("post create", flag.ContinueOnError)
	text := fs.String("text", "", "Post text content (required)")
	visibility := fs.String("visibility", "PUBLIC", "Visibility: PUBLIC or CONNECTIONS")
	image := fs.String("image", "", "Path to image file to attach")
	video := fs.String("video", "", "Path to video file to attach")
	document := fs.String("document", "", "Path to PDF document for carousel post")
	title := fs.String("title", "", "Title for document/carousel post")
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

	// Resolve the full person URN for the author (required by API v202601+).
	if err := requireAuth(deps.Profile); err == nil {
		if profile, err := deps.Profile.Me(ctx); err == nil {
			req.AuthorURN = "urn:li:person:" + profile.ID
		}
	}

	if err := attachMedia(ctx, deps, req, *image, *video, *document, *title); err != nil {
		return err
	}

	post, err := deps.Posts.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("post create: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Post created: %s\n", post.ID)
	return nil
}

// attachMedia uploads the given image, video, or document and sets the MediaURN on req.
func attachMedia(ctx context.Context, deps *Deps, req *model.CreatePostRequest, imagePath, videoPath, documentPath, documentTitle string) error {
	if imagePath == "" && videoPath == "" && documentPath == "" {
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
	if documentPath != "" {
		filePath = documentPath
		mediaType = "DOCUMENT"
	}

	// The document API requires the full person URN (urn:li:person:ID),
	// while images and videos accept "me".
	owner := "me"
	if mediaType == "DOCUMENT" {
		if err := requireAuth(deps.Profile); err != nil {
			return fmt.Errorf("post create: profile required for document upload: %w", err)
		}
		profile, err := deps.Profile.Me(ctx)
		if err != nil {
			return fmt.Errorf("post create: resolve owner: %w", err)
		}
		owner = "urn:li:person:" + profile.ID
	}

	upload, err := deps.Media.InitUpload(ctx, owner, mediaType)
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
	if documentTitle != "" {
		req.MediaTitle = documentTitle
	}
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
