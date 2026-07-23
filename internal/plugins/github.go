// story: e03s02
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// GitHubPlugin creates GitHub releases via the API.
type GitHubPlugin struct {
	client              HTTPClient
	apiBaseURL          string // for testing; empty = use default GitHub API
	uploadBaseURL       string // for testing; empty = use default uploads host
	assets              []algorithm.AssetConfig
	draftRelease        bool
	releaseNameTemplate string
	releaseBodyTemplate string
	successComment      string
	releasedLabels      []string
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
	TagName              string `json:"tag_name"`
	Name                 string `json:"name"`
	Body                 string `json:"body"`
	Draft                bool   `json:"draft"`
	Prerelease           bool   `json:"prerelease"`
	GenerateReleaseNotes bool   `json:"generate_release_notes"`
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

	name, err := p.renderTemplate("release name", p.releaseNameTemplate, "v"+version, tctx)
	if err != nil {
		return nil, err
	}
	body, err := p.renderTemplate("release body", p.releaseBodyTemplate, notes, tctx)
	if err != nil {
		return nil, err
	}

	req := &createReleaseRequest{
		TagName:              version,
		Name:                 name,
		Body:                 body,
		Draft:                p.draftRelease,
		Prerelease:           releaseType == "prerelease",
		GenerateReleaseNotes: notes == "",
	}
	return json.Marshal(req)
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

	// Upload configured assets (if any). Missing files warn but do not fail.
	if len(p.assets) > 0 && releaseID != 0 {
		if err := p.uploadAssets(repo, releaseID); err != nil {
			return nil, err
		}
	}

	// Draft releases are published atomically after asset upload (#13).
	if p.draftRelease && releaseID != 0 {
		if err := p.publishDraft(repo, releaseID); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// publishDraft flips a draft release to published via PATCH.
func (p *GitHubPlugin) publishDraft(repo string, releaseID int64) error {
	patchURL := fmt.Sprintf("%s/repos/%s/releases/%d", p.releaseURLBase(), repo, releaseID)
	payload, err := json.Marshal(map[string]bool{"draft": false})
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

// Success is called after a successful release. It posts a comment to
// each issue/PR referenced in the release commits (via fixes/closes/
// resolves #N), so contributors and users get notified (#12). Commenting
// failures (403/404) are logged and non-fatal — they must not fail an
// otherwise-successful release.
func (p *GitHubPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun {
		return nil
	}
	if state.NextRelease == nil || len(ctx.Commits) == 0 {
		return nil
	}

	// Collect referenced issue numbers from all release commit messages.
	seen := map[int]bool{}
	for _, c := range ctx.Commits {
		for _, n := range parseReferencedIssues(c.Message) {
			seen[n] = true
		}
	}
	if len(seen) == 0 {
		return nil
	}

	comment, err := p.buildSuccessComment(state.NextRelease.Version)
	if err != nil {
		// Invalid comment template is non-fatal in the success hook.
		fmt.Fprintf(os.Stderr, "warning: invalid successComment template: %v\n", err)
		return nil
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	for n := range seen {
		if err := p.postIssueComment(repo, n, comment); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not comment on issue #%d: %v\n", n, err)
		}
	}
	return nil
}

// buildSuccessComment renders the success comment template, or falls back
// to a sensible default when no template is configured.
func (p *GitHubPlugin) buildSuccessComment(version string) (string, error) {
	tmpl := p.successComment
	fallback := fmt.Sprintf("🎉 Released in version %s", version)
	return p.renderTemplate("success comment", tmpl, fallback, templateContext{Version: version})
}

// postIssueComment posts a comment to /repos/{repo}/issues/{n}/comments.
// 403/404 responses are returned as errors but are non-fatal at the
// caller (Success logs and continues).
func (p *GitHubPlugin) postIssueComment(repo string, issueNum int, body string) error {
	url := fmt.Sprintf("%s/repos/%s/issues/%d/comments", p.releaseURLBase(), repo, issueNum)
	payload, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create comment request: %w", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	// 403/404 are non-fatal — caller (Success) logs and continues.
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("comment on issue #%d failed (HTTP %d): %s", issueNum, resp.StatusCode, string(respBody))
}

// issueRefPattern captures issue numbers from conventional commit
// messages: bare (#123) or keyword-prefixed (fixes/closes/resolves #N).
var issueRefPattern = regexp.MustCompile(`(?i)(?:fixes|closes|resolves)?\s*#(\d+)`)

// parseReferencedIssues returns the deduplicated, ordered list of issue
// numbers referenced in a commit message via #N, fixes #N, closes #N,
// or resolves #N.
func parseReferencedIssues(message string) []int {
	matches := issueRefPattern.FindAllStringSubmatch(message, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := map[int]bool{}
	var out []int
	for _, m := range matches {
		n, err := strconv.Atoi(m[1])
		if err != nil || n <= 0 || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	return out
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

// uploadAssets expands configured asset globs and uploads each file to
// the given release via the GitHub uploads endpoint. Missing files are
// logged as warnings and skipped (non-fatal).
func (p *GitHubPlugin) uploadAssets(repo string, releaseID int64) error {
	assets, errs := expandAssetGlobs(p.assets)
	// Missing globs/files are warnings, not failures (#10 acceptance).
	for _, errMsg := range errs {
		fmt.Fprintf(os.Stderr, "warning: %s\n", errMsg)
	}

	for _, asset := range assets {
		if err := p.uploadOneAsset(repo, releaseID, asset); err != nil {
			return err
		}
	}
	return nil
}

// uploadOneAsset uploads a single file to the release.
func (p *GitHubPlugin) uploadOneAsset(repo string, releaseID int64, asset algorithm.AssetConfig) error {
	f, err := os.Open(asset.Path)
	if err != nil {
		// Missing file is a warning, not a failure.
		fmt.Fprintf(os.Stderr, "warning: could not open asset %q: %v\n", asset.Path, err)
		return nil
	}
	defer func() { _ = f.Close() }()

	name := asset.Label
	if name == "" {
		name = filepath.Base(asset.Path)
	}
	uploadURL := fmt.Sprintf("%s/repos/%s/releases/%d/assets?name=%s",
		p.uploadsHost(), repo, releaseID, url.QueryEscape(name))

	req, err := http.NewRequest(http.MethodPost, uploadURL, f)
	if err != nil {
		return fmt.Errorf("failed to create asset upload request for %q: %w", name, err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", mimeTypeForAsset(name))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload asset %q: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("asset upload %q failed (HTTP %d): %s", name, resp.StatusCode, string(body))
	}
	return nil
}

// uploadsHost returns the GitHub uploads base URL.
func (p *GitHubPlugin) uploadsHost() string {
	if p.uploadBaseURL != "" {
		return p.uploadBaseURL
	}
	return "https://uploads.github.com"
}

// expandAssetGlobs expands glob patterns in asset paths to concrete files.
// Returns the expanded list plus a slice of error messages for any pattern
// that matched no files (non-fatal — caller may log and continue).
func expandAssetGlobs(assets []algorithm.AssetConfig) ([]algorithm.AssetConfig, []string) {
	var out []algorithm.AssetConfig
	var errs []string

	for _, asset := range assets {
		matches, err := filepath.Glob(asset.Path)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid glob %q: %v", asset.Path, err))
			continue
		}
		if len(matches) == 0 {
			// Could be a literal path that doesn't exist, or an empty glob.
			// Keep it as-is so the uploader can emit a missing-file warning.
			out = append(out, asset)
			continue
		}
		for _, m := range matches {
			out = append(out, algorithm.AssetConfig{Path: m, Label: asset.Label})
		}
	}
	return out, errs
}

// mimeTypeForAsset returns the MIME type for an asset based on its
// filename extension. Falls back to application/octet-stream.
func mimeTypeForAsset(name string) string {
	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".gz", ".tgz":
		return "application/gzip"
	case ".zip":
		return "application/zip"
	case ".exe":
		return "application/vnd.microsoft.portable-executable"
	case ".dmg":
		return "application/x-apple-diskimage"
	case ".deb":
		return "application/vnd.debian.binary-package"
	case ".rpm":
		return "application/x-rpm"
	default:
		if t := mime.TypeByExtension(ext); t != "" {
			return t
		}
		return "application/octet-stream"
	}
}

func init() {
	Register(NewGitHubPlugin())
}
