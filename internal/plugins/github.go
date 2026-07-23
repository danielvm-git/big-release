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
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// GitHubPlugin creates GitHub releases via the API.
type GitHubPlugin struct {
	client        HTTPClient
	apiBaseURL    string // for testing; empty = use default GitHub API
	uploadBaseURL string // for testing; empty = use default uploads host
	assets        []algorithm.AssetConfig
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
	Prerelease           bool   `json:"prerelease"`
	GenerateReleaseNotes bool   `json:"generate_release_notes"`
}

func (p *GitHubPlugin) buildReleasePayload(version, notes, releaseType string) ([]byte, error) {
	body := &createReleaseRequest{
		TagName:              version,
		Name:                 fmt.Sprintf("v%s", version),
		Body:                 notes,
		Prerelease:           releaseType == "prerelease",
		GenerateReleaseNotes: notes == "",
	}
	return json.Marshal(body)
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
	baseURL := p.apiBaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return fmt.Sprintf("%s/repos/%s/releases", baseURL, repo)
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
	payload, err := p.buildReleasePayload(state.NextRelease.Version, state.NextRelease.Notes, string(state.NextRelease.Type))
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
	return nil, nil
}

// Success is called after a successful release.
func (p *GitHubPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
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
