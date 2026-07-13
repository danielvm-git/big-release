// story: e03s02
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// GitHubPlugin creates GitHub releases via the API.
type GitHubPlugin struct {
	client     HTTPClient
	apiBaseURL string // for testing; empty = use default GitHub API
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

// Name returns the plugin name.
func (p *GitHubPlugin) Name() string {
	return "github"
}

// VerifyConditions checks that required environment variables are set.
func (p *GitHubPlugin) VerifyConditions(ctx *algorithm.Context) error {
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
func (p *GitHubPlugin) AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error) {
	return "", nil
}

// GenerateNotes is not applicable for the GitHub plugin.
func (p *GitHubPlugin) GenerateNotes(ctx *algorithm.Context) (string, error) {
	return "", nil
}

// Prepare is not applicable for the GitHub plugin.
func (p *GitHubPlugin) Prepare(ctx *algorithm.Context) error {
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

func (p *GitHubPlugin) handleReleaseResponse(resp *http.Response, repo, version string) error {
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("GitHub authentication failed: invalid or expired GITHUB_TOKEN (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusForbidden:
		return fmt.Errorf("GitHub permission denied (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusNotFound:
		return fmt.Errorf("GitHub repository %q not found (HTTP %d): %s", repo, resp.StatusCode, string(respBody))
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("GitHub release already exists for tag %s (HTTP %d): %s", version, resp.StatusCode, string(respBody))
	default:
		return fmt.Errorf("GitHub API returned unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
}

func (p *GitHubPlugin) releaseURL(repo string) string {
	baseURL := p.apiBaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return fmt.Sprintf("%s/repos/%s/releases", baseURL, repo)
}

func (p *GitHubPlugin) sendReleaseRequest(payload []byte, repo, version string) error {
	req, err := p.createHTTPRequest(p.releaseURL(repo), payload)
	if err != nil {
		return err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create GitHub release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return p.handleReleaseResponse(resp, repo, version)
}

// Publish creates a GitHub release via the API.
func (p *GitHubPlugin) Publish(ctx *algorithm.Context) (*algorithm.Release, error) {
	if ctx.DryRun {
		return nil, nil
	}
	payload, err := p.buildReleasePayload(ctx.NextRelease.Version, ctx.NextRelease.Notes, string(ctx.NextRelease.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal release body: %w", err)
	}
	repo := os.Getenv("GITHUB_REPOSITORY")
	if err := p.sendReleaseRequest(payload, repo, ctx.NextRelease.Version); err != nil {
		return nil, err
	}
	return nil, nil
}

// Success is called after a successful release.
func (p *GitHubPlugin) Success(ctx *algorithm.Context) error {
	return nil
}

// Fail is called on release failure.
func (p *GitHubPlugin) Fail(ctx *algorithm.Context, err error) error {
	return nil
}

func init() {
	Register(NewGitHubPlugin())
}
