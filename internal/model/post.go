package model

import "time"

// Post represents a LinkedIn post (also known as a share or UGC post).
type Post struct {
	ID             string    `json:"id"`
	Author         string    `json:"author"`
	Text           string    `json:"text"`
	MediaCategory  string    `json:"mediaCategory"`
	Visibility     string    `json:"visibility"`
	CreatedAt      time.Time `json:"createdAt"`
	LifecycleState string    `json:"lifecycleState"`
}

// CreatePostRequest contains the fields needed to create a new post.
type CreatePostRequest struct {
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
	MediaURN   string `json:"mediaUrn,omitempty"`
	MediaTitle string `json:"mediaTitle,omitempty"`
	AuthorURN  string `json:"authorUrn,omitempty"`
}

// PostList is a paginated list of posts.
type PostList struct {
	Elements []Post  `json:"elements"`
	Paging   *Paging `json:"paging,omitempty"`
}

// Paging holds pagination metadata for list responses.
type Paging struct {
	Count int `json:"count"`
	Start int `json:"start"`
	Total int `json:"total"`
}
