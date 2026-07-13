// story: e03s03
package plugins

import (
	"fmt"
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

func TestExecPluginName(t *testing.T) {
	t.Run("SC-e03s03-P1-01: Name returns 'exec'", func(t *testing.T) {
		p := NewExecPlugin()
		if name := p.Name(); name != "exec" {
			t.Errorf("expected Name() == %q, got %q", "exec", name)
		}
	})
}

func TestExecPluginVerifyConditions(t *testing.T) {
	t.Run("SC-e03s03-P1-02: VerifyConditions passes with valid config", func(t *testing.T) {
		p := NewExecPlugin()
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Enabled: true, Options: map[string]string{"commands": "echo hello"}},
				},
			},
		}
		if err := p.VerifyConditions(ctx); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
	t.Run("SC-e03s03-P1-03: VerifyConditions fails with nil config", func(t *testing.T) {
		if err := NewExecPlugin().VerifyConditions(&algorithm.Context{}); err == nil {
			t.Error("expected error with nil config, got nil")
		}
	})
	t.Run("SC-e03s03-P1-04: VerifyConditions fails when exec publisher not configured", func(t *testing.T) {
		ctx := &algorithm.Context{Config: &algorithm.Config{Publishers: map[string]algorithm.PublisherConfig{}}}
		if err := NewExecPlugin().VerifyConditions(ctx); err == nil {
			t.Error("expected error when exec not configured, got nil")
		}
	})
	t.Run("SC-e03s03-P1-05: VerifyConditions fails when exec publisher not enabled", func(t *testing.T) {
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Enabled: false},
				},
			},
		}
		if err := NewExecPlugin().VerifyConditions(ctx); err == nil {
			t.Error("expected error when exec not enabled, got nil")
		}
	})
	t.Run("SC-e03s03-P1-06: VerifyConditions fails with no commands", func(t *testing.T) {
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Enabled: true, Options: map[string]string{}},
				},
			},
		}
		if err := NewExecPlugin().VerifyConditions(ctx); err == nil {
			t.Error("expected error with no commands, got nil")
		}
	})
}

func TestExecPluginNoop(t *testing.T) {
	p := NewExecPlugin()
	t.Run("SC-e03s03-P1-07: AnalyzeCommits returns empty", func(t *testing.T) {
		rt, err := p.AnalyzeCommits(&algorithm.Context{})
		if err != nil || rt != "" {
			t.Errorf("expected empty, got err=%v rt=%q", err, rt)
		}
	})
	t.Run("SC-e03s03-P1-08: GenerateNotes returns empty", func(t *testing.T) {
		notes, err := p.GenerateNotes(&algorithm.Context{})
		if err != nil || notes != "" {
			t.Errorf("expected empty, got err=%v notes=%q", err, notes)
		}
	})
}

func TestExecPluginPrepare(t *testing.T) {
	t.Run("SC-e03s03-P1-09: Prepare does nothing in dry-run mode", func(t *testing.T) {
		if err := NewExecPlugin().Prepare(&algorithm.Context{DryRun: true}); err != nil {
			t.Errorf("expected no error in dry-run, got: %v", err)
		}
	})
	t.Run("SC-e03s03-P1-10: Prepare runs configured commands", func(t *testing.T) {
		var executedCommands []string
		p := &ExecPlugin{runner: &mockCommandRunner{
			runFunc: func(name string, _ ...string) (string, error) {
				executedCommands = append(executedCommands, name)
				return "output", nil
			},
		}}
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Options: map[string]string{"commands": "echo hello\nmake build"}},
				},
			},
		}
		if err := p.Prepare(ctx); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if len(executedCommands) != 2 {
			t.Errorf("expected 2 commands, got %d: %v", len(executedCommands), executedCommands)
		}
	})
	t.Run("SC-e03s03-P1-11: Prepare skips comments and empty lines", func(t *testing.T) {
		var executedCommands []string
		p := &ExecPlugin{runner: &mockCommandRunner{
			runFunc: func(name string, _ ...string) (string, error) {
				executedCommands = append(executedCommands, name)
				return "output", nil
			},
		}}
		ctx := &algorithm.Config{
			Publishers: map[string]algorithm.PublisherConfig{
				"exec": {Options: map[string]string{"commands": "# comment\n\necho hello\n\nmake build"}},
			},
		}
		if err := p.Prepare(&algorithm.Context{Config: ctx}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if len(executedCommands) != 2 {
			t.Errorf("expected 2 commands, got %d", len(executedCommands))
		}
	})
	t.Run("SC-e03s03-P1-12: Prepare returns error on command failure", func(t *testing.T) {
		p := &ExecPlugin{runner: &mockCommandRunner{
			runFunc: func(_ string, _ ...string) (string, error) {
				return "", fmt.Errorf("command failed")
			},
		}}
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Options: map[string]string{"commands": "false"}},
				},
			},
		}
		if err := p.Prepare(ctx); err == nil {
			t.Error("expected error on command failure, got nil")
		}
	})
	t.Run("SC-e03s03-P1-13: Prepare passes arguments to commands", func(t *testing.T) {
		var capturedName string
		var capturedArgs []string
		p := &ExecPlugin{runner: &mockCommandRunner{
			runFunc: func(name string, args ...string) (string, error) {
				capturedName = name
				capturedArgs = args
				return "output", nil
			},
		}}
		ctx := &algorithm.Context{
			Config: &algorithm.Config{
				Publishers: map[string]algorithm.PublisherConfig{
					"exec": {Options: map[string]string{"commands": "echo hello world"}},
				},
			},
		}
		if err := p.Prepare(ctx); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if capturedName != "echo" {
			t.Errorf("expected 'echo', got %q", capturedName)
		}
		if len(capturedArgs) != 2 || capturedArgs[0] != "hello" || capturedArgs[1] != "world" {
			t.Errorf("expected [hello world], got %v", capturedArgs)
		}
	})
}

func TestExecPluginLifecycle(t *testing.T) {
	p := NewExecPlugin()
	t.Run("SC-e03s03-P1-14: Publish returns nil", func(t *testing.T) {
		release, err := p.Publish(&algorithm.Context{})
		if err != nil || release != nil {
			t.Errorf("expected nil, got err=%v release=%v", err, release)
		}
	})
	t.Run("SC-e03s03-P1-15: Success returns nil", func(t *testing.T) {
		if err := p.Success(&algorithm.Context{}); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
	t.Run("SC-e03s03-P1-16: Fail returns nil", func(t *testing.T) {
		if err := p.Fail(&algorithm.Context{}, nil); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}

func TestExecPluginAutoRegistration(t *testing.T) {
	t.Run("SC-e03s03-P1-17: ExecPlugin auto-registered in global registry", func(t *testing.T) {
		found := false
		for _, name := range List() {
			if name == "exec" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected 'exec' to be registered in global registry")
		}
	})
}
