// story: e08s02
package git

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/danielvm-git/big-release/internal/secure"
)

// AuthURL injects token into HTTPS remote URLs without mutating disk config.
// SSH URLs are returned unchanged. URLs that already contain credentials are
// not double-injected.
func AuthURL(remoteURL, token string) (string, error) {
	if token == "" {
		return remoteURL, nil
	}
	if remoteURL == "" {
		return "", fmt.Errorf("empty remote URL")
	}

	if strings.HasPrefix(remoteURL, "git@") || strings.HasPrefix(remoteURL, "ssh://") {
		return remoteURL, nil
	}

	parsed, err := url.Parse(remoteURL)
	if err != nil {
		return "", fmt.Errorf("invalid remote URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return remoteURL, nil
	}
	if parsed.User != nil {
		return remoteURL, nil
	}

	parsed.User = url.UserPassword("x-access-token", token)
	return parsed.String(), nil
}

// AuthURLError wraps an error with redacted message text.
func AuthURLError(err error, token string) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if token != "" {
		msg = strings.ReplaceAll(msg, token, secure.Redacted)
	}
	return fmt.Errorf("%s", secure.RedactKnownSecrets(msg))
}
