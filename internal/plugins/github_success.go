package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// Success is called after a successful release. It posts a comment to
// each issue/PR referenced in the release commits (via fixes/closes/
// resolves #N), so contributors and users get notified (#12). Commenting
// failures (403/404) are logged and non-fatal — they must not fail an
// otherwise-successful release.
func (p *GitHubPlugin) Success(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if ctx.DryRun {
		return nil
	}
	if state.NextRelease == nil || len(ctx.Commits) == 0 {
		return nil
	}

	// Collect referenced issue numbers from all release commit messages.
	seen := map[int]bool{}
	for _, c := range ctx.Commits {
		for _, n := range parseReferencedIssues(c.Message) {
			seen[n] = true
		}
	}
	if len(seen) == 0 {
		return nil
	}

	comment, err := p.buildSuccessComment(state.NextRelease.Version)
	if err != nil {
		// Invalid comment template is non-fatal in the success hook.
		fmt.Fprintf(os.Stderr, "warning: invalid successComment template: %v\n", err)
		return nil
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	for n := range seen {
		if err := p.postIssueComment(repo, n, comment); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not comment on issue #%d: %v\n", n, err)
		}
	}
	return nil
}

// buildSuccessComment renders the success comment template, or falls back
// to a sensible default when no template is configured.
func (p *GitHubPlugin) buildSuccessComment(version string) (string, error) {
	tmpl := p.successComment
	fallback := fmt.Sprintf("🎉 Released in version %s", version)
	return p.renderTemplate("success comment", tmpl, fallback, templateContext{Version: version})
}

// postIssueComment posts a comment to /repos/{repo}/issues/{n}/comments.
// 403/404 responses are returned as errors but are non-fatal at the
// caller (Success logs and continues).
func (p *GitHubPlugin) postIssueComment(repo string, issueNum int, body string) error {
	url := fmt.Sprintf("%s/repos/%s/issues/%d/comments", p.releaseURLBase(), repo, issueNum)
	payload, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create comment request: %w", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	// 403/404 are non-fatal — caller (Success) logs and continues.
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("comment on issue #%d failed (HTTP %d): %s", issueNum, resp.StatusCode, string(respBody))
}

// parseReferencedIssues returns the deduplicated, ordered list of issue
// numbers referenced in a commit message via #N, fixes #N, closes #N,
// or resolves #N. Uses the shared algorithm.IssueRefPattern.
func parseReferencedIssues(message string) []int {
	matches := algorithm.IssueRefPattern.FindAllStringSubmatch(message, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := map[int]bool{}
	var out []int
	for _, m := range matches {
		n, err := strconv.Atoi(m[2])
		if err != nil || n <= 0 || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	return out
}
