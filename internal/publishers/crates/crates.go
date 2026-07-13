package crates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/danielvm-git/big-release/internal/publishers"
)

const (
	// DefaultRegistryURL is the default crates.io API endpoint.
	DefaultRegistryURL = "https://crates.io/api/v1/crates/new"

	// verifyBaseURL is the base URL for crates.io's versions API.
	verifyBaseURL = "https://crates.io/api/v1/crates"

	// maxRetries is the maximum number of retry attempts on 429.
	maxRetries = 3

	// retryBase is the base backoff duration in seconds.
	retryBase = 1 * time.Second

	// envToken is the environment variable for the crates.io token.
	envToken = "CARGO_TOKEN"
)

// Publisher publishes Rust crates to crates.io.
type Publisher struct {
	// RegistryURL is the crates.io upload endpoint. Defaults to DefaultRegistryURL.
	RegistryURL string
	// HTTPClient is the HTTP client used for API calls. Defaults to http.DefaultClient.
	HTTPClient *http.Client
	// DryRun, when true, skips actual HTTP requests.
	DryRun bool
	// VerifyURL is the base URL for crates.io's versions API. Defaults to verifyBaseURL.
	VerifyURL string
}

// NewPublisher creates a new crates.io Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		RegistryURL: DefaultRegistryURL,
		HTTPClient:  http.DefaultClient,
		VerifyURL:   verifyBaseURL,
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "crates"
}

// Detect returns true when Cargo.toml exists in the working directory.
func (p *Publisher) Detect() bool {
	_, err := os.Stat("Cargo.toml")
	return err == nil
}

// Prepare updates the version field in Cargo.toml using TOML-compliant parsing.
func (p *Publisher) Prepare(version string) error {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return fmt.Errorf("crates: failed to read Cargo.toml: %w", err)
	}

	var cfg map[string]interface{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("crates: failed to parse Cargo.toml: %w", err)
	}

	pkg, ok := cfg["package"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("crates: no [package] section found in Cargo.toml")
	}

	pkg["version"] = version

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("crates: failed to marshal Cargo.toml: %w", err)
	}

	if err := os.WriteFile("Cargo.toml", buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("crates: failed to write Cargo.toml: %w", err)
	}

	return nil
}

// Publish uploads the crate to the crates.io registry.
func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("crates: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	req, err := http.NewRequest(http.MethodPut, p.RegistryURL, nil)
	if err != nil {
		return fmt.Errorf("crates: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", token)

	var lastErr error
	backoff := retryBase

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 2)
		}

		resp, doErr := p.HTTPClient.Do(req)
		if doErr != nil {
			lastErr = fmt.Errorf("crates: request failed: %w", doErr)
			continue
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		switch {
		case resp.StatusCode == http.StatusOK:
			return nil
		case resp.StatusCode == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("crates: rate limited (HTTP %d)", resp.StatusCode)
			continue
		case resp.StatusCode == http.StatusUnauthorized:
			return fmt.Errorf("crates: authentication failed (HTTP %d): check CARGO_TOKEN", resp.StatusCode)
		case resp.StatusCode == http.StatusForbidden:
			return fmt.Errorf("crates: forbidden (HTTP %d): insufficient permissions", resp.StatusCode)
		case resp.StatusCode >= 500:
			return fmt.Errorf("crates: server error (HTTP %d)", resp.StatusCode)
		default:
			return fmt.Errorf("crates: unexpected status (HTTP %d)", resp.StatusCode)
		}
	}

	return fmt.Errorf("crates: publish failed after %d retries: %w", maxRetries, lastErr)
}

// Verify checks that the given version has been published to crates.io.
func (p *Publisher) Verify(version string) error {
	pkgName, err := readPackageName()
	if err != nil {
		return fmt.Errorf("crates: %w", err)
	}

	url := fmt.Sprintf("%s/%s/versions", p.VerifyURL, pkgName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("crates: failed to create verify request: %w", err)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("crates: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("crates: crate %q not found (HTTP 404)", pkgName)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("crates: verify failed with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("crates: failed to read verify response: %w", err)
	}

	var versionsResp struct {
		Versions []struct {
			Num string `json:"num"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(body, &versionsResp); err != nil {
		return fmt.Errorf("crates: failed to parse verify response: %w", err)
	}

	for _, v := range versionsResp.Versions {
		if v.Num == version {
			return nil
		}
	}

	return fmt.Errorf("crates: version %q not found in published versions", version)
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

func init() {
	publishers.Register(NewPublisher())
}

// --- internal helpers ---

// readPackageName reads the package name from Cargo.toml.
func readPackageName() (string, error) {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return "", fmt.Errorf("failed to read Cargo.toml: %w", err)
	}

	var cfg struct {
		Package struct {
			Name string `toml:"name"`
		} `toml:"package"`
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse Cargo.toml: %w", err)
	}

	if cfg.Package.Name == "" {
		return "", fmt.Errorf("package name not found in Cargo.toml")
	}

	return cfg.Package.Name, nil
}
