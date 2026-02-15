package command

import (
	"bytes"
	"context"
	"io"

	"github.com/Softorize/lcli/internal/model"
	"github.com/Softorize/lcli/internal/output"
)

// testDeps creates a Deps with captured stdout/stderr and optional service mocks.
func testDeps() (*Deps, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	return &Deps{
		Output: output.NewPrinter(stdout, output.FormatTable),
		Stdout: stdout,
		Stderr: stderr,
	}, stdout, stderr
}

// mockProfiler implements Profiler for testing.
type mockProfiler struct {
	meFunc    func(ctx context.Context) (*model.Profile, error)
	getByIDFunc func(ctx context.Context, id string) (*model.Profile, error)
}

func (m *mockProfiler) Me(ctx context.Context) (*model.Profile, error) {
	return m.meFunc(ctx)
}

func (m *mockProfiler) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	return m.getByIDFunc(ctx, id)
}

// mockPoster implements Poster for testing.
type mockPoster struct {
	createFunc      func(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	getFunc         func(ctx context.Context, urn string) (*model.Post, error)
	deleteFunc      func(ctx context.Context, urn string) error
	listByAuthorFunc func(ctx context.Context, authorURN string, start, count int) (*model.PostList, error)
}

func (m *mockPoster) Create(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	return m.createFunc(ctx, req)
}

func (m *mockPoster) Get(ctx context.Context, urn string) (*model.Post, error) {
	return m.getFunc(ctx, urn)
}

func (m *mockPoster) Delete(ctx context.Context, urn string) error {
	return m.deleteFunc(ctx, urn)
}

func (m *mockPoster) ListByAuthor(ctx context.Context, authorURN string, start, count int) (*model.PostList, error) {
	return m.listByAuthorFunc(ctx, authorURN, start, count)
}

// mockCommenter implements Commenter for testing.
type mockCommenter struct {
	createFunc func(ctx context.Context, req *model.CreateCommentRequest) (*model.Comment, error)
	listFunc   func(ctx context.Context, postURN string, start, count int) (*model.CommentList, error)
	deleteFunc func(ctx context.Context, commentURN string) error
}

func (m *mockCommenter) Create(ctx context.Context, req *model.CreateCommentRequest) (*model.Comment, error) {
	return m.createFunc(ctx, req)
}

func (m *mockCommenter) List(ctx context.Context, postURN string, start, count int) (*model.CommentList, error) {
	return m.listFunc(ctx, postURN, start, count)
}

func (m *mockCommenter) Delete(ctx context.Context, commentURN string) error {
	return m.deleteFunc(ctx, commentURN)
}

// mockReacter implements Reacter for testing.
type mockReacter struct {
	reactFunc   func(ctx context.Context, actorURN, entityURN string, reaction model.ReactionType) error
	unreactFunc func(ctx context.Context, actorURN, entityURN string) error
	listFunc    func(ctx context.Context, entityURN string, start, count int) (*model.ReactionList, error)
}

func (m *mockReacter) React(ctx context.Context, actorURN, entityURN string, reaction model.ReactionType) error {
	return m.reactFunc(ctx, actorURN, entityURN, reaction)
}

func (m *mockReacter) Unreact(ctx context.Context, actorURN, entityURN string) error {
	return m.unreactFunc(ctx, actorURN, entityURN)
}

func (m *mockReacter) List(ctx context.Context, entityURN string, start, count int) (*model.ReactionList, error) {
	return m.listFunc(ctx, entityURN, start, count)
}

// mockMediaUploader implements MediaUploader for testing.
type mockMediaUploader struct {
	initUploadFunc func(ctx context.Context, owner string, mediaType string) (*model.MediaUpload, error)
	uploadFunc     func(ctx context.Context, uploadURL string, data io.Reader) error
	getStatusFunc  func(ctx context.Context, mediaURN string) (*model.MediaStatus, error)
}

func (m *mockMediaUploader) InitUpload(ctx context.Context, owner string, mediaType string) (*model.MediaUpload, error) {
	return m.initUploadFunc(ctx, owner, mediaType)
}

func (m *mockMediaUploader) Upload(ctx context.Context, uploadURL string, data io.Reader) error {
	return m.uploadFunc(ctx, uploadURL, data)
}

func (m *mockMediaUploader) GetStatus(ctx context.Context, mediaURN string) (*model.MediaStatus, error) {
	return m.getStatusFunc(ctx, mediaURN)
}

// mockOrgReader implements OrgReader for testing.
type mockOrgReader struct {
	getFunc           func(ctx context.Context, id int64) (*model.Organization, error)
	getByVanityFunc   func(ctx context.Context, vanityName string) (*model.Organization, error)
	followerStatsFunc func(ctx context.Context, orgURN string) (*model.OrgFollowerStats, error)
	pageStatsFunc     func(ctx context.Context, orgURN string) (*model.OrgPageStats, error)
}

func (m *mockOrgReader) Get(ctx context.Context, id int64) (*model.Organization, error) {
	return m.getFunc(ctx, id)
}

func (m *mockOrgReader) GetByVanity(ctx context.Context, vanityName string) (*model.Organization, error) {
	return m.getByVanityFunc(ctx, vanityName)
}

func (m *mockOrgReader) FollowerStats(ctx context.Context, orgURN string) (*model.OrgFollowerStats, error) {
	return m.followerStatsFunc(ctx, orgURN)
}

func (m *mockOrgReader) PageStats(ctx context.Context, orgURN string) (*model.OrgPageStats, error) {
	return m.pageStatsFunc(ctx, orgURN)
}

// mockAnalyticsReader implements AnalyticsReader for testing.
type mockAnalyticsReader struct {
	postAnalyticsFunc func(ctx context.Context, postURN string) (map[string]any, error)
	profileViewsFunc  func(ctx context.Context) (int, error)
}

func (m *mockAnalyticsReader) PostAnalytics(ctx context.Context, postURN string) (map[string]any, error) {
	return m.postAnalyticsFunc(ctx, postURN)
}

func (m *mockAnalyticsReader) ProfileViews(ctx context.Context) (int, error) {
	return m.profileViewsFunc(ctx)
}
