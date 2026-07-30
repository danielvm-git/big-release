package pypi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/httputil"
)

const (
	// DefaultRegistryURL is the default PyPI upload endpoint.
	DefaultRegistryURL = "https://upload.pypi.org/legacy/"

	// verifyBaseURL is the base URL for PyPI's JSON API.
	verifyBaseURL = "https://pypi.org/pypi"

	// envToken is the environment variable for the PyPI token.
	envToken = "PYPI_TOKEN"
)

// versionLinePattern matches version assignments in setup.cfg (version = X.Y.Z)
// and pyproject.toml (version = "X.Y.Z").
var versionLinePattern = regexp.MustCompile(`^(version\s*=\s*)["']?([^"'\s#]+)["']?\s*(#.*)?$`)

// Publisher publishes Python packages to PyPI.
type Publisher struct {
	// RegistryURL is the PyPI upload endpoint. Defaults to DefaultRegistryURL.
	RegistryURL string
	// Client is the retry HTTP client used for API calls.
	Client *httputil.RetryClient
	// DryRun, when true, skips actual HTTP requests.
	DryRun bool
	// VerifyURL is the base URL for PyPI's JSON API. Defaults to verifyBaseURL.
	VerifyURL string
}

// NewPublisher creates a new PyPI Publisher with default settings.
func NewPublisher() *Publisher {
	return &Publisher{
		RegistryURL: DefaultRegistryURL,
		Client:      httputil.NewRetryClient(http.DefaultClient),
		VerifyURL:   verifyBaseURL,
	}
}

// Name returns the publisher name.
func (p *Publisher) Name() string {
	return "pypi"
}

// Detect returns true when setup.py or pyproject.toml exists in the working
// directory.
func (p *Publisher) Detect() bool {
	if _, err := os.Stat("setup.py"); err == nil {
		return true
	}
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return true
	}
	return false
}

// Prepare updates the version field in the project configuration file.
// It updates setup.cfg or pyproject.toml depending on which exists.
// pyproject.toml takes precedence when both exist.
func (p *Publisher) Prepare(version string) error {
	// Try pyproject.toml first (modern Python packaging prefers it)
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return updateVersionInFile("pyproject.toml", version)
	}
	// Fall back to setup.cfg
	if _, err := os.Stat("setup.cfg"); err == nil {
		return updateVersionInFile("setup.cfg", version)
	}
	// Fall back to setup.py as last resort
	if _, err := os.Stat("setup.py"); err == nil {
		return nil // setup.py version is embedded in code, skip updating
	}
	return fmt.Errorf("pypi: no supported config file found (pyproject.toml, setup.cfg, setup.py)")
}

// Publish uploads the package to the PyPI registry.
func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("pypi: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	// Build multipart form data for PyPI upload.
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// PyPI requires the distribution files to be attached.
	// We look for dist/*.tar.gz and dist/*.whl files.
	distFiles, err := filepath.Glob("dist/*")
	if err != nil {
		return fmt.Errorf("pypi: failed to list dist/ directory: %w", err)
	}
	if len(distFiles) == 0 {
		return fmt.Errorf("pypi: no distribution files found in dist/ (run python -m build first)")
	}

	for _, f := range distFiles {
		if err := p.addDistFile(w, f); err != nil {
			return err
		}
	}

	// Add required protocol version field.
	if err := w.WriteField("protocol_version", "1"); err != nil {
		return fmt.Errorf("pypi: failed to write protocol_version: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("pypi: failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.RegistryURL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf("pypi: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "token "+token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("pypi: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain body to allow connection reuse.
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

// Verify checks that the given version has been published to PyPI.
func (p *Publisher) Verify(version string) error {
	pkgName, err := readPackageName()
	if err != nil {
		return fmt.Errorf("pypi: %w", err)
	}

	url := fmt.Sprintf("%s/%s/json", p.VerifyURL, pkgName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("pypi: failed to create verify request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("pypi: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("pypi: package %q not found (HTTP 404)", pkgName)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pypi: verify failed with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("pypi: failed to read verify response: %w", err)
	}

	var pkgInfo struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	if err := json.Unmarshal(body, &pkgInfo); err != nil {
		return fmt.Errorf("pypi: failed to parse verify response: %w", err)
	}

	if pkgInfo.Info.Version != version {
		return fmt.Errorf("pypi: published version %q does not match expected version %q",
			pkgInfo.Info.Version, version)
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

// addDistFile opens, copies, and closes a single dist file (BUG-pypi-defer-in-loop).
func (p *Publisher) addDistFile(w *multipart.Writer, f string) error {
	fh, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("pypi: failed to open %s: %w", f, err)
	}
	defer func() { _ = fh.Close() }()

	part, err := w.CreateFormFile("content", filepath.Base(f))
	if err != nil {
		return fmt.Errorf("pypi: failed to create form file %s: %w", f, err)
	}
	if _, err := io.Copy(part, fh); err != nil {
		return fmt.Errorf("pypi: failed to copy %s: %w", f, err)
	}
	return nil
}

// updateVersionInFile reads the file at the given path, finds the version
// assignment in the metadata section, replaces it, and writes the file back.
func updateVersionInFile(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("pypi: failed to read %s: %w", path, err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	found := false

	for i, line := range lines {
		if versionLinePattern.MatchString(line) {
			// Preserve the original key + assignment operator, replace the value.
			matches := versionLinePattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				prefix := matches[1]                  // "version = " or "version="
				quoted := strings.Contains(line, `"`) // TOML uses quotes
				comment := ""
				if len(matches) >= 4 && matches[3] != "" {
					comment = " " + matches[3]
				}

				if quoted {
					lines[i] = prefix + `"` + version + `"` + comment
				} else {
					lines[i] = prefix + version + comment
				}
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("pypi: no version field found in %s", path)
	}

	output := strings.Join(lines, "\n")
	if err := os.WriteFile(path, []byte(output), 0644); err != nil {
		return fmt.Errorf("pypi: failed to write %s: %w", path, err)
	}

	return nil
}

// readPackageName reads the package name from pyproject.toml or setup.cfg.
func readPackageName() (string, error) {
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return extractPackageName("pyproject.toml")
	}
	if _, err := os.Stat("setup.cfg"); err == nil {
		return extractPackageName("setup.cfg")
	}
	return "", fmt.Errorf("no supported config file found (pyproject.toml, setup.cfg)")
}

// extractPackageName extracts the package name from a config file.
// It looks for a name field in the metadata/project section.
func extractPackageName(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	content := string(data)
	inMetadataSection := false
	inProjectSection := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if trimmed == "[metadata]" {
			inMetadataSection = true
			inProjectSection = false
			continue
		}
		if trimmed == "[project]" {
			inProjectSection = true
			inMetadataSection = false
			continue
		}
		if strings.HasPrefix(trimmed, "[") {
			inMetadataSection = false
			inProjectSection = false
			continue
		}

		if inMetadataSection || inProjectSection {
			if strings.HasPrefix(trimmed, "name") && strings.Contains(trimmed, "=") {
				parts := strings.SplitN(trimmed, "=", 2)
				if len(parts) == 2 {
					name := strings.TrimSpace(parts[1])
					name = strings.Trim(name, `"'`)
					if name != "" {
						return name, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("package name not found in %s", path)
}
