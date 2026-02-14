// Package command implements the lcli CLI subcommands.
package command

import (
	"io"
	"reflect"

	"github.com/Softorize/lcli/internal/config"
	"github.com/Softorize/lcli/internal/linkedin"
	"github.com/Softorize/lcli/internal/output"
)

// Deps holds injected dependencies for all commands.
type Deps struct {
	// Cfg provides access to application configuration.
	Cfg *config.Config
	// Profile provides access to LinkedIn profile endpoints.
	Profile *linkedin.ProfileService
	// Posts provides access to LinkedIn post endpoints.
	Posts *linkedin.PostService
	// Comments provides access to LinkedIn comment endpoints.
	Comments *linkedin.CommentService
	// Reactions provides access to LinkedIn reaction endpoints.
	Reactions *linkedin.ReactionService
	// Media provides access to LinkedIn media upload endpoints.
	Media *linkedin.MediaService
	// Orgs provides access to LinkedIn organization endpoints.
	Orgs *linkedin.OrgService
	// Analytics provides access to LinkedIn analytics endpoints.
	Analytics *linkedin.AnalyticsService
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
