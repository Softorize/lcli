package model

import "time"

// ReactionType represents the kind of reaction on a LinkedIn post.
type ReactionType string

const (
	// ReactionLike is a standard like reaction.
	ReactionLike ReactionType = "LIKE"
	// ReactionCelebrate is a celebratory reaction.
	ReactionCelebrate ReactionType = "CELEBRATE"
	// ReactionSupport is a supportive reaction.
	ReactionSupport ReactionType = "SUPPORT"
	// ReactionLove is a love reaction.
	ReactionLove ReactionType = "LOVE"
	// ReactionInsightful marks content as insightful.
	ReactionInsightful ReactionType = "INSIGHTFUL"
	// ReactionFunny marks content as funny.
	ReactionFunny ReactionType = "FUNNY"
)

// Reaction represents a single reaction on a LinkedIn entity.
type Reaction struct {
	Actor     string       `json:"actor"`
	Type      ReactionType `json:"type"`
	CreatedAt time.Time    `json:"createdAt"`
}

// ReactionSummary aggregates reaction counts for a post.
type ReactionSummary struct {
	PostURN    string              `json:"postUrn"`
	TotalCount int                 `json:"totalCount"`
	ByType     map[ReactionType]int `json:"byType"`
}

// ReactionList is a paginated list of reactions.
type ReactionList struct {
	Elements []Reaction `json:"elements"`
	Paging   *Paging    `json:"paging,omitempty"`
}
