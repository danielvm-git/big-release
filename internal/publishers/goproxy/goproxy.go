package goproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/danielvm-git/big-release/internal/publishers"
)

const (
	// DefaultProxyURL is the default Go module mirror URL.
	DefaultProxyURL = "https://proxy.golang.org"

	// maxRetries is the maximum number of retry attempts on 429.
	maxRetries = 3

	// retryBase is the base backoff duration in seconds.
	retryBase = 1 * time.Second

	// maxResponseSize is the maximum response body size in bytes to read (1 MB).
	maxResponseSize = 1 * 1024 * 1024
)

// Publisher publishes Go modules to the Go module mirror.
type Publisher struct {
	// ProxyURL is the Go module proxy URL. Defaults to DefaultProxyURL.
	ProxyURL string
	// HTTPClient is the HTTP client used for API calls. Defaults to http.DefaultClient.
	HTTPClient *http.Client
	// DryRun, when true, skips actual exec and HTTP calls.
	DryRun bool
	// ExecCommand is the function used to run external commands. Defaults to exec.Command.
	ExecCommand func(name string, arg ...string) *exec.Cmd
}

// NewPublisher creates a new Go Proxy Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		ProxyURL:    DefaultProxyURL,
		HTTPClient:  http.DefaultClient,
		ExecCommand: exec.Command,
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "goproxy"
}

// Detect returns true when go.mod exists in the working directory.
func (p *Publisher) Detect() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

// Prepare returns nil — Go proxy uses tag-based versioning, no file mutation needed.
func (p *Publisher) Prepare(version string) error {
	return nil
}

// Publish pushes a versioned git tag, triggers the Go proxy via go list -m,
// and waits for the proxy to acknowledge the new version with retry on 429.
func (p *Publisher) Publish(version string) error {
	modulePath, err := readModulePath()
	if err != nil {
		return err
	}

	if p.DryRun {
		return nil
	}

	// Step 1: Create git tag.
	tag := "v" + version
	cmd := p.ExecCommand("git", "tag", tag)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goproxy: failed to create git tag %s: %w: %s", tag, err, strings.TrimSpace(stderr.String()))
	}

	// Step 2: Push git tag.
	cmd = p.ExecCommand("git", "push", "origin", tag)
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goproxy: failed to push git tag %s: %w: %s", tag, err, strings.TrimSpace(stderr.String()))
	}

	// Step 3: Trigger Go proxy via go list -m.
	proxyURL := os.Getenv("GOPROXY")
	if proxyURL == "" {
		proxyURL = p.ProxyURL
	}

	goModule := fmt.Sprintf("%s@%s", modulePath, version)
	cmd = p.ExecCommand("go", "list", "-m", goModule)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPROXY=%s", proxyURL))
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goproxy: failed to list module %s: %w: %s", goModule, err, strings.TrimSpace(stderr.String()))
	}

	// Step 4: Poll proxy to confirm the version is available (with retry on 429).
	verifyURL := fmt.Sprintf("%s/%s/@v/%s.info", proxyURL, modulePath, version)

	var lastErr error
	backoff := retryBase

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 2)
		}

		req, reqErr := http.NewRequest(http.MethodGet, verifyURL, nil)
		if reqErr != nil {
			lastErr = fmt.Errorf("goproxy: failed to create verify request: %w", reqErr)
			continue
		}

		resp, doErr := p.HTTPClient.Do(req)
		if doErr != nil {
			lastErr = fmt.Errorf("goproxy: verify request failed: %w", doErr)
			continue
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		switch {
		case resp.StatusCode == http.StatusOK:
			return nil
		case resp.StatusCode == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("goproxy: rate limited (HTTP %d)", resp.StatusCode)
			continue
		case resp.StatusCode >= 500:
			return fmt.Errorf("goproxy: server error (HTTP %d)", resp.StatusCode)
		default:
			return fmt.Errorf("goproxy: unexpected proxy status (HTTP %d)", resp.StatusCode)
		}
	}

	return fmt.Errorf("goproxy: publish verification failed after %d retries: %w", maxRetries, lastErr)
}

// Verify checks that the given version is available on the Go module proxy.
func (p *Publisher) Verify(version string) error {
	modulePath, err := readModulePath()
	if err != nil {
		return err
	}

	proxyURL := os.Getenv("GOPROXY")
	if proxyURL == "" {
		proxyURL = p.ProxyURL
	}

	verifyURL := fmt.Sprintf("%s/%s/@v/%s.info", proxyURL, modulePath, version)
	req, err := http.NewRequest(http.MethodGet, verifyURL, nil)
	if err != nil {
		return fmt.Errorf("goproxy: failed to create verify request: %w", err)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("goproxy: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("goproxy: module %q version %q not found (HTTP 404)", modulePath, version)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("goproxy: verify failed with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
	if err != nil {
		return fmt.Errorf("goproxy: failed to read verify response: %w", err)
	}
	if len(body) > maxResponseSize {
		return fmt.Errorf("goproxy: response body too large")
	}

	var info struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return fmt.Errorf("goproxy: failed to parse verify response: %w", err)
	}

	if info.Version != version {
		return fmt.Errorf("goproxy: version mismatch: expected %q, got %q", version, info.Version)
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

// readModulePath reads the module path from go.mod.
func readModulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("goproxy: failed to read go.mod: %w", err)
	}

	// Parse the first line: `module <path>`
	lines := strings.SplitN(string(data), "\n", 2)
	if len(lines) == 0 {
		return "", fmt.Errorf("goproxy: go.mod is empty")
	}

	parts := strings.Fields(lines[0])
	if len(parts) < 2 || parts[0] != "module" {
		return "", fmt.Errorf("goproxy: could not parse module path from go.mod")
	}

	return parts[1], nil
}
