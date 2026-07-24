// story: e23s01 e23s02 e23s03
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// GitLabConfig configures the gitlab plugin. Loaded from the
// PluginConfigs["gitlab"] entry.
type GitLabConfig struct {
	Assets          []algorithm.AssetConfig `yaml:"assets,omitempty" json:"assets,omitempty"`
	ReleaseName     string                  `yaml:"releaseName,omitempty" json:"releaseName,omitempty"`
	ReleaseBody     string                  `yaml:"releaseBody,omitempty" json:"releaseBody,omitempty"`
	GitlabURL       string                  `yaml:"gitlabUrl,omitempty" json:"gitlabUrl,omitempty"`
	GitlabProjectID string                  `yaml:"gitlabProjectId,omitempty" json:"gitlabProjectId,omitempty"`
}

// GitLabPlugin creates GitLab releases via the API.
type GitLabPlugin struct {
	client              HTTPClient
	apiBaseURL          string // for testing; empty = use CI_API_V4_URL or default
	assets              []algorithm.AssetConfig
	releaseNameTemplate string
	releaseBodyTemplate string
	gitlabURL           string
	gitlabProjectID     string
}

// NewGitLabPlugin creates a new GitLabPlugin.
func NewGitLabPlugin() *GitLabPlugin {
	return &GitLabPlugin{
		client: http.DefaultClient,
	}
}

// Configure decodes the gitlab plugin's typed config from the raw
// PluginConfigs map entry. Implements ConfigurablePlugin.
func (p *GitLabPlugin) Configure(raw map[string]interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("gitlab plugin: failed to re-marshal config: %w", err)
	}
	var cfg GitLabConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("gitlab plugin: invalid config: %w", err)
	}
	p.assets = cfg.Assets
	p.releaseNameTemplate = cfg.ReleaseName
	p.releaseBodyTemplate = cfg.ReleaseBody
	p.gitlabURL = cfg.GitlabURL
	p.gitlabProjectID = cfg.GitlabProjectID
	return nil
}

// Name returns the plugin name.
func (p *GitLabPlugin) Name() string {
	return "gitlab"
}

// VerifyConditions checks that required environment variables are set.
func (p *GitLabPlugin) VerifyConditions(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun {
		return nil
	}
	if gitlabToken() == "" {
		return fmt.Errorf("GITLAB_TOKEN or CI_JOB_TOKEN environment variable is required")
	}
	if p.projectID() == "" {
		return fmt.Errorf("CI_PROJECT_ID or GITLAB_PROJECT_ID environment variable is required")
	}
	return nil
}

