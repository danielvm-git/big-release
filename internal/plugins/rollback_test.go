// story: NT-3 NT-7
package plugins

import (
	"testing"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// rollbackRecorder is a stub that proves Rollbacker is a capability a plugin
// can opt into without changing the base Plugin interface.
type rollbackRecorder struct {
	called bool
	state  *algorithm.MutableState
}

func (r *rollbackRecorder) Name() string { return "rollback-recorder" }

func (r *rollbackRecorder) Rollback(_ *algorithm.ReadOnlyContext, state *algorithm.MutableState) error {
	r.called = true
	r.state = state
	return nil
}

var _ Plugin = (*rollbackRecorder)(nil)
var _ Rollbacker = (*rollbackRecorder)(nil)

func TestRollbackerReceivesLifecycleState(t *testing.T) {
	rec := &rollbackRecorder{}
	state := &algorithm.MutableState{Status: algorithm.ReleaseStatusFailed}

	var rb Rollbacker = rec
	if err := rb.Rollback(&algorithm.ReadOnlyContext{}, state); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if !rec.called {
		t.Error("expected Rollback to be called")
	}
	if rec.state != state {
		t.Error("expected Rollback to receive the same MutableState pointer")
	}
}
