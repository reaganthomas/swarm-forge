package runner_test

import (
	"errors"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/runner"
)

func TestRegisterAndMatch(t *testing.T) {
	reg := runner.NewRegistry()
	called := false
	err := reg.Register(`^a user exists$`, func(_ []string) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}

	fn, matches, err := reg.Match("a user exists")
	if err != nil {
		t.Fatalf("Match error: %v", err)
	}
	if len(matches) != 0 {
		t.Errorf("expected 0 captures, got %d", len(matches))
	}
	if err := fn(matches); err != nil {
		t.Fatalf("StepFunc error: %v", err)
	}
	if !called {
		t.Error("StepFunc was not called")
	}
}

func TestMatchWithCaptureGroups(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^the user "([^"]*)" exists$`, func(m []string) error {
		return nil
	})

	_, matches, err := reg.Match(`the user "Alice" exists`)
	if err != nil {
		t.Fatalf("Match error: %v", err)
	}
	if len(matches) != 1 || matches[0] != "Alice" {
		t.Errorf("captures = %v, want [Alice]", matches)
	}
}

func TestMatchNoMatch(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^a user exists$`, func(_ []string) error {
		return nil
	})

	_, _, err := reg.Match("something else")
	if err == nil {
		t.Fatal("expected error for no match")
	}
}

func TestMatchAmbiguous(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^a user`, func(_ []string) error { return nil })
	_ = reg.Register(`^a user exists`, func(_ []string) error { return nil })

	_, _, err := reg.Match("a user exists")
	if err == nil {
		t.Fatal("expected error for ambiguous match")
	}
}

func TestRegisterInvalidRegex(t *testing.T) {
	reg := runner.NewRegistry()
	err := reg.Register(`[invalid`, func(_ []string) error { return nil })
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestStepFuncReturnsError(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^it fails$`, func(_ []string) error {
		return errors.New("assertion failed")
	})

	fn, matches, _ := reg.Match("it fails")
	err := fn(matches)
	if err == nil || err.Error() != "assertion failed" {
		t.Errorf("expected assertion error, got: %v", err)
	}
}
