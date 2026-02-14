package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runMedia dispatches to media subcommands: upload.
func runMedia(args []string, deps *Deps) error {
	if len(args) == 0 {
		printMediaUsage(deps)
		return nil
	}

	switch args[0] {
	case "upload":
		return runMediaUpload(args[1:], deps)
	case "-help", "--help", "-h":
		printMediaUsage(deps)
		return nil
	default:
		return fmt.Errorf("media: unknown subcommand %q", args[0])
	}
}

// printMediaUsage writes media command help text.
func printMediaUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli media <subcommand> [flags]

Subcommands:
  upload    Upload an image or video file

Use "lcli media <subcommand> -help" for more information.
`)
}

// runMediaUpload handles the media upload subcommand.
func runMediaUpload(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("media upload", flag.ContinueOnError)
	mediaType := fs.String("type", "", "Media type: image or video (auto-detected if not set)")
	owner := fs.String("owner", "me", "Owner URN (defaults to 'me')")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("media upload: file path argument is required")
	}

	filePath := fs.Arg(0)
	detectedType := *mediaType
	if detectedType == "" {
		detectedType = detectMediaType(filePath)
	}

	if detectedType != "image" && detectedType != "video" {
		return fmt.Errorf("media upload: unable to detect type for %q, use --type", filePath)
	}

	if err := requireAuth(deps.Media); err != nil {
		return err
	}

	apiType := strings.ToUpper(detectedType)
	ctx := context.Background()

	upload, err := deps.Media.InitUpload(ctx, *owner, apiType)
	if err != nil {
		return fmt.Errorf("media upload: %w", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("media upload: open file: %w", err)
	}
	defer f.Close()

	if err := deps.Media.Upload(ctx, upload.UploadURL, f); err != nil {
		return fmt.Errorf("media upload: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Upload complete.\n")
	fmt.Fprintf(deps.Stdout, "Media URN: %s\n", upload.MediaURN)
	return nil
}

// detectMediaType guesses the media type from the file extension.
func detectMediaType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".mp4", ".mov", ".avi", ".wmv", ".webm":
		return "video"
	default:
		return ""
	}
}
