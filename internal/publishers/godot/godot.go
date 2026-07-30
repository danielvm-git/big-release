package godot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/httputil"
)

const (
	// DefaultGitHubAPI is the default GitHub API endpoint.
	DefaultGitHubAPI = "https://api.github.com"

	// envToken is the environment variable for the GitHub token.
	envToken = "GITHUB_TOKEN"

	// envOwner is the environment variable for the GitHub owner.
	envOwner = "GITHUB_OWNER"

	// envRepo is the environment variable for the GitHub repo.
	envRepo = "GITHUB_REPO"
)

// versionKeyPattern matches config/version in project.godot INI-style files.
var versionKeyPattern = regexp.MustCompile(`^(config/version\s*=\s*)"?([^"\n]+)"?\s*$`)

// Publisher publishes Godot projects via GitHub Releases.
type Publisher struct {
	// GitHubAPI is the GitHub API base URL. Defaults to DefaultGitHubAPI.
	GitHubAPI string
	// Client is the retry HTTP client used for API calls.
	Client *httputil.RetryClient
	// DryRun, when true, skips actual HTTP requests.
	DryRun bool
}

// NewPublisher creates a new Godot Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		GitHubAPI: DefaultGitHubAPI,
		Client:    httputil.NewRetryClient(http.DefaultClient),
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "godot"
}

// Detect returns true when project.godot exists in the working directory.
func (p *Publisher) Detect() bool {
	_, err := os.Stat("project.godot")
	return err == nil
}

// Prepare updates the config/version key in project.godot.
func (p *Publisher) Prepare(version string) error {
	if _, err := os.Stat("project.godot"); err != nil {
		return fmt.Errorf("godot: project.godot not found: %w", err)
	}
	return updateVersionInProjectGodot(version)
}

// Publish creates a GitHub Release for the given version.
func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("godot: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	owner := os.Getenv(envOwner)
	repo := os.Getenv(envRepo)
	if owner == "" || repo == "" {
		return fmt.Errorf("godot: %s and %s environment variables must be set", envOwner, envRepo)
	}

	// Build release payload. Use the git tagFormat for the tag name
	// instead of hardcoding "v" prefix (BUG-godot-hardcoded-v-prefix).
	tag := p.resolveTag(version)
	payload := map[string]interface{}{
		"tag_name":         tag,
		"name":             tag,
		"target_commitish": "main",
		"draft":            false,
		"prerelease":       false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("godot: failed to marshal release payload: %w", err)
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases", p.GitHubAPI, owner, repo)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("godot: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("godot: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain body to allow connection reuse.
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

// Verify checks that a GitHub Release exists for the given version tag.
func (p *Publisher) Verify(version string) error {
	owner := os.Getenv(envOwner)
	repo := os.Getenv(envRepo)
	if owner == "" || repo == "" {
		return fmt.Errorf("godot: %s and %s environment variables must be set", envOwner, envRepo)
	}

	tag := p.resolveTag(version)
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", p.GitHubAPI, owner, repo, tag)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("godot: failed to create verify request: %w", err)
	}

	token := os.Getenv(envToken)
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("godot: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("godot: release for version %q not found (HTTP 404)", version)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("godot: verify failed with HTTP %d", resp.StatusCode)
	}

	return nil
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

func init() {
	publishers.Register(NewPublisher())
}

// --- internal helpers ---

// resolveTag returns the git tag for the given version. It reads the
// tagFormat from the .big-release.yml config, falling back to "v" prefix
// (BUG-godot-hardcoded-v-prefix).
func (p *Publisher) resolveTag(version string) string {
	// Try to read tagFormat from config.
	cfgPath := ".big-release.yml"
	if _, err := os.Stat(cfgPath); err == nil {
		data, err := os.ReadFile(cfgPath)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "tagFormat:") {
					parts := strings.SplitN(trimmed, ":", 2)
					if len(parts) == 2 {
						fmt := strings.TrimSpace(parts[1])
						fmt = strings.Trim(fmt, `"'`)
						if idx := strings.Index(fmt, "${version}"); idx > 0 {
							return fmt[:idx] + version
						}
					}
				}
			}
		}
	}
	return "v" + version
}

// updateVersionInProjectGodot reads project.godot, finds config/version,
// replaces it, and writes the file back.
func updateVersionInProjectGodot(version string) error {
	data, err := os.ReadFile("project.godot")
	if err != nil {
		return fmt.Errorf("godot: failed to read project.godot: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	found := false

	for i, line := range lines {
		if versionKeyPattern.MatchString(line) {
			matches := versionKeyPattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				prefix := matches[1]
				lines[i] = prefix + `"` + version + `"`
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("godot: config/version key not found in project.godot")
	}

	output := strings.Join(lines, "\n")
	if err := os.WriteFile("project.godot", []byte(output), 0644); err != nil {
		return fmt.Errorf("godot: failed to write project.godot: %w", err)
	}

	return nil
}
