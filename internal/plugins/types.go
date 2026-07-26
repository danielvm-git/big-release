// story: e21s03 e21s04
package plugins

// GitConfig configures the git plugin (release commit message and assets).
// Loaded from PluginConfigs["git"].
type GitConfig struct {
	Message       string   `yaml:"message,omitempty" json:"message,omitempty"`
	Assets        []string `yaml:"assets,omitempty" json:"assets,omitempty"`
	PostTagAssets []string `yaml:"postTagAssets,omitempty" json:"postTagAssets,omitempty"`

	// TagOnly publishes tags without pushing the release commit.
	//
	// For setups where the release identity cannot push to the default branch
	// (GitHub branch protection with no bypass available, GitLab protected
	// branches, any pre-receive policy). Opt in explicitly — the default keeps
	// semantic-release parity, where a failed commit push fails the release.
	TagOnly bool `yaml:"tagOnly,omitempty" json:"tagOnly,omitempty"`
}
