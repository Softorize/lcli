package linkedin

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Softorize/lcli/internal/model"
)

// MediaService provides access to LinkedIn media upload endpoints.
type MediaService struct {
	doer Doer
}

// NewMediaService creates a MediaService backed by the given Doer.
func NewMediaService(d Doer) *MediaService {
	return &MediaService{doer: d}
}

// initUploadRequest is the request body for initializing a media upload.
type initUploadRequest struct {
	InitializeUploadRequest initUploadOwner `json:"initializeUploadRequest"`
}

// initUploadOwner wraps the owner URN for the upload init call.
type initUploadOwner struct {
	Owner string `json:"owner"`
}

// initUploadResponse is the raw response from the upload init endpoint.
type initUploadResponse struct {
	Value struct {
		UploadURL          string `json:"uploadUrl"`
		MediaURNOrImage    string `json:"image,omitempty"`
		MediaURNOrVideo    string `json:"video,omitempty"`
		MediaURNOrDocument string `json:"document,omitempty"`
		UploadToken        string `json:"uploadToken,omitempty"`
	} `json:"value"`
}

// InitUpload initializes a media upload and returns the upload details.
// mediaType must be "IMAGE", "VIDEO", or "DOCUMENT".
func (s *MediaService) InitUpload(ctx context.Context, owner string, mediaType string) (*model.MediaUpload, error) {
	var path string
	switch mediaType {
	case "IMAGE":
		path = "/images?action=initializeUpload"
	case "VIDEO":
		path = "/videos?action=initializeUpload"
	case "DOCUMENT":
		path = "/documents?action=initializeUpload"
	default:
		return nil, fmt.Errorf("init upload: unsupported media type %q", mediaType)
	}

	body := initUploadRequest{
		InitializeUploadRequest: initUploadOwner{Owner: owner},
	}

	resp, err := s.doer.Do(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, fmt.Errorf("init upload: %w", err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("init upload: %w", err)
	}

	var raw initUploadResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("init upload: %w", err)
	}

	mediaURN := raw.Value.MediaURNOrImage
	if mediaURN == "" {
		mediaURN = raw.Value.MediaURNOrVideo
	}
	if mediaURN == "" {
		mediaURN = raw.Value.MediaURNOrDocument
	}

	return &model.MediaUpload{
		UploadURL:   raw.Value.UploadURL,
		MediaURN:    mediaURN,
		UploadToken: raw.Value.UploadToken,
	}, nil
}

// Upload sends binary data to the provided upload URL.
func (s *MediaService) Upload(ctx context.Context, uploadURL string, data io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, data)
	if err != nil {
		return fmt.Errorf("build upload request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload media: status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// GetStatus retrieves the processing status of an uploaded media asset.
func (s *MediaService) GetStatus(ctx context.Context, mediaURN string) (*model.MediaStatus, error) {
	path := fmt.Sprintf("/assets/%s", mediaURN)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get media status %s: %w", mediaURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get media status %s: %w", mediaURN, err)
	}

	var raw struct {
		URN    string `json:"id"`
		Status string `json:"status"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get media status %s: %w", mediaURN, err)
	}

	return &model.MediaStatus{
		URN:    raw.URN,
		Status: raw.Status,
	}, nil
}
