package algorithm

// Commit represents a parsed git commit
type Commit struct {
	Hash     string
	Message  string
	Type     string
	Scope    string
	Subject  string
	Breaking bool
	Author   string
	Email    string
	Date     string
	Body     string
	Footers  map[string]string
}

// ReleaseType represents the type of release
type ReleaseType string

const (
	ReleaseTypeMajor      ReleaseType = "major"
	ReleaseTypeMinor      ReleaseType = "minor"
	ReleaseTypePatch      ReleaseType = "patch"
	ReleaseTypePrerelease ReleaseType = "prerelease"
)

// Release represents a release to be created
type Release struct {
	Version string
	GitTag  string
	GitHead string
	Channel string
	Notes   string
	Type    ReleaseType
	Branch  string
	Assets  []Asset
}

// Asset represents a release asset
type Asset struct {
	Name string
	Path string
	Type string
}

// Branch represents a git branch configuration
type Branch struct {
	Name       string
	Type       BranchType
	Range      string
	Channel    string
	Prerelease string
	Tags       []Tag
}

// BranchType represents the type of branch
type BranchType string

const (
	BranchTypeRelease     BranchType = "release"
	BranchTypeMaintenance BranchType = "maintenance"
	BranchTypePrerelease  BranchType = "prerelease"
)

// Tag represents a git tag
type Tag struct {
	Version  string
	GitTag   string
	GitHead  string
	Channels []string
}

// Config represents the release configuration
type Config struct {
	Branches       []BranchConfig                    `yaml:"branches" json:"branches"`
	TagFormat      string                            `yaml:"tagFormat" json:"tagFormat"`
	Plugins        []string                          `yaml:"plugins" json:"plugins"`
	Publishers     map[string]PublisherConfig        `yaml:"publishers" json:"publishers"`
	CommitTypes    []CommitTypeConfig                `yaml:"commitTypes,omitempty" json:"commitTypes,omitempty"`
	PluginConfigs  map[string]map[string]interface{} `yaml:"pluginConfigs,omitempty" json:"pluginConfigs,omitempty"`
	InitialVersion string                            `yaml:"initialVersion,omitempty" json:"initialVersion,omitempty"`
	ChangelogTitle string                            `yaml:"changelogTitle,omitempty" json:"changelogTitle,omitempty"`
	VersionFile    *VersionFileConfig                `yaml:"versionFile,omitempty" json:"versionFile,omitempty"`
}

// VersionFileConfig configures the versionFile plugin.
type VersionFileConfig struct {
	Path     string `yaml:"path,omitempty" json:"path,omitempty"`
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
}

// CommitTypeConfig configures how a single conventional commit type is
// rendered (and whether it is rendered at all) in release notes. Hidden
// types are excluded from the changelog unless their commit carries a
// breaking change.
type CommitTypeConfig struct {
	Type    string `yaml:"type" json:"type"`
	Section string `yaml:"section,omitempty" json:"section,omitempty"` // display title; empty = hidden
	Hidden  bool   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
}

// BranchConfig represents branch configuration
type BranchConfig struct {
	Name       string `yaml:"name" json:"name"`
	Type       string `yaml:"type,omitempty" json:"type,omitempty"`
	Channel    string `yaml:"channel,omitempty" json:"channel,omitempty"`
	Prerelease string `yaml:"prerelease,omitempty" json:"prerelease,omitempty"`
}

// PublisherConfig represents publisher configuration
type PublisherConfig struct {
	Enabled bool              `yaml:"enabled" json:"enabled"`
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}

// AssetConfig declares a binary artifact to attach to a GitHub release.
// Path may be a glob pattern (e.g. "dist/*.tar.gz"); Label is an optional
// display name overriding the file name on the release page.
type AssetConfig struct {
	Path  string `yaml:"path" json:"path"`
	Label string `yaml:"label,omitempty" json:"label,omitempty"`
}

// GitHubConfig configures the github plugin (release creation, assets,
// commenting). Loaded from the PluginConfigs["github"] entry.
type GitHubConfig struct {
	Assets         []AssetConfig `yaml:"assets,omitempty" json:"assets,omitempty"`
	DraftRelease   bool          `yaml:"draftRelease,omitempty" json:"draftRelease,omitempty"`
	ReleaseName    string        `yaml:"releaseName,omitempty" json:"releaseName,omitempty"`
	ReleaseBody    string        `yaml:"releaseBody,omitempty" json:"releaseBody,omitempty"`
	SuccessComment string        `yaml:"successComment,omitempty" json:"successComment,omitempty"`
	ReleasedLabels []string      `yaml:"releasedLabels,omitempty" json:"releasedLabels,omitempty"`
}

// ReadOnlyContext holds immutable inputs to the release pipeline.
// Constructed once and never modified.
type ReadOnlyContext struct {
	Config        *Config
	Branch        *Branch
	Commits       []*Commit
	Releases      []*Release
	RepositoryURL string
	DryRun        bool
}

// MutableState holds values written by specific pipeline stages.
// Passed as a pointer so stages can update it.
type MutableState struct {
	LastRelease *Release
	NextRelease *Release
	Notes       string
	Assets      []Asset
}
