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

// CommentService provides access to LinkedIn comment endpoints.
type CommentService struct {
	doer Doer
}

// NewCommentService creates a CommentService backed by the given Doer.
func NewCommentService(d Doer) *CommentService {
	return &CommentService{doer: d}
}

// commentBody is the request payload for creating a comment.
type commentBody struct {
	Actor   string `json:"actor"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

// commentResponse is the raw API response for a single comment.
type commentResponse struct {
	ID            string `json:"$URN"`
	Actor         string `json:"actor"`
	CreatedAt     int64  `json:"created"`
	Message       struct {
		Text string `json:"text"`
	} `json:"message"`
	ParentComment string `json:"parentComment,omitempty"`
}

// toComment converts a raw API response into a domain Comment.
func (r *commentResponse) toComment() *model.Comment {
	c := &model.Comment{
		ID:            r.ID,
		Author:        r.Actor,
		Text:          r.Message.Text,
		ParentComment: r.ParentComment,
	}

	if r.CreatedAt > 0 {
		c.CreatedAt = time.UnixMilli(r.CreatedAt)
	}

	return c
}

// Create adds a new comment to a post.
func (s *CommentService) Create(ctx context.Context, req *model.CreateCommentRequest) (*model.Comment, error) {
	path := fmt.Sprintf("/socialActions/%s/comments", url.PathEscape(req.PostURN))

	body := commentBody{Actor: "me"}
	body.Message.Text = req.Text

	resp, err := s.doer.Do(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, fmt.Errorf("create comment on %s: %w", req.PostURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("create comment on %s: %w", req.PostURN, err)
	}

	var raw commentResponse
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("create comment on %s: %w", req.PostURN, err)
	}

	return raw.toComment(), nil
}

// List retrieves comments for a post with pagination.
func (s *CommentService) List(ctx context.Context, postURN string, start, count int) (*model.CommentList, error) {
	path := fmt.Sprintf("/socialActions/%s/comments?start=%s&count=%s",
		url.PathEscape(postURN),
		strconv.Itoa(start),
		strconv.Itoa(count),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list comments on %s: %w", postURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("list comments on %s: %w", postURN, err)
	}

	var raw struct {
		Elements []commentResponse `json:"elements"`
		Paging   *model.Paging     `json:"paging"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("list comments on %s: %w", postURN, err)
	}

	list := &model.CommentList{Paging: raw.Paging}
	for i := range raw.Elements {
		list.Elements = append(list.Elements, *raw.Elements[i].toComment())
	}

	return list, nil
}

// Delete removes a comment by its URN.
func (s *CommentService) Delete(ctx context.Context, commentURN string) error {
	path := fmt.Sprintf("/socialActions/%s/comments/%s",
		url.PathEscape(commentURN),
		url.PathEscape(commentURN),
	)

	resp, err := s.doer.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete comment %s: %w", commentURN, err)
	}

	if err := checkError(resp); err != nil {
		return fmt.Errorf("delete comment %s: %w", commentURN, err)
	}

	drainBody(resp)
	return nil
}
