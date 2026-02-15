// Package command implements the lcli CLI subcommands.
package command

import (
	"context"
	"io"
	"reflect"

	"github.com/Softorize/lcli/internal/config"
	"github.com/Softorize/lcli/internal/model"
	"github.com/Softorize/lcli/internal/output"
)

// Profiler retrieves LinkedIn profile data.
type Profiler interface {
	Me(ctx context.Context) (*model.Profile, error)
	GetByID(ctx context.Context, id string) (*model.Profile, error)
}

// Poster manages LinkedIn posts.
type Poster interface {
	Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	Get(ctx context.Context, urn string) (*model.Post, error)
	Delete(ctx context.Context, urn string) error
	ListByAuthor(ctx context.Context, authorURN string, start, count int) (*model.PostList, error)
}

// Commenter manages comments on LinkedIn posts.
type Commenter interface {
	Create(ctx context.Context, req *model.CreateCommentRequest) (*model.Comment, error)
	List(ctx context.Context, postURN string, start, count int) (*model.CommentList, error)
	Delete(ctx context.Context, commentURN string) error
}

// Reacter manages reactions on LinkedIn entities.
type Reacter interface {
	React(ctx context.Context, actorURN, entityURN string, reaction model.ReactionType) error
	Unreact(ctx context.Context, actorURN, entityURN string) error
	List(ctx context.Context, entityURN string, start, count int) (*model.ReactionList, error)
}

// MediaUploader handles media upload workflows.
type MediaUploader interface {
	InitUpload(ctx context.Context, owner string, mediaType string) (*model.MediaUpload, error)
	Upload(ctx context.Context, uploadURL string, data io.Reader) error
	GetStatus(ctx context.Context, mediaURN string) (*model.MediaStatus, error)
}

// OrgReader retrieves organization data and statistics.
type OrgReader interface {
	Get(ctx context.Context, id int64) (*model.Organization, error)
	GetByVanity(ctx context.Context, vanityName string) (*model.Organization, error)
	FollowerStats(ctx context.Context, orgURN string) (*model.OrgFollowerStats, error)
	PageStats(ctx context.Context, orgURN string) (*model.OrgPageStats, error)
}

// AnalyticsReader retrieves LinkedIn analytics data.
type AnalyticsReader interface {
	PostAnalytics(ctx context.Context, postURN string) (map[string]any, error)
	ProfileViews(ctx context.Context) (int, error)
}

// Deps holds injected dependencies for all commands.
type Deps struct {
	// Cfg provides access to application configuration.
	Cfg *config.Config
	// Profile provides access to LinkedIn profile endpoints.
	Profile Profiler
	// Posts provides access to LinkedIn post endpoints.
	Posts Poster
	// Comments provides access to LinkedIn comment endpoints.
	Comments Commenter
	// Reactions provides access to LinkedIn reaction endpoints.
	Reactions Reacter
	// Media provides access to LinkedIn media upload endpoints.
	Media MediaUploader
	// Orgs provides access to LinkedIn organization endpoints.
	Orgs OrgReader
	// Analytics provides access to LinkedIn analytics endpoints.
	Analytics AnalyticsReader
	// Output is the configured printer for structured results.
	Output *output.Printer
	// Stdout is the writer for command results.
	Stdout io.Writer
	// Stderr is the writer for progress messages and errors.
	Stderr io.Writer
}

// requireAuth returns an error if the given service pointer is nil,
// indicating the user has not authenticated yet.
func requireAuth(svc any) error {
	if svc == nil {
		return errNotAuthenticated
	}
	v := reflect.ValueOf(svc)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return errNotAuthenticated
	}
	return nil
}
