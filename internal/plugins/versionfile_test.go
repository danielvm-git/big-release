package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func TestVersionFilePlugin_Name(t *testing.T) {
	p := NewVersionFilePlugin()
	if p.Name() != "versionfile" {
		t.Errorf("expected Name() == 'versionfile', got %q", p.Name())
	}
}

func TestVersionFilePlugin_Configure(t *testing.T) {
	p := NewVersionFilePlugin()

	// Test with custom config
	raw := map[string]interface{}{
		"path":     "dist/VERSION",
		"template": "v{{.Version}}",
	}
	if err := p.Configure(raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.path != "dist/VERSION" {
		t.Errorf("expected path 'dist/VERSION', got %q", p.path)
	}
	if p.template != "v{{.Version}}" {
		t.Errorf("expected template 'v{{.Version}}', got %q", p.template)
	}
}

func TestVersionFilePlugin_Configure_Defaults(t *testing.T) {
	p := NewVersionFilePlugin()

	// Test with empty config (should keep defaults)
	raw := map[string]interface{}{}
	if err := p.Configure(raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.path != "VERSION" {
		t.Errorf("expected default path 'VERSION', got %q", p.path)
	}
	if p.template != "{{.Version}}" {
		t.Errorf("expected default template '{{.Version}}', got %q", p.template)
	}
}

func TestVersionFilePlugin_PostPublish(t *testing.T) {
	// Create temp dir for test output
	tmpDir := t.TempDir()
	versionFile := filepath.Join(tmpDir, "VERSION")

	p := NewVersionFilePlugin()
	p.path = versionFile

	ctx := &algorithm.ReadOnlyContext{}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{
			Version: "1.2.3",
			Type:    algorithm.ReleaseTypeMinor,
		},
	}

	if err := p.PostPublish(ctx, state); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("failed to read VERSION file: %v", err)
	}

	if string(content) != "1.2.3" {
		t.Errorf("expected VERSION file to contain '1.2.3', got %q", string(content))
	}
}

func TestVersionFilePlugin_PostPublish_CustomTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	versionFile := filepath.Join(tmpDir, "version.txt")

	p := NewVersionFilePlugin()
	p.path = versionFile
	p.template = "v{{.Version}} ({{.Type}})"

	ctx := &algorithm.ReadOnlyContext{}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{
			Version: "2.0.0",
			Type:    algorithm.ReleaseTypeMajor,
		},
	}

	if err := p.PostPublish(ctx, state); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}

	expected := "v2.0.0 (major)"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

func TestVersionFilePlugin_PostPublish_NestedDir(t *testing.T) {
	tmpDir := t.TempDir()
	versionFile := filepath.Join(tmpDir, "dist", "build", "VERSION")

	p := NewVersionFilePlugin()
	p.path = versionFile

	ctx := &algorithm.ReadOnlyContext{}
	state := &algorithm.MutableState{
		NextRelease: &algorithm.Release{
			Version: "0.1.0",
			Type:    algorithm.ReleaseTypePatch,
		},
	}

	if err := p.PostPublish(ctx, state); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("failed to read VERSION file: %v", err)
	}

	if string(content) != "0.1.0" {
		t.Errorf("expected '0.1.0', got %q", string(content))
	}
}

func TestVersionFilePlugin_PostPublish_NilRelease(t *testing.T) {
	p := NewVersionFilePlugin()

	ctx := &algorithm.ReadOnlyContext{}
	state := &algorithm.MutableState{
		NextRelease: nil,
	}

	// Should be a no-op, not an error
	if err := p.PostPublish(ctx, state); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersionFilePlugin_AutoRegistration(t *testing.T) {
	// Verify the plugin is auto-registered via init()
	plugin, err := Get("versionfile")
	if err != nil {
		t.Fatalf("versionfile plugin not registered: %v", err)
	}
	if plugin.Name() != "versionfile" {
		t.Errorf("expected plugin name 'versionfile', got %q", plugin.Name())
	}
}
