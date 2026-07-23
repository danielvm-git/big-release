package git

import "github.com/danielvm-git/big-release/internal/algorithm"

// GitAPI defines the interface for git operations used by the orchestrator.
type GitAPI interface {
	GetCommits(from, to string) ([]*algorithm.Commit, error)
	GetTags() ([]string, error)
	GetTagHead(tag string) (string, error)
	GetHead() (string, error)
	GetLastRelease(tagFormat string) (*algorithm.Release, error)
	GetRepositoryURL() (string, error)
	GetCurrentBranch() (string, error)
	IsGitRepo() bool
	StageChanges() error
	GetModifiedFiles() ([]string, error)
	StagePaths(paths []string) error
	HasChangesToCommit() (bool, error)
	Commit(message string) error
	CreateTag(tag, message string) error
	Push(remote string) error
	PushTags(remote string) error
	DeleteTag(tag string) error
}
