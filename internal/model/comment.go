package model

import "time"

// Comment represents a comment on a LinkedIn post.
type Comment struct {
	ID            string    `json:"id"`
	Author        string    `json:"author"`
	Text          string    `json:"text"`
	CreatedAt     time.Time `json:"createdAt"`
	ParentComment string    `json:"parentComment,omitempty"`
}

// CreateCommentRequest contains the fields needed to create a comment.
type CreateCommentRequest struct {
	PostURN string `json:"postUrn"`
	Text    string `json:"text"`
}

// CommentList is a paginated list of comments.
type CommentList struct {
	Elements []Comment `json:"elements"`
	Paging   *Paging   `json:"paging,omitempty"`
}
