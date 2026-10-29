package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func TestDefaultConfig_InitialVersion(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.InitialVersion != "0.1.0" {
		t.Errorf("expected InitialVersion 0.1.0, got %q", cfg.InitialVersion)
	}
}

func TestDefaultConfig_Branches(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Branches) == 0 {
		t.Error("expected default branches, got none")
	}
}

func TestDefaultConfig_PassesValidate(t *testing.T) {
	if err := ValidateConfig(DefaultConfig()); err != nil {
		t.Fatalf("DefaultConfig should pass ValidateConfig: %v", err)
	}
}

func TestDefaultConfig_TagFormat(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TagFormat != "v${version}" {
		t.Errorf("expected tag format v${version}, got %q", cfg.TagFormat)
	}
}

func TestValidateConfig_RejectsInvalidBranchType(t *testing.T) {
	cfg := &algorithm.Config{
		Branches: []algorithm.BranchConfig{
			{Name: "main", Type: "invalid-type"},
		},
		TagFormat: "v${version}",
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid branch type, got nil")
	}
	if !strings.Contains(err.Error(), "invalid branch type") {
		t.Errorf("expected 'invalid branch type' error, got: %v", err)
	}
}

func TestValidateConfig_AcceptsValidBranchTypes(t *testing.T) {
	validTypes := []string{"", "release", "maintenance", "prerelease"}
	for _, bt := range validTypes {
		cfg := &algorithm.Config{
			Branches: []algorithm.BranchConfig{
				{Name: "main", Type: bt},
			},
			TagFormat: "v${version}",
		}
		err := ValidateConfig(cfg)
		if err != nil {
			t.Errorf("expected no error for branch type %q, got: %v", bt, err)
		}
	}
}

func TestValidateConfig_RejectsDuplicateBranchNames(t *testing.T) {
	cfg := &algorithm.Config{
		Branches: []algorithm.BranchConfig{
			{Name: "main", Type: "release"},
			{Name: "main", Type: "prerelease"},
		},
		TagFormat: "v${version}",
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate branch names, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate branch name") {
		t.Errorf("expected 'duplicate branch name' error, got: %v", err)
	}
}

// --- e18s01: configurable commit type sections & visibility ---

func TestDefaultConfig_SeedsCommitTypes(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.CommitTypes) == 0 {
		t.Fatal("expected DefaultConfig to seed CommitTypes, got empty slice")
	}
	seen := map[string]bool{}
	for _, ct := range cfg.CommitTypes {
		seen[ct.Type] = true
	}
	for _, want := range []string{"feat", "fix", "perf", "docs", "refactor", "chore", "style", "test"} {
		if !seen[want] {
			t.Errorf("DefaultConfig.CommitTypes missing %q: %+v", want, seen)
		}
	}
}

// --- e19s01 (#10): plugin config round-trip ---

func TestLoadPluginConfigs_AssetsRoundTrip(t *testing.T) {
	yml := `
branches:
  - name: main
tagFormat: "v${version}"
plugins:
  - github
pluginConfigs:
  github:
    assets:
      - path: dist/*.tar.gz
        label: "big-release (linux-amd64)"
      - path: bin/big-release
`
	tmp := t.TempDir() + "/.big-release.yml"
	if err := os.WriteFile(tmp, []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	gh, ok := cfg.PluginConfigs["github"]
	if !ok {
		t.Fatal("expected github pluginConfigs entry")
	}
	assets, _ := gh["assets"].([]interface{})
	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d: %+v", len(assets), gh["assets"])
	}
}

// --- e08s06: configuration file loading ---

func TestFileLoading_YAML(t *testing.T) {
	yml := `
branches:
  - name: main
tagFormat: "release-${version}"
`
	path := filepath.Join(t.TempDir(), ".big-release.yml")
	if err := os.WriteFile(path, []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TagFormat != "release-${version}" {
		t.Fatalf("expected custom tagFormat, got %q", cfg.TagFormat)
	}
}

func TestFileLoading_JSON(t *testing.T) {
	js := `{"branches":[{"name":"main"}],"tagFormat":"json-${version}"}`
	path := filepath.Join(t.TempDir(), ".big-release.json")
	if err := os.WriteFile(path, []byte(js), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TagFormat != "json-${version}" {
		t.Fatalf("expected json tagFormat, got %q", cfg.TagFormat)
	}
}

func TestFileLoading_ParentDiscovery(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "sub", "dir")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	yml := `
branches:
  - name: main
tagFormat: "parent-${version}"
`
	if err := os.WriteFile(filepath.Join(root, ".big-release.yml"), []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TagFormat != "parent-${version}" {
		t.Fatalf("expected parent config, got %q", cfg.TagFormat)
	}
}

func TestFileLoading_ExplicitPath(t *testing.T) {
	dir := t.TempDir()
	explicit := filepath.Join(dir, "custom.yml")
	other := filepath.Join(dir, ".big-release.yml")
	if err := os.WriteFile(explicit, []byte("branches:\n  - name: main\ntagFormat: explicit-${version}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(other, []byte("branches:\n  - name: main\ntagFormat: default-${version}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(explicit)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TagFormat != "explicit-${version}" {
		t.Fatalf("expected explicit path config, got %q", cfg.TagFormat)
	}
}

func TestLoad_ExplicitNonexistentFile_ReturnsError(t *testing.T) {
	// BUG-config-explicit-fallback: when the user explicitly provides a config
	// path that doesn't exist, Load should return an error — not silently fall
	// back to defaults. A typo in --config should not go undetected.
	_, err := Load("/tmp/nonexistent-big-release-config-12345.yml")
	if err == nil {
		t.Fatal("expected error for nonexistent explicit config file, got nil")
	}
}

// --- BUG-branch-string-form: README documents `branches: - main` ---

func loadConfigText(t *testing.T, name, text string) *algorithm.Config {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%s) failed: %v", name, err)
	}
	return cfg
}

func TestFileLoading_YAMLBranchStringForm(t *testing.T) {
	cfg := loadConfigText(t, ".big-release.yml", `
branches:
  - main
  - "N.x"
  - name: beta
    prerelease: true
tagFormat: "v${version}"
`)
	want := []string{"main", "N.x", "beta"}
	if len(cfg.Branches) != len(want) {
		t.Fatalf("expected %d branches, got %d: %+v", len(want), len(cfg.Branches), cfg.Branches)
	}
	for i, name := range want {
		if cfg.Branches[i].Name != name {
			t.Errorf("branch %d: expected name %q, got %q", i, name, cfg.Branches[i].Name)
		}
	}
	if cfg.Branches[2].Prerelease != "true" {
		t.Errorf("object form must keep its fields, got prerelease %q", cfg.Branches[2].Prerelease)
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("string-form branches must validate: %v", err)
	}
}

func TestFileLoading_JSONBranchStringForm(t *testing.T) {
	cfg := loadConfigText(t, ".big-release.json",
		`{"branches":["main",{"name":"next","channel":"next"}],"tagFormat":"v${version}"}`)
	if len(cfg.Branches) != 2 || cfg.Branches[0].Name != "main" || cfg.Branches[1].Channel != "next" {
		t.Fatalf("unexpected branches: %+v", cfg.Branches)
	}
}

func TestFileLoading_YAMLBranchRejectsSequence(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".big-release.yml")
	if err := os.WriteFile(path, []byte("branches:\n  - [main]\ntagFormat: \"v${version}\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("a nested list is not a branch and must fail to parse")
	}
}
