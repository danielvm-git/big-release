package plugins

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// VersionFilePlugin writes a version file after the tag is created.
// Implements PostPublisher to run after Phase 6 (Publish).
type VersionFilePlugin struct {
	path     string
	template string
}

// NewVersionFilePlugin creates a new VersionFilePlugin with defaults.
func NewVersionFilePlugin() *VersionFilePlugin {
	return &VersionFilePlugin{
		path:     "VERSION",
		template: "{{.Version}}",
	}
}

// Name returns the plugin name.
func (p *VersionFilePlugin) Name() string {
	return "versionfile"
}

// Configure sets the plugin configuration.
func (p *VersionFilePlugin) Configure(raw map[string]interface{}) error {
	if path, ok := raw["path"].(string); ok && path != "" {
		p.path = path
	}
	if tmpl, ok := raw["template"].(string); ok && tmpl != "" {
		p.template = tmpl
	}
	return nil
}

// PostPublish writes the version file after the tag is created.
func (p *VersionFilePlugin) PostPublish(ctx *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	if state.NextRelease == nil {
		return nil
	}

	tmpl, err := template.New("version").Parse(p.template)
	if err != nil {
		return fmt.Errorf("versionfile: invalid template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, state.NextRelease); err != nil {
		return fmt.Errorf("versionfile: template render failed: %w", err)
	}

	dir := filepath.Dir(p.path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("versionfile: failed to create directory %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(p.path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("versionfile: failed to write %s: %w", p.path, err)
	}

	return nil
}

func init() {
	Register(NewVersionFilePlugin())
}
