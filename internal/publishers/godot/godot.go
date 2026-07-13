package godot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/danielvm-git/big-release/internal/publishers"
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

	// maxRetries is the maximum number of retry attempts on 429.
	maxRetries = 3

	// retryBase is the base backoff duration in seconds.
	retryBase = 1 * time.Second
)

// versionKeyPattern matches config/version in project.godot INI-style files.
var versionKeyPattern = regexp.MustCompile(`^(config/version\s*=\s*)"?([^"\n]+)"?\s*$`)

// Publisher publishes Godot projects via GitHub Releases.
type Publisher struct {
	// GitHubAPI is the GitHub API base URL. Defaults to DefaultGitHubAPI.
	GitHubAPI string
	// HTTPClient is the HTTP client used for API calls. Defaults to http.DefaultClient.
	HTTPClient *http.Client
	// DryRun, when true, skips actual HTTP requests.
	DryRun bool
}

// NewPublisher creates a new Godot Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		GitHubAPI:  DefaultGitHubAPI,
		HTTPClient: http.DefaultClient,
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

	// Build release payload.
	payload := map[string]interface{}{
		"tag_name":         "v" + version,
		"name":             "v" + version,
		"target_commitish": "main",
		"draft":            false,
		"prerelease":       false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("godot: failed to marshal release payload: %w", err)
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases", p.GitHubAPI, owner, repo)

	var lastErr error
	backoff := retryBase

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(
				float64(backoff)*2,
				float64(retryBase*time.Duration(math.Pow(2, float64(maxRetries)))),
			))
		}

		req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if reqErr != nil {
			lastErr = fmt.Errorf("godot: failed to create request: %w", reqErr)
			continue
		}
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, doErr := p.HTTPClient.Do(req)
		if doErr != nil {
			lastErr = fmt.Errorf("godot: request failed: %w", doErr)
			continue
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		switch {
		case resp.StatusCode == http.StatusCreated:
			return nil
		case resp.StatusCode == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("godot: rate limited (HTTP %d)", resp.StatusCode)
			continue
		case resp.StatusCode == http.StatusUnauthorized:
			return fmt.Errorf("godot: authentication failed (HTTP %d): check GITHUB_TOKEN", resp.StatusCode)
		case resp.StatusCode >= 500:
			return fmt.Errorf("godot: server error (HTTP %d)", resp.StatusCode)
		default:
			return fmt.Errorf("godot: unexpected status (HTTP %d)", resp.StatusCode)
		}
	}

	return fmt.Errorf("godot: publish failed after %d retries: %w", maxRetries, lastErr)
}

// Verify checks that a GitHub Release exists for the given version tag.
func (p *Publisher) Verify(version string) error {
	owner := os.Getenv(envOwner)
	repo := os.Getenv(envRepo)
	if owner == "" || repo == "" {
		return fmt.Errorf("godot: %s and %s environment variables must be set", envOwner, envRepo)
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/v%s", p.GitHubAPI, owner, repo, version)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("godot: failed to create verify request: %w", err)
	}

	token := os.Getenv(envToken)
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := p.HTTPClient.Do(req)
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

func init() {
	publishers.Register(NewPublisher())
}

// --- internal helpers ---

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
