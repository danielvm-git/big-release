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
	Branches       []BranchConfig             `yaml:"branches" json:"branches"`
	TagFormat      string                     `yaml:"tagFormat" json:"tagFormat"`
	Plugins        []string                   `yaml:"plugins" json:"plugins"`
	Publishers     map[string]PublisherConfig `yaml:"publishers" json:"publishers"`
	InitialVersion string                     `yaml:"initialVersion,omitempty" json:"initialVersion,omitempty"`
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

// Context represents the release context
type Context struct {
	Config        *Config
	Branch        *Branch
	LastRelease   *Release
	NextRelease   *Release
	Commits       []*Commit
	Releases      []*Release
	RepositoryURL string
	DryRun        bool
}
