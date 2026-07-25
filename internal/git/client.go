package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// Client provides git operations
type Client struct {
	authorName  string
	authorEmail string
}

// gitCmd builds a git command with the calling process's git worktree/hook
// context stripped from its environment (GIT_DIR, GIT_WORK_TREE,
// GIT_INDEX_FILE, etc.). Without this, a Client invoked from inside a git
// hook (e.g. big-release running from a pre-commit hook, or any test run
// while a commit is in flight) inherits those variables and silently
// operates on the hook's repository instead of the one at the command's Dir
// (see BUG-tag-ignores-tagformat investigation).
func gitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Env = scrubGitEnv()
	return cmd
}

// runGit runs cmd and, on failure, appends any stderr git produced to the
// error. cmd.Run() alone connects a nil Stderr to /dev/null, so callers
// otherwise only ever see the bare "exit status 1" — the actual reason
// (rejected push, auth failure, etc.) is silently discarded.
func runGit(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return fmt.Errorf("%w: %s", err, msg)
		}
		return err
	}
	return nil
}

func scrubGitEnv() []string {
	env := os.Environ()
	out := make([]string, 0, len(env))
	for _, e := range env {
		switch {
		case strings.HasPrefix(e, "GIT_DIR="),
			strings.HasPrefix(e, "GIT_WORK_TREE="),
			strings.HasPrefix(e, "GIT_INDEX_FILE="),
			strings.HasPrefix(e, "GIT_PREFIX="),
			strings.HasPrefix(e, "GIT_COMMON_DIR="),
			strings.HasPrefix(e, "GIT_OBJECT_DIRECTORY="),
			strings.HasPrefix(e, "GIT_ALTERNATE_OBJECT_DIRECTORIES="):
			continue
		default:
			out = append(out, e)
		}
	}
	return out
}

// NewClient creates a new git Client
func NewClient() (*Client, error) {
	return &Client{
		authorName:  "big-release[bot]",
		authorEmail: "big-release@users.noreply.github.com",
	}, nil
}

// GetCommits retrieves commits between two refs
func (c *Client) GetCommits(from, to string) ([]*algorithm.Commit, error) {
	range_spec := ""
	if from != "" {
		range_spec = from + ".." + to
	} else {
		range_spec = to
	}

	cmd := gitCmd("log", range_spec, "--pretty=format:%H|%s|%an|%ae|%ai|%b")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	var commits []*algorithm.Commit
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 6)
		if len(parts) < 6 {
			continue
		}

		commits = append(commits, &algorithm.Commit{
			Hash:    parts[0],
			Message: parts[1],
			Author:  parts[2],
			Email:   parts[3],
			Date:    parts[4],
			Body:    parts[5],
		})
	}

	return commits, nil
}

