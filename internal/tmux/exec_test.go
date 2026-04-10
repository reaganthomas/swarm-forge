package tmux_test

import (
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/tmux"
)

func TestExecCommanderSatisfiesInterface(t *testing.T) {
	var cmd tmux.Commander = tmux.NewExecCommander()
	if cmd == nil {
		t.Fatal("NewExecCommander returned nil")
	}
}

func TestExecCommanderRunIsCallable(t *testing.T) {
	cmd := tmux.NewExecCommander()
	// list-sessions may fail if tmux isn't running — that's fine,
	// we just verify the method exists and is callable
	_ = cmd.Run("list-sessions")
}

func TestExecCommanderHasSessionReturnsFalse(t *testing.T) {
	cmd := tmux.NewExecCommander()
	// nonexistent session should return false
	if cmd.HasSession("nonexistent-session-xyzzy-12345") {
		t.Fatal("HasSession should return false for nonexistent session")
	}
}
