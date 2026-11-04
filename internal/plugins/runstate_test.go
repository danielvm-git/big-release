package plugins

import "testing"

func TestGitHubPluginResetRunState_ClearsRememberedRelease(t *testing.T) {
	p := NewGitHubPlugin()
	p.recordRelease(4242)
	if got := p.rememberedReleaseID(); got != 4242 {
		t.Fatalf("rememberedReleaseID = %d, want 4242", got)
	}

	p.ResetRunState()

	if got := p.rememberedReleaseID(); got != 0 {
		t.Fatalf("after ResetRunState rememberedReleaseID = %d, want 0", got)
	}
}

func TestGitLabPluginResetRunState_ClearsRememberedRelease(t *testing.T) {
	p := NewGitLabPlugin()
	p.recordRelease("1.0.0")
	if got := p.rememberedReleaseTag(); got != "1.0.0" {
		t.Fatalf("rememberedReleaseTag = %q, want 1.0.0", got)
	}

	p.ResetRunState()

	if got := p.rememberedReleaseTag(); got != "" {
		t.Fatalf("after ResetRunState rememberedReleaseTag = %q, want empty", got)
	}
}
