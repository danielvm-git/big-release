package plugins

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// uploadAssets expands configured asset globs and uploads each file to
// the given release via the GitHub uploads endpoint. Missing files are
// logged as warnings and skipped (non-fatal).
func (p *GitHubPlugin) uploadAssets(repo string, releaseID int64) error {
	assets, errs := expandAssetGlobs(p.assets)
	// Missing globs/files are warnings, not failures (#10 acceptance).
	for _, errMsg := range errs {
		fmt.Fprintf(os.Stderr, "warning: %s\n", errMsg)
	}

	for _, asset := range assets {
		if err := p.uploadOneAsset(repo, releaseID, asset); err != nil {
			return err
		}
	}
	return nil
}

// uploadOneAsset uploads a single file to the release.
func (p *GitHubPlugin) uploadOneAsset(repo string, releaseID int64, asset algorithm.AssetConfig) error {
	f, err := os.Open(asset.Path)
	if err != nil {
		// Missing file is a warning, not a failure.
		fmt.Fprintf(os.Stderr, "warning: could not open asset %q: %v\n", asset.Path, err)
		return nil
	}
	defer func() { _ = f.Close() }()

	name := asset.Label
	if name == "" {
		name = filepath.Base(asset.Path)
	}
	uploadURL := fmt.Sprintf("%s/repos/%s/releases/%d/assets?name=%s",
		p.uploadsHost(), repo, releaseID, url.QueryEscape(name))

	req, err := http.NewRequest(http.MethodPost, uploadURL, f)
	if err != nil {
		return fmt.Errorf("failed to create asset upload request for %q: %w", name, err)
	}
	if stat, err := f.Stat(); err == nil {
		req.ContentLength = stat.Size()
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", mimeTypeForAsset(name))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload asset %q: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("asset upload %q failed (HTTP %d): %s", name, resp.StatusCode, string(body))
	}
	return nil
}

// uploadsHost returns the GitHub uploads base URL.
func (p *GitHubPlugin) uploadsHost() string {
	if p.uploadBaseURL != "" {
		return p.uploadBaseURL
	}
	return "https://uploads.github.com"
}

// expandAssetGlobs expands glob patterns in asset paths to concrete files.
// Returns the expanded list plus a slice of error messages for any pattern
// that matched no files (non-fatal — caller may log and continue).
func expandAssetGlobs(assets []algorithm.AssetConfig) ([]algorithm.AssetConfig, []string) {
	var out []algorithm.AssetConfig
	var errs []string

	for _, asset := range assets {
		matches, err := filepath.Glob(asset.Path)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid glob %q: %v", asset.Path, err))
			continue
		}
		if len(matches) == 0 {
			// Could be a literal path that doesn't exist, or an empty glob.
			// Keep it as-is so the uploader can emit a missing-file warning.
			out = append(out, asset)
			continue
		}
		for _, m := range matches {
			out = append(out, algorithm.AssetConfig{Path: m, Label: asset.Label})
		}
	}
	return out, errs
}

// mimeTypeForAsset returns the MIME type for an asset based on its
// filename extension. Falls back to application/octet-stream.
func mimeTypeForAsset(name string) string {
	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".gz", ".tgz":
		return "application/gzip"
	case ".zip":
		return "application/zip"
	case ".exe":
		return "application/vnd.microsoft.portable-executable"
	case ".dmg":
		return "application/x-apple-diskimage"
	case ".deb":
		return "application/vnd.debian.binary-package"
	case ".rpm":
		return "application/x-rpm"
	default:
		if t := mime.TypeByExtension(ext); t != "" {
			return t
		}
		return "application/octet-stream"
	}
}
