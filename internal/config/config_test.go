package config

import (
	"os"
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
