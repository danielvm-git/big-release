package packagist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/danielvm-git/big-release/internal/publishers"
)

const (
	DefaultAPIURL = "https://packagist.org"

	updateEndpoint = "/api/update-package"

	maxRetries = 3

	retryBase = 1 * time.Second

	envToken = "PACKAGIST_TOKEN"
)

type Publisher struct {
	APIURL     string
	HTTPClient *http.Client
	DryRun     bool
}

func NewPublisher() *Publisher {
	return &Publisher{
		APIURL:     DefaultAPIURL,
		HTTPClient: http.DefaultClient,
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

	body := map[string]interface{}{
		"repository": map[string]string{
			"url": p.APIURL,
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

		resp, doErr := p.HTTPClient.Do(req)
		if doErr != nil {
			lastErr = fmt.Errorf("packagist: request failed: %w", doErr)
			continue
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		switch {
		case resp.StatusCode == http.StatusOK:
			return nil
		case resp.StatusCode == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("packagist: rate limited (HTTP %d)", resp.StatusCode)
			continue
		case resp.StatusCode == http.StatusUnauthorized:
			return fmt.Errorf("packagist: authentication failed (HTTP %d): check PACKAGIST_TOKEN", resp.StatusCode)
		case resp.StatusCode >= 500:
			return fmt.Errorf("packagist: server error (HTTP %d)", resp.StatusCode)
		default:
			return fmt.Errorf("packagist: unexpected status (HTTP %d)", resp.StatusCode)
		}
	}

	return fmt.Errorf("packagist: publish failed after %d retries: %w", maxRetries, lastErr)
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

	resp, err := p.HTTPClient.Do(req)
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
