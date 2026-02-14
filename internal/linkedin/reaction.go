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

// ReactionService provides access to LinkedIn reaction endpoints.
type ReactionService struct {
	doer Doer
}

// NewReactionService creates a ReactionService backed by the given Doer.
func NewReactionService(d Doer) *ReactionService {
	return &ReactionService{doer: d}
}

// reactionBody is the request payload for creating a reaction.
type reactionBody struct {
	Root  string `json:"root"`
	Type  string `json:"reactionType"`
	Actor string `json:"actor"`
}

// reactionResponse is the raw API response for a single reaction.
type reactionResponse struct {
	Actor     string `json:"actor"`
	Type      string `json:"reactionType"`
	CreatedAt int64  `json:"created"`
}

// toReaction converts a raw API response into a domain Reaction.
func (r *reactionResponse) toReaction() *model.Reaction {
	rx := &model.Reaction{
		Actor: r.Actor,
		Type:  model.ReactionType(r.Type),
	}

	if r.CreatedAt > 0 {
		rx.CreatedAt = time.UnixMilli(r.CreatedAt)
	}

	return rx
}

// React adds a reaction to a LinkedIn entity.
func (s *ReactionService) React(ctx context.Context, actorURN, entityURN string, reaction model.ReactionType) error {
	body := reactionBody{
		Root:  entityURN,
		Type:  string(reaction),
		Actor: actorURN,
	}

	resp, err := s.doer.Do(ctx, http.MethodPost, "/reactions", body)
	if err != nil {
		return fmt.Errorf("react on %s: %w", entityURN, err)
	}

	if err := checkError(resp); err != nil {
		return fmt.Errorf("react on %s: %w", entityURN, err)
	}

	drainBody(resp)
	return nil
}

// Unreact removes the actor's reaction from a LinkedIn entity.
func (s *ReactionService) Unreact(ctx context.Context, actorURN, entityURN string) error {
	path := fmt.Sprintf("/reactions/(actor:%s,entity:%s)",
		url.PathEscape(actorURN),
		url.PathEscape(entityURN),
	)

	resp, err := s.doer.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("unreact on %s: %w", entityURN, err)
	}

	if err := checkError(resp); err != nil {
		return fmt.Errorf("unreact on %s: %w", entityURN, err)
	}

	drainBody(resp)
	return nil
}

// List retrieves reactions for an entity with pagination.
func (s *ReactionService) List(ctx context.Context, entityURN string, start, count int) (*model.ReactionList, error) {
	path := fmt.Sprintf("/reactions/(entity:%s)?start=%s&count=%s",
		url.PathEscape(entityURN),
		strconv.Itoa(start),
		strconv.Itoa(count),
	)

	resp, err := s.doer.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list reactions on %s: %w", entityURN, err)
	}

	if err := checkError(resp); err != nil {
		return nil, fmt.Errorf("list reactions on %s: %w", entityURN, err)
	}

	var raw struct {
		Elements []reactionResponse `json:"elements"`
		Paging   *model.Paging      `json:"paging"`
	}
	if err := decodeJSON(resp, &raw); err != nil {
		return nil, fmt.Errorf("list reactions on %s: %w", entityURN, err)
	}

	list := &model.ReactionList{Paging: raw.Paging}
	for i := range raw.Elements {
		list.Elements = append(list.Elements, *raw.Elements[i].toReaction())
	}

	return list, nil
}