// AnalyzeCommits is not applicable for the GitLab plugin.
func (p *GitLabPlugin) AnalyzeCommits(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (algorithm.ReleaseType, error) {
	return "", nil
}

// VerifyRelease is not applicable for the GitLab plugin.
func (p *GitLabPlugin) VerifyRelease(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// GenerateNotes is not applicable for the GitLab plugin.
func (p *GitLabPlugin) GenerateNotes(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (string, error) {
	return "", nil
}

// Prepare is not applicable for the GitLab plugin.
func (p *GitLabPlugin) Prepare(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

type gitlabReleaseRequest struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
	Ref         string `json:"ref,omitempty"`
}

type gitlabAssetLinkRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	LinkType string `json:"link_type,omitempty"`
}

type gitlabUploadResponse struct {
	URL      string `json:"url"`
	FullPath string `json:"full_path"`
}

func (p *GitLabPlugin) buildReleasePayload(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) ([]byte, error) {
	version := state.NextRelease.Version
	notes := state.NextRelease.Notes

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

	req := gitlabReleaseRequest{
		Name:        name,
		TagName:     version,
		Description: body,
		Ref:         p.releaseRef(ctx),
	}
	return json.Marshal(req)
}

func (p *GitLabPlugin) renderTemplate(label, tmpl, fallback string, tctx templateContext) (string, error) {
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

func (p *GitLabPlugin) releaseRef(ctx *algorithm.ReadOnlyContext) string {
	if ref := os.Getenv("CI_COMMIT_REF_NAME"); ref != "" {
		return ref
	}
	if ctx.Branch != nil {
		return ctx.Branch.Name
	}
	return "main"
}

func (p *GitLabPlugin) apiBase() string {
	if p.apiBaseURL != "" {
		return strings.TrimSuffix(p.apiBaseURL, "/")
	}
	if base := os.Getenv("CI_API_V4_URL"); base != "" {
		return strings.TrimSuffix(base, "/")
	}
	if p.gitlabURL != "" {
		return strings.TrimSuffix(p.gitlabURL, "/") + "/api/v4"
	}
	return "https://gitlab.com/api/v4"
}

func (p *GitLabPlugin) projectID() string {
	if p.gitlabProjectID != "" {
		return p.gitlabProjectID
	}
	if id := os.Getenv("CI_PROJECT_ID"); id != "" {
		return id
	}
	return os.Getenv("GITLAB_PROJECT_ID")
}

func gitlabToken() string {
	if t := os.Getenv("GITLAB_TOKEN"); t != "" {
		return t
	}
	return os.Getenv("CI_JOB_TOKEN")
}

func usesJobToken() bool {
	return os.Getenv("GITLAB_TOKEN") == "" && os.Getenv("CI_JOB_TOKEN") != ""
}

func (p *GitLabPlugin) setAuth(req *http.Request) {
	if usesJobToken() {
		req.Header.Set("Job-Token", os.Getenv("CI_JOB_TOKEN"))
		return
	}
	req.Header.Set("Private-Token", os.Getenv("GITLAB_TOKEN"))
}

func (p *GitLabPlugin) releasesURL() string {
	return fmt.Sprintf("%s/projects/%s/releases", p.apiBase(), p.projectID())
}

func (p *GitLabPlugin) createHTTPRequest(url string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	p.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (p *GitLabPlugin) handleReleaseResponse(resp *http.Response, version string) error {
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("GitLab authentication failed: invalid or expired token (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusForbidden:
		return fmt.Errorf("GitLab permission denied (HTTP %d): %s", resp.StatusCode, string(respBody))
	case http.StatusNotFound:
		return fmt.Errorf("GitLab project %q not found (HTTP %d): %s", p.projectID(), resp.StatusCode, string(respBody))
	case http.StatusConflict:
		return fmt.Errorf("GitLab release already exists for tag %s (HTTP %d): %s", version, resp.StatusCode, string(respBody))
	default:
		return fmt.Errorf("GitLab API returned unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
}

func (p *GitLabPlugin) sendReleaseRequest(payload []byte, version string) error {
	req, err := p.createHTTPRequest(p.releasesURL(), payload)
	if err != nil {
		return err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create GitLab release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return p.handleReleaseResponse(resp, version)
}

// Publish creates a GitLab release via the API and uploads configured assets.
func (p *GitLabPlugin) Publish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) (*algorithm.Release, error) {
	if ctx.DryRun {
		return nil, nil
	}
	payload, err := p.buildReleasePayload(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal release body: %w", err)
	}
	version := state.NextRelease.Version
	if err := p.sendReleaseRequest(payload, version); err != nil {
		return nil, err
	}

	if len(p.assets) > 0 {
		if err := p.uploadAssets(version); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (p *GitLabPlugin) uploadAssets(tagName string) error {
	assets, errs := expandAssetGlobs(p.assets)
	for _, errMsg := range errs {
		fmt.Fprintf(os.Stderr, "warning: %s\n", errMsg)
	}

	for _, asset := range assets {
		if err := p.uploadOneAsset(tagName, asset); err != nil {
			return err
		}
	}
	return nil
}

func (p *GitLabPlugin) uploadOneAsset(tagName string, asset algorithm.AssetConfig) error {
	f, err := os.Open(asset.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not open asset %q: %v\n", asset.Path, err)
		return nil
	}
	defer func() { _ = f.Close() }()

	name := asset.Label
	if name == "" {
		name = filepath.Base(asset.Path)
	}

	uploadURL, err := p.uploadFile(f, name)
	if err != nil {
		return err
	}
	return p.createAssetLink(tagName, name, uploadURL)
}

func (p *GitLabPlugin) uploadFile(r io.Reader, filename string) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create upload form: %w", err)
	}
	if _, err := io.Copy(part, r); err != nil {
		return "", fmt.Errorf("failed to copy asset into form: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize upload form: %w", err)
	}

	url := fmt.Sprintf("%s/projects/%s/uploads", p.apiBase(), p.projectID())
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}
	p.setAuth(req)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload asset %q: %w", filename, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("asset upload %q failed (HTTP %d): %s", filename, resp.StatusCode, string(respBody))
	}

	var upload gitlabUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&upload); err != nil {
		return "", fmt.Errorf("failed to decode upload response: %w", err)
	}
	return p.fullAssetURL(upload.URL), nil
}

func (p *GitLabPlugin) fullAssetURL(uploadPath string) string {
	if strings.HasPrefix(uploadPath, "http://") || strings.HasPrefix(uploadPath, "https://") {
		return uploadPath
	}
	if projectURL := os.Getenv("CI_PROJECT_URL"); projectURL != "" {
		return strings.TrimSuffix(projectURL, "/") + uploadPath
	}
	if p.gitlabURL != "" {
		return strings.TrimSuffix(p.gitlabURL, "/") + uploadPath
	}
	// Fallback: derive web URL from API base (gitlab.com/api/v4 → gitlab.com).
	base := strings.TrimSuffix(p.apiBase(), "/api/v4")
	return base + uploadPath
}

func (p *GitLabPlugin) createAssetLink(tagName, name, assetURL string) error {
	link := gitlabAssetLinkRequest{
		Name:     name,
		URL:      assetURL,
		LinkType: "package",
	}
	payload, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("failed to marshal asset link: %w", err)
	}
	url := fmt.Sprintf("%s/projects/%s/releases/%s/assets/links",
		p.apiBase(), p.projectID(), tagName)
	req, err := p.createHTTPRequest(url, payload)
	if err != nil {
		return err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create asset link %q: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("asset link %q failed (HTTP %d): %s", name, resp.StatusCode, string(respBody))
	}
	return nil
}

// Success is a no-op for the GitLab plugin.
func (p *GitLabPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	return nil
}

// Fail is called on release failure.
func (p *GitLabPlugin) Fail(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState, err error) error {
	return nil
}

func init() {
	Register(NewGitLabPlugin())
}
