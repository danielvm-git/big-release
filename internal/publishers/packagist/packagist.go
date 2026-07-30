package packagist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/httputil"
)

const (
	DefaultAPIURL = "https://packagist.org"

	updateEndpoint = "/api/update-package"

	envToken = "PACKAGIST_TOKEN"
)

type Publisher struct {
	APIURL string
	Client *httputil.RetryClient
	DryRun bool
}

func NewPublisher() *Publisher {
	return &Publisher{
		APIURL: DefaultAPIURL,
		Client: httputil.NewRetryClient(http.DefaultClient),
	}
}

func (p *Publisher) Name() string {
	return "packagist"
}

func (p *Publisher) Detect() bool {
	_, err := os.Stat("composer.json")
	return err == nil
}

func (p *Publisher) Prepare(version string) error {
	data, err := os.ReadFile("composer.json")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("packagist: composer.json not found")
		}
		return fmt.Errorf("packagist: failed to read composer.json: %w", err)
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("packagist: failed to parse composer.json: %w", err)
	}

	cfg["version"] = version

	updated, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("packagist: failed to marshal composer.json: %w", err)
	}

	updated = append(updated, '\n')

	if err := os.WriteFile("composer.json", updated, 0644); err != nil {
		return fmt.Errorf("packagist: failed to write composer.json: %w", err)
	}

	return nil
}

func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("packagist: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	// The repository URL must be the git repo URL, not the Packagist API URL.
	// Use GITHUB_REPOSITORY env var (owner/repo format) to build the GitHub URL,
	// falling back to git remote origin URL (BUG-packagist-wrong-update-url).
	repoURL := p.resolveRepositoryURL()

	body := map[string]interface{}{
		"repository": map[string]string{
			"url": repoURL,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("packagist: failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.APIURL+updateEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("packagist: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("packagist: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain body to allow connection reuse.
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

func (p *Publisher) Verify(version string) error {
	vendor, pkg, err := readPackageVendorAndName()
	if err != nil {
		return fmt.Errorf("packagist: %w", err)
	}

	url := fmt.Sprintf("%s/packages/%s/%s.json", p.APIURL, vendor, pkg)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("packagist: failed to create verify request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("packagist: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("packagist: package %s/%s not found (HTTP 404)", vendor, pkg)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("packagist: verify failed with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("packagist: failed to read verify response: %w", err)
	}

	var pkgInfo struct {
		Package struct {
			Versions map[string]interface{} `json:"versions"`
		} `json:"package"`
	}
	if err := json.Unmarshal(body, &pkgInfo); err != nil {
		return fmt.Errorf("packagist: failed to parse verify response: %w", err)
	}

	if _, ok := pkgInfo.Package.Versions[version]; !ok {
		return fmt.Errorf("packagist: version %q not found in published versions", version)
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

func readPackageVendorAndName() (string, string, error) {
	data, err := os.ReadFile("composer.json")
	if err != nil {
		return "", "", fmt.Errorf("failed to read composer.json: %w", err)
	}

	var cfg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", "", fmt.Errorf("failed to parse composer.json: %w", err)
	}

	if cfg.Name == "" {
		return "", "", fmt.Errorf("package name not found in composer.json")
	}

	parts := splitPackageName(cfg.Name)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid package name %q in composer.json (expected vendor/package)", cfg.Name)
	}

	return parts[0], parts[1], nil
}

func splitPackageName(name string) []string {
	for i := 0; i < len(name); i++ {
		if name[i] == '/' {
			return []string{name[:i], name[i+1:]}
		}
	}
	return []string{name}
}

// resolveRepositoryURL returns the git repository URL for the Packagist
// update request. It tries GITHUB_REPOSITORY (owner/repo format) first,
// then falls back to git remote origin URL (BUG-packagist-wrong-update-url).
func (p *Publisher) resolveRepositoryURL() string {
	if repo := os.Getenv("GITHUB_REPOSITORY"); repo != "" {
		return "https://github.com/" + repo
	}
	// Fallback: read from git remote.
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err == nil {
		url := strings.TrimSpace(string(out))
		if url != "" {
			return url
		}
	}
	return p.APIURL
}