// GetTags retrieves all tags
func (c *Client) GetTags() ([]string, error) {
	cmd := gitCmd("tag", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	var tags []string
	for _, tag := range strings.Split(string(output), "\n") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// GetTagHead retrieves the commit hash for a tag
func (c *Client) GetTagHead(tag string) (string, error) {
	cmd := gitCmd("rev-list", "-1", tag)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tag head: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetHead retrieves the current HEAD commit hash
func (c *Client) GetHead() (string, error) {
	cmd := gitCmd("rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateTag creates an annotated git tag
func (c *Client) CreateTag(tag, message string) error {
	cmd := gitCmd("tag", "-a", tag, "-m", message)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// Push pushes commits to the remote
func (c *Client) Push(remote string) error {
	cmd := gitCmd("push", remote)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// AddNote adds a note to a ref
func (c *Client) AddNote(note, ref string) error {
	cmd := gitCmd("notes", "--ref", "big-release-"+ref, "add", "-f", "-m", note, ref)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to add note: %w", err)
	}

	return nil
}

// PushNotes pushes notes to remote
func (c *Client) PushNotes(remote, ref string) error {
	cmd := gitCmd("push", remote, "refs/notes/big-release-"+ref)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to push notes: %w", err)
	}

	return nil
}

// VerifyAuth verifies push authentication
func (c *Client) VerifyAuth(remote, branch string) error {
	cmd := gitCmd("push", "--dry-run", "--no-verify", remote, "HEAD:"+branch)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

// IsBranchUpToDate checks if branch is up to date with remote
func (c *Client) IsBranchUpToDate(remote, branch string) (bool, error) {
	// Get local HEAD
	localHead, err := c.GetHead()
	if err != nil {
		return false, err
	}

	// Get remote HEAD
	cmd := gitCmd("ls-remote", "--heads", remote, branch)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to get remote HEAD: %w", err)
	}

	// Parse remote HEAD
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			remoteHead := parts[0]
			if localHead == remoteHead {
				return true, nil
			}
		}
	}

	return false, nil
}

// Commit creates a commit from whatever the caller has already staged.
// It intentionally does not stage anything itself (previously an unconditional
// `git add -A`) — that silently re-broadened every commit to the entire
// working tree regardless of what a caller had deliberately staged, defeating
// selective staging via StagePaths/matchModifiedAssets (see
// BUG-git-push-error-swallowed). Callers that want everything staged still
// call StageChanges explicitly first.
func (c *Client) Commit(message string) error {
	commitCmd := gitCmd("commit", "-m", message, "--author="+c.authorName+" <"+c.authorEmail+">")
	if err := runGit(commitCmd); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// GetRepositoryURL retrieves the repository URL
func (c *Client) GetRepositoryURL() (string, error) {
	cmd := gitCmd("config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repository URL: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch retrieves the current branch name
func (c *Client) GetCurrentBranch() (string, error) {
	cmd := gitCmd("rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsGitRepo checks if the current directory is a git repository
func (c *Client) IsGitRepo() bool {
	cmd := gitCmd("rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetLastRelease retrieves the last release tag
func (c *Client) GetLastRelease(tagFormat string) (*algorithm.Release, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}

	prefix := tagPrefix(tagFormat)

	type candidate struct {
		tag     string
		version *semver.Version
	}
	var matches []candidate

	for _, tag := range tags {
		if !strings.HasPrefix(tag, prefix) {
			continue
		}
		versionStr := strings.TrimPrefix(tag, prefix)
		if versionStr == "" || !strings.Contains(versionStr, ".") {
			continue
		}
		v, err := semver.NewVersion(versionStr)
		if err != nil {
			continue
		}
		matches = append(matches, candidate{tag: tag, version: v})
	}

	if len(matches) == 0 {
		return nil, nil
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].version.GreaterThan(matches[j].version)
	})

	best := matches[0]
	head, err := c.GetTagHead(best.tag)
	if err != nil {
		return nil, err
	}

	return &algorithm.Release{
		Version: best.version.String(),
		GitTag:  best.tag,
		GitHead: head,
	}, nil
}

// FormatTag applies tagFormat to a semver version string, e.g.
// FormatTag("v${version}", "1.0.0") -> "v1.0.0". Empty tagFormat falls back
// to the bare version, matching the pre-existing default behavior.
// Shared by GetLastRelease (read path) and plugins.GitPlugin (write path) so
// the two can't drift out of sync again (see BUG-tag-ignores-tagformat).
func FormatTag(tagFormat, version string) string {
	return tagPrefix(tagFormat) + version
}

func tagPrefix(tagFormat string) string {
	if idx := strings.Index(tagFormat, "${version}"); idx > 0 {
		return tagFormat[:idx]
	}
	return ""
}

// GetCurrentTime returns the current time formatted for git
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05 -0700")
}

// StageChanges stages all changes in the working directory.
func (c *Client) StageChanges() error {
	cmd := gitCmd("add", ".")
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}
	return nil
}

// GetModifiedFiles returns paths with unstaged or staged modifications.
func (c *Client) GetModifiedFiles() ([]string, error) {
	cmd := gitCmd("status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 4 {
			continue
		}
		path := strings.TrimSpace(line[3:])
		if idx := strings.Index(path, " -> "); idx >= 0 {
			path = path[idx+4:]
		}
		if path != "" {
			files = append(files, path)
		}
	}
	return files, nil
}

// StagePaths stages the given paths for commit.
func (c *Client) StagePaths(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	args := append([]string{"add", "--"}, paths...)
	cmd := gitCmd(args...)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to stage paths: %w", err)
	}
	return nil
}

// HasChangesToCommit checks whether anything is staged for commit. This
// deliberately looks at the index (git diff --cached), not overall working
// tree dirtiness (git status --porcelain would also be true from files an
// earlier plugin wrote but nobody chose to stage) — otherwise a caller that
// intentionally staged nothing would still trigger a commit attempt (see
// BUG-git-push-error-swallowed).
func (c *Client) HasChangesToCommit() (bool, error) {
	cmd := gitCmd("diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check status: %w", err)
	}
	return strings.TrimSpace(string(output)) != "", nil
}

// PushTags pushes tags to the remote.
func (c *Client) PushTags(remote string) error {
	cmd := gitCmd("push", remote, "--tags")
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to push tags: %w", err)
	}
	return nil
}

// DeleteTag deletes a local git tag.
func (c *Client) DeleteTag(tag string) error {
	cmd := gitCmd("tag", "-d", tag)
	if err := runGit(cmd); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
