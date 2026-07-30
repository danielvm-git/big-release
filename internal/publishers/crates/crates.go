package crates

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/httputil"
)

const (
	// DefaultRegistryURL is the default crates.io API endpoint.
	DefaultRegistryURL = "https://crates.io/api/v1/crates/new"

	// verifyBaseURL is the base URL for crates.io's versions API.
	verifyBaseURL = "https://crates.io/api/v1/crates"

	// envToken is the environment variable for the crates.io token.
	envToken = "CARGO_TOKEN"
)

// Publisher publishes Rust crates to crates.io.
type Publisher struct {
	// RegistryURL is the crates.io upload endpoint. Defaults to DefaultRegistryURL.
	RegistryURL string
	// Client is the retry HTTP client used for API calls.
	Client *httputil.RetryClient
	// DryRun, when true, skips actual HTTP requests.
	DryRun bool
	// VerifyURL is the base URL for crates.io's versions API. Defaults to verifyBaseURL.
	VerifyURL string
}

// NewPublisher creates a new crates.io Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		RegistryURL: DefaultRegistryURL,
		Client:      httputil.NewRetryClient(http.DefaultClient),
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
// The crates.io new-crate endpoint requires a tar.gz body containing the
// crate source. Without it the API returns 400 Bad Request (BUG-crates-empty-body-publish).
func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("crates: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	// Build the tar.gz body that crates.io expects.
	body, err := p.buildUploadBody()
	if err != nil {
		return fmt.Errorf("crates: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, p.RegistryURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("crates: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/gzip")

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("crates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain body to allow connection reuse.
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

// buildUploadBody creates a minimal tar.gz archive containing Cargo.toml
// and a placeholder lib.rs — the minimum crates.io requires for a new-crate upload.
func (p *Publisher) buildUploadBody() (*bytes.Buffer, error) {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to read Cargo.toml for upload: %w", err)
	}

	// Read package name from Cargo.toml for the lib.rs placeholder.
	cfg := struct {
		Package struct {
			Name string `toml:"name"`
		} `toml:"package"`
	}{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse Cargo.toml for upload: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	// Write Cargo.toml
	if err := writeTarEntry(tw, "Cargo.toml", data); err != nil {
		return nil, err
	}

	// Write minimal src/lib.rs
	libContent := []byte(fmt.Sprintf("//! %s\n", cfg.Package.Name))
	if err := writeTarEntry(tw, "src/lib.rs", libContent); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return &buf, nil
}

func writeTarEntry(tw *tar.Writer, name string, content []byte) error {
	hdr := &tar.Header{
		Name: name,
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("failed to write tar header for %s: %w", name, err)
	}
	if _, err := tw.Write(content); err != nil {
		return fmt.Errorf("failed to write tar content for %s: %w", name, err)
	}
	return nil
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

	resp, err := p.Client.Do(req)
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
