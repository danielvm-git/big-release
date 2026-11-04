// story: e03s02 e21s01 e21s02 e22s03
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
)

// GitHubPlugin creates GitHub releases via the API.
type GitHubPlugin struct {
	client                 HTTPClient
	apiBaseURL             string // for testing; empty = use default GitHub API
	uploadBaseURL          string // for testing; empty = use default uploads host
	assets                 []algorithm.AssetConfig
	draftRelease           bool
	releaseNameTemplate    string
	releaseBodyTemplate    string
	successComment         string
	releasedLabels         []string
	discussionCategoryName string
	makeLatest             *bool

	releaseMu        sync.Mutex
	createdReleaseID int64
}

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewGitHubPlugin creates a new GitHubPlugin.
func NewGitHubPlugin() *GitHubPlugin {
	return &GitHubPlugin{
		client: http.DefaultClient,
	}
}

// Configure decodes the github plugin's typed config (assets) from the
// raw PluginConfigs map entry. Implements ConfigurablePlugin.
func (p *GitHubPlugin) Configure(raw map[string]interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	// Re-marshal to YAML/JSON then decode into the typed struct, so users
	// can write assets as either a list of strings or {path,label} maps.
	data, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("github plugin: failed to re-marshal config: %w", err)
	}
	var cfg algorithm.GitHubConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("github plugin: invalid config: %w", err)
	}
	p.assets = cfg.Assets
	p.draftRelease = cfg.DraftRelease
	p.releaseNameTemplate = cfg.ReleaseName
	p.releaseBodyTemplate = cfg.ReleaseBody
	p.successComment = cfg.SuccessComment
	p.releasedLabels = cfg.ReleasedLabels
	if v, ok := raw["discussionCategoryName"].(string); ok {
		p.discussionCategoryName = strings.TrimSpace(v)
	}
	if v, ok := raw["makeLatest"].(bool); ok {
		p.makeLatest = &v
	}
	return nil
}

// Name returns the plugin name.
func (p *GitHubPlugin) Name() string {
	return "github"
}

// VerifyConditions checks that required environment variables are set.
func (p *GitHubPlugin) VerifyConditions(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	// Skip credential checks in dry-run mode — we won't be publishing.
	if ctx.DryRun {
		return nil
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}
	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return fmt.Errorf("GITHUB_REPOSITORY environment variable is required (format: owner/repo)")
	}
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("GITHUB_REPOSITORY must be in format owner/repo, got %q", repo)
	}
	return nil
}

// AnalyzeCommits is not applicable for the GitHub plugin.
func (p *GitHubPlugin) AnalyzeCommits(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return "", nil
}

// VerifyRelease is not applicable for the GitHub plugin.
func (p *GitHubPlugin) VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// GenerateNotes is not applicable for the GitHub plugin.
func (p *GitHubPlugin) GenerateNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	return "", nil
}

// Prepare is not applicable for the GitHub plugin.
func (p *GitHubPlugin) Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// createReleaseRequest represents a GitHub release creation request.
type createReleaseRequest struct {
	TagName                string `json:"tag_name"`
	Name                   string `json:"name"`
	Body                   string `json:"body"`
	Draft                  bool   `json:"draft"`
	Prerelease             bool   `json:"prerelease"`
	GenerateReleaseNotes   bool   `json:"generate_release_notes"`
	DiscussionCategoryName string `json:"discussion_category_name,omitempty"`
	MakeLatest             string `json:"make_latest,omitempty"`
}

// templateContext is the variable scope passed to release name/body
// templates. Mirrors @semantic-release/github's lodash template context.
type templateContext struct {
	Version     string
	Date        string
	Branch      string
	Notes       string
	NextRelease *algorithm.Release
	LastRelease *algorithm.Release
}

