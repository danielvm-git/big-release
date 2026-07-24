// story: e08s05
package release

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/danielvm-git/big-release/internal/algorithm"
	"github.com/danielvm-git/big-release/internal/git"
	"github.com/danielvm-git/big-release/internal/plugins"
)

type failingVerifyPlugin struct {
	name string
	msg  string
}

func (p *failingVerifyPlugin) Name() string { return p.name }
func (p *failingVerifyPlugin) VerifyConditions(_ *algorithm.ReadOnlyContext, _ *algorithm.MutableState) error {
	return fmt.Errorf("%s", p.msg)
}

var _ plugins.Plugin = (*failingVerifyPlugin)(nil)
var _ plugins.ConditionVerifier = (*failingVerifyPlugin)(nil)

func TestAggregateErrors_TwoFailures(t *testing.T) {
	plugins.Register(&failingVerifyPlugin{name: "fail-a", msg: "condition A failed"})
	plugins.Register(&failingVerifyPlugin{name: "fail-b", msg: "condition B failed"})

	gitClient, err := git.NewClient()
	if err != nil {
		t.Skipf("skipping: %v", err)
	}

	ctx := &Context{
		Config: &algorithm.Config{
			Branches:    []algorithm.BranchConfig{{Name: "main"}},
			TagFormat:   "v${version}",
			Plugins:     []string{"fail-a", "fail-b"},
			CommitTypes: algorithm.DefaultCommitTypes(),
		},
		Git:    gitClient,
		Logger: zap.NewNop(),
		DryRun: true,
	}

	r := New(ctx)
	algoCtx := &algorithm.ReadOnlyContext{
		Config: ctx.Config,
		Branch: &algorithm.Branch{Name: "main"},
		Commits: []*algorithm.Commit{
			{Type: "feat", Subject: "trigger release", Message: "feat: trigger release"},
		},
		DryRun: true,
	}
	state := &algorithm.MutableState{}

	err = r.runPluginLifecycle(algoCtx, state)
	if err == nil {
		t.Fatal("expected aggregate error")
	}
	if !strings.Contains(err.Error(), "condition A failed") {
		t.Errorf("missing error A: %v", err)
	}
	if !strings.Contains(err.Error(), "condition B failed") {
		t.Errorf("missing error B: %v", err)
	}

	var agg *AggregateError
	if !errors.As(err, &agg) {
		t.Fatalf("expected AggregateError, got %T: %v", err, err)
	}
	if len(agg.Unwrap()) < 2 {
		t.Fatalf("expected at least 2 wrapped errors, got %d", len(agg.Unwrap()))
	}
}

func TestAggregateErrors_Unwrap(t *testing.T) {
	e1 := fmt.Errorf("first failure")
	e2 := fmt.Errorf("second failure")
	agg := NewAggregateError(e1, e2)
	unwrapped := agg.Unwrap()
	if len(unwrapped) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(unwrapped))
	}
	if !errors.Is(errors.Join(unwrapped...), e1) {
		t.Error("expected first error in join")
	}
}
