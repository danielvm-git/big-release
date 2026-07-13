package config

import "testing"

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