func (p *GitHubPlugin) buildReleasePayload(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) ([]byte, error) {
	version := state.NextRelease.Version
	notes := state.NextRelease.Notes
	releaseType := string(state.NextRelease.Type)

	tagFormat := ""
	if ctx.Config != nil {
		tagFormat = ctx.Config.TagFormat
	}
	// Tag actually created by the git plugin (BUG-tag-ignores-tagformat) —
	// must match exactly, or GitHub creates a second, orphan tag for this release.
	tag := git.FormatTag(tagFormat, version)

	tctx := templateContext{
		Version:     version,
		Date:        time.Now().UTC().Format(time.RFC3339),
		Notes:       notes,
		NextRelease: state.NextRelease,
		LastRelease: state.LastRelease,
	}
	if ctx.Branch != nil {
		tctx.Branch = ctx.Branch.Name
	}

	// Default display name is always "v{version}" regardless of tagFormat —
	// a cosmetic convention, independent of the actual tag_name below.
	name, err := p.renderTemplate("release name", p.releaseNameTemplate, "v"+version, tctx)
	if err != nil {
		return nil, err
	}
	body, err := p.renderTemplate("release body", p.releaseBodyTemplate, notes, tctx)
	if err != nil {
		return nil, err
	}

	req := &createReleaseRequest{
		TagName:              tag,
		Name:                 name,
		Body:                 body,
		Draft:                p.draftRelease,
		Prerelease:           releaseType == "prerelease",
		GenerateReleaseNotes: notes == "",
	}
	if p.discussionCategoryName != "" {
		req.DiscussionCategoryName = p.discussionCategoryName
	}
	if ml := p.resolveMakeLatest(ctx); ml != "" {
		req.MakeLatest = ml
	}
	return json.Marshal(req)
}

// resolveMakeLatest returns the GitHub make_latest flag ("true"/"false").
func (p *GitHubPlugin) resolveMakeLatest(ctx *algorithm.ReadOnlyContext) string {
	if p.makeLatest != nil {
		if *p.makeLatest {
			return "true"
		}
		return "false"
	}
	if ctx.Branch != nil && ctx.Branch.Channel != "" && ctx.Branch.Channel != "latest" {
		return "false"
	}
	if ctx.Branch != nil &&
		ctx.Branch.Type == algorithm.BranchTypeRelease &&
		ctx.Branch.Name == "main" {
		return "true"
	}
	return "false"
}

// renderTemplate executes a Go text/template against tctx. When tmpl is
// empty, the fallback value is returned unchanged (preserving default
// behavior). A parse or execution error is returned with context.
func (p *GitHubPlugin) renderTemplate(label, tmpl, fallback string, tctx templateContext) (string, error) {
	if tmpl == "" {
		return fallback, nil
	}
	t, err := template.New(label).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("invalid %s template: %w", label, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, tctx); err != nil {
		return "", fmt.Errorf("failed to render %s template: %w", label, err)
	}
	return buf.String(), nil
}

func (p *GitHubPlugin) createHTTPRequest(url string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	return req, nil
}

