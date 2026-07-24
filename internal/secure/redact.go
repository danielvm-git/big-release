// story: e08s01
package secure

import (
	"os"
	"regexp"
	"strings"
)

// Redacted is the replacement string for sensitive values.
const Redacted = "[secure]"

var (
	sensitivePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(token|password|credential|secret|private)[=:]\s*\S+`),
		regexp.MustCompile(`(?i)(ghp_[a-zA-Z0-9]{20,})`),
		regexp.MustCompile(`(?i)(gho_[a-zA-Z0-9]{20,})`),
		regexp.MustCompile(`(?i)(npm_[a-zA-Z0-9]{20,})`),
	}

	knownSecretEnvVars = []string{
		"GH_TOKEN",
		"GITHUB_TOKEN",
		"NPM_TOKEN",
		"NODE_AUTH_TOKEN",
		"GITLAB_TOKEN",
		"CI_JOB_TOKEN",
	}
)

// Redact replaces sensitive patterns in text with Redacted.
func Redact(text string) string {
	for _, pattern := range sensitivePatterns {
		text = pattern.ReplaceAllString(text, Redacted)
	}
	return text
}

// RedactKnownSecrets replaces values of known secret environment variables in text.
func RedactKnownSecrets(text string) string {
	out := text
	for _, name := range knownSecretEnvVars {
		val := os.Getenv(name)
		if val == "" || len(val) < 4 {
			continue
		}
		out = strings.ReplaceAll(out, val, Redacted)
	}
	return Redact(out)
}
