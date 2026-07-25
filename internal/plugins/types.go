// story: e21s03 e21s04
package plugins

// GitConfig configures the git plugin (release commit message and assets).
// Loaded from PluginConfigs["git"].
type GitConfig struct {
	Message       string   `yaml:"message,omitempty" json:"message,omitempty"`
	Assets        []string `yaml:"assets,omitempty" json:"assets,omitempty"`
	PostTagAssets []string `yaml:"postTagAssets,omitempty" json:"postTagAssets,omitempty"`
}
