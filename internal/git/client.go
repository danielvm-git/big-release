package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// Client provides git operations
type Client struct {
	authorName  string
	authorEmail string
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

	cmd := exec.Command("git", "log", range_spec, "--pretty=format:%H|%s|%an|%ae|%ai|%b")
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
	cmd := exec.Command("git", "tag", "--list")
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
	cmd := exec.Command("git", "rev-list", "-1", tag)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get tag head: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetHead retrieves the current HEAD commit hash
func (c *Client) GetHead() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateTag creates an annotated git tag
func (c *Client) CreateTag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// Push pushes commits to the remote
func (c *Client) Push(remote string) error {
	cmd := exec.Command("git", "push", remote)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// AddNote adds a note to a ref
func (c *Client) AddNote(note, ref string) error {
	cmd := exec.Command("git", "notes", "--ref", "big-release-"+ref, "add", "-f", "-m", note, ref)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add note: %w", err)
	}

	return nil
}

// PushNotes pushes notes to remote
func (c *Client) PushNotes(remote, ref string) error {
	cmd := exec.Command("git", "push", remote, "refs/notes/big-release-"+ref)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push notes: %w", err)
	}

	return nil
}

// VerifyAuth verifies push authentication
func (c *Client) VerifyAuth(remote, branch string) error {
	cmd := exec.Command("git", "push", "--dry-run", "--no-verify", remote, "HEAD:"+branch)
	if err := cmd.Run(); err != nil {
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
	cmd := exec.Command("git", "ls-remote", "--heads", remote, branch)
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

// Commit creates a commit with the configured author
func (c *Client) Commit(message string) error {
	// Stage all changes
	stageCmd := exec.Command("git", "add", "-A")
	if err := stageCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create commit
	commitCmd := exec.Command("git", "commit", "-m", message, "--author="+c.authorName+" <"+c.authorEmail+">")
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// GetRepositoryURL retrieves the repository URL
func (c *Client) GetRepositoryURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repository URL: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch retrieves the current branch name
func (c *Client) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsGitRepo checks if the current directory is a git repository
func (c *Client) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetLastRelease retrieves the last release tag
func (c *Client) GetLastRelease(tagFormat string) (*algorithm.Release, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}

	// Parse tag format to extract version
	// Simple implementation: assume tags are like "v1.2.3"
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if strings.HasPrefix(tag, "v") {
			version := strings.TrimPrefix(tag, "v")
			head, err := c.GetTagHead(tag)
			if err != nil {
				continue
			}

			return &algorithm.Release{
				Version: version,
				GitTag:  tag,
				GitHead: head,
			}, nil
		}
	}

	return nil, nil
}

// GetCurrentTime returns the current time formatted for git
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05 -0700")
}

// StageChanges stages all changes in the working directory.
func (c *Client) StageChanges() error {
	cmd := exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}
	return nil
}

// HasChangesToCommit checks if there are any changes to commit.
func (c *Client) HasChangesToCommit() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check status: %w", err)
	}
	return strings.TrimSpace(string(output)) != "", nil
}

// PushTags pushes tags to the remote.
func (c *Client) PushTags(remote string) error {
	cmd := exec.Command("git", "push", remote, "--tags")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push tags: %w", err)
	}
	return nil
}

// DeleteTag deletes a local git tag.
func (c *Client) DeleteTag(tag string) error {
	cmd := exec.Command("git", "tag", "-d", tag)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
