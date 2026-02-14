package linkedin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/toto/lcli/internal/model"
)

// PostService provides access to LinkedIn post (UGC) endpoints.
type PostService struct {
	doer Doer
}

// NewPostService creates a PostService backed by the given Doer.
func NewPostService(d Doer) *PostService {
	return &PostService{doer: d}
}

// postBody is the request payload for creating a post via the REST API.
type postBody struct {
	Author       string          `json:"author"`
	Commentary   string          `json:"commentary"`
	Visibility   string          `json:"visibility"`
	Distribution postDistribution `json:"distribution"`
	Content      *postContent    `json:"content,omitempty"`
}

// postDistribution controls the feed distribution of a post.
type postDistribution struct {
	FeedDistribution string `json:"feedDistribution"`
}

// postContent holds optional media content for a post.
type postContent struct {
	Media *postMedia `json:"media,omitempty"`
}

// postMedia references an uploaded media asset.
type postMedia struct {
	ID string `json:"id"`
}

// postResponse is the raw API response for a single post.
type postResponse struct {
	ID             string `json:"id"`
	Author         string `json:"author"`
	Commentary     string `json:"commentary"`
	Visibility     string `json:"visibility"`
	CreatedAt      int64  `json:"createdAt"`
	LifecycleState string `json:"lifecycleState"`
	Distribution   struct {
		FeedDistribution string `json:"feedDistribution"`
	} `json:"distribution"`
	Content *struct {
		Media *struct {
			ID string `json:"id"`
		} `json:"media"`
	} `json:"content"`
}

// toPost converts a raw API response into a domain Post.
func (r *postResponse) toPost() *model.Post {
	p := &model.Post{
		ID:             r.ID,
		Author:         r.Author,
		Text:           r.Commentary,
		Visibility:     r.Visibility,
		LifecycleState: r.LifecycleState,
		MediaCategory:  "NONE",
	}

	if r.CreatedAt > 0 {
		p.CreatedAt = time.UnixMilli(r.CreatedAt)
	}

	if r.Content != nil && r.Content.Media != nil {
		p.MediaCategory = "IMAGE"
	}

	return p
}

// Create publishes a new post on LinkedIn.
func (s *PostService) Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	body := postBody{
		Author:     "me",
		Commentary: req.Text,
		Visibility: req.Visibility,
		Distribution: postDistribution{
			FeedDistribution: "MAIN_FEED",
		},
	}

	if req.MediaURN != "" {
		body.Content = &postContent{
			Media: &postMedia{ID: req.MediaURN},
		}
	}

	resp, err := s.doer.Do(ctx, http.MethodPost, "/posts", body)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	var raw postResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return raw.toPost(), nil
}

// Get retrieves a single post by its URN.
func (s *PostService) Get(ctx context.Context, urn string) (*model.Post, error) {
	path := fmt.Sprintf("/posts/%s", url.PathEscape(urn))

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get post %s: %w", urn, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("get post %s: %w", urn, err)
	}

	var raw postResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("get post %s: %w", urn, err)
	}

	return raw.toPost(), nil
}

// Delete removes a post by its URN.
func (s *PostService) Delete(ctx context.Context, urn string) error {
	path := fmt.Sprintf("/posts/%s", url.PathEscape(urn))

	resp, err := s.doer.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete post %s: %w", urn, err)
	}

	if err := checkError(resp); err != nil {
		return fmt.Errorf("delete post %s: %w", urn, err)
	}

	drainBody(resp)
	return nil
}

// ListByAuthor returns posts authored by the given URN with pagination.
func (s *PostService) ListByAuthor(ctx context.Context, authorURN string, start, count int) (*model.PostList, error) {
	path := fmt.Sprintf("/posts?author=%s&q=author&start=%s&count=%s",
		url.QueryEscape(authorURN),
		strconv.Itoa(start),
		strconv.Itoa(count),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list posts by %s: %w", authorURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("list posts by %s: %w", authorURN, err)
	}

	var raw struct {
		Elements []postResponse `json:"elements"`
		Paging   *model.Paging  `json:"paging"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("list posts by %s: %w", authorURN, err)
	}

	list := &model.PostList{Paging: raw.Paging}
	for i := range raw.Elements {
		list.Elements = append(list.Elements, *raw.Elements[i].toPost())
	}

	return list, nil
}