func (p *GitHubPlugin) handleReleaseResponse(resp *http.Response, repo, version string) (int64, error) {
	if resp.StatusCode == http.StatusCreated {
		id, _ := parseReleaseID(resp.Body)
		return id, nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return 0, fmt.Errorf("GitHub authentication failed: invalid or expired GITHUB_TOKEN (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusForbidden:
		return 0, fmt.Errorf("GitHub permission denied (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusNotFound:
		return 0, fmt.Errorf("GitHub repository %q not found (HTTP %d): %s", repo, resp.StatusCode, string(respBody))
	case http.StatusUnprocessableEntity:
		return 0, fmt.Errorf("GitHub release already exists for tag %s (HTTP %d): %s", version, resp.StatusCode, string(respBody))
	default:
		return 0, fmt.Errorf("GitHub API returned unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
}

func (p *GitHubPlugin) releaseURL(repo string) string {
	return fmt.Sprintf("%s/repos/%s/releases", p.releaseURLBase(), repo)
}

func (p *GitHubPlugin) sendReleaseRequest(payload []byte, repo, version string) (int64, error) {
	req, err := p.createHTTPRequest(p.releaseURL(repo), payload)
	if err != nil {
		return 0, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to create GitHub release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return p.handleReleaseResponse(resp, repo, version)
}

// Publish creates a GitHub release via the API and, if assets are
// configured, uploads them to the release.
func (p *GitHubPlugin) Publish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error) {
	if ctx.DryRun {
		return nil, nil
	}
	payload, err := p.buildReleasePayload(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal release body: %w", err)
	}
	repo := os.Getenv("GITHUB_REPOSITORY")
	releaseID, err := p.sendReleaseRequest(payload, repo, state.NextRelease.Version)
	if err != nil {
		return nil, err
	}
	p.recordRelease(releaseID)

	// Upload configured assets (if any). Missing files warn but do not fail.
	if len(p.assets) > 0 && releaseID != 0 {
		if err := p.uploadAssets(repo, releaseID); err != nil {
			// NT-5: the release exists but has no complete asset set; remove
			// it so no half-published release is left behind.
			p.deleteRelease(repo, releaseID)
			p.forgetRelease()
			return nil, err
		}
	}

	// Draft releases are published atomically after asset upload (#13).
	if p.draftRelease && releaseID != 0 {
		if err := p.publishDraft(ctx, repo, releaseID); err != nil {
			// NT-6: the draft was created but never published; delete it so
			// no invisible orphan draft remains.
			p.deleteRelease(repo, releaseID)
			p.forgetRelease()
			return nil, err
		}
	}
	return nil, nil
}

// recordRelease remembers the release created by Publish so Rollback can
// delete it if the lifecycle fails later.
func (p *GitHubPlugin) recordRelease(releaseID int64) {
	p.releaseMu.Lock()
	p.createdReleaseID = releaseID
	p.releaseMu.Unlock()
}

// forgetRelease clears the remembered release after it has been removed.
func (p *GitHubPlugin) forgetRelease() {
	p.recordRelease(0)
}

// ResetRunState clears the release remembered by a previous run. The registry
// reuses this plugin instance, so without the reset a rollback in a later run
// could delete a release it did not create (NT-11).
func (p *GitHubPlugin) ResetRunState() { p.forgetRelease() }

// rememberedReleaseID returns the release created by Publish, or 0.
func (p *GitHubPlugin) rememberedReleaseID() int64 {
	p.releaseMu.Lock()
	defer p.releaseMu.Unlock()
	return p.createdReleaseID
}

// Rollback deletes the GitHub release created by a failed release. It is a
// no-op when Publish never created one, and idempotent once the release is
// gone (a 404 counts as already deleted).
func (p *GitHubPlugin) Rollback(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	releaseID := p.rememberedReleaseID()
	if releaseID == 0 {
		return nil
	}
	if err := p.deleteReleaseByID(os.Getenv("GITHUB_REPOSITORY"), releaseID); err != nil {
		return err
	}
	p.forgetRelease()
	return nil
}

// deleteRelease attempts to DELETE a release that was created earlier in a
// publish attempt that has now failed (NT-5/NT-6). Cleanup is best-effort:
// failures are logged to stderr and never replace the caller's original error.
func (p *GitHubPlugin) deleteRelease(repo string, releaseID int64) {
	if err := p.deleteReleaseByID(repo, releaseID); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %s\n", err)
	}
}

// deleteReleaseByID performs the DELETE and reports whether cleanup actually
// failed. A 404 means the release is already gone and counts as success, so
// Rollback stays idempotent.
func (p *GitHubPlugin) deleteReleaseByID(repo string, releaseID int64) error {
	deleteURL := fmt.Sprintf("%s/repos/%s/releases/%d", p.releaseURLBase(), repo, releaseID)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to build release cleanup request for %d: %w", releaseID, err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("release cleanup failed for %d: %w", releaseID, err)
	}
	defer func() { _ = resp.Body.Close() }()
	switch resp.StatusCode {
	case http.StatusNoContent, http.StatusOK, http.StatusNotFound:
		return nil
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("release cleanup failed for %d (HTTP %d): %s", releaseID, resp.StatusCode, string(body))
	}
}

// publishDraft flips a draft release to published via PATCH.
func (p *GitHubPlugin) publishDraft(ctx *algorithm.ReadOnlyContext, repo string, releaseID int64) error {
	patchURL := fmt.Sprintf("%s/repos/%s/releases/%d", p.releaseURLBase(), repo, releaseID)
	body := map[string]interface{}{"draft": false}
	if ml := p.resolveMakeLatest(ctx); ml != "" {
		body["make_latest"] = ml
	}
	if p.discussionCategoryName != "" {
		body["discussion_category_name"] = p.discussionCategoryName
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal draft publish payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPatch, patchURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create draft publish request: %w", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to publish draft release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("draft publish failed (HTTP %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// releaseURLBase returns the API base without the /releases suffix, used
// to build per-release URLs (PATCH, assets).
func (p *GitHubPlugin) releaseURLBase() string {
	baseURL := p.apiBaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return baseURL
}

// Fail is called on release failure.
func (p *GitHubPlugin) Fail(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, err error) error {
	return nil
}

// releaseResponse is the subset of the GitHub release payload we read.
type releaseResponse struct {
	ID int64 `json:"id"`
}

// parseReleaseID reads the release ID from a 201 response body.
func parseReleaseID(r io.Reader) (int64, error) {
	var rr releaseResponse
	if err := json.NewDecoder(r).Decode(&rr); err != nil {
		return 0, err
	}
	return rr.ID, nil
}

func init() {
	Register(NewGitHubPlugin())
}
