package runner_test

import (
	"errors"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/gherkin"
	"github.com/swarm-forge/swarm-forge/internal/runner"
)

func TestExecutorPassingScenario(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^step one$`, func(_ []string) error { return nil })
	_ = reg.Register(`^step two$`, func(_ []string) error { return nil })

	feature := gherkin.Feature{
		Name: "Test",
		Scenarios: []gherkin.Scenario{
			{
				Name: "All pass",
				Steps: []gherkin.Step{
					{Type: gherkin.Given, Text: "step one"},
					{Type: gherkin.Then, Text: "step two"},
				},
			},
		},
	}

	exec := runner.NewExecutor(reg)
	results := exec.Run(feature)

	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Status != runner.Passed {
		t.Errorf("scenario status = %v, want Passed", results[0].Status)
	}
	for i, sr := range results[0].Steps {
		if sr.Status != runner.Passed {
			t.Errorf("step %d status = %v, want Passed", i, sr.Status)
		}
	}
}

func TestExecutorFailingStep(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^it passes$`, func(_ []string) error { return nil })
	_ = reg.Register(`^it fails$`, func(_ []string) error {
		return errors.New("boom")
	})
	_ = reg.Register(`^after fail$`, func(_ []string) error { return nil })

	feature := gherkin.Feature{
		Name: "Test",
		Scenarios: []gherkin.Scenario{
			{
				Name: "Failure",
				Steps: []gherkin.Step{
					{Type: gherkin.Given, Text: "it passes"},
					{Type: gherkin.When, Text: "it fails"},
					{Type: gherkin.Then, Text: "after fail"},
				},
			},
		},
	}

	exec := runner.NewExecutor(reg)
	results := exec.Run(feature)

	sr := results[0]
	if sr.Status != runner.Failed {
		t.Errorf("scenario status = %v, want Failed", sr.Status)
	}
	if sr.Steps[0].Status != runner.Passed {
		t.Errorf("step 0 = %v, want Passed", sr.Steps[0].Status)
	}
	if sr.Steps[1].Status != runner.Failed {
		t.Errorf("step 1 = %v, want Failed", sr.Steps[1].Status)
	}
	if sr.Steps[1].Error != "boom" {
		t.Errorf("step 1 error = %q, want %q", sr.Steps[1].Error, "boom")
	}
	if sr.Steps[2].Status != runner.Skipped {
		t.Errorf("step 2 = %v, want Skipped", sr.Steps[2].Status)
	}
}

func TestExecutorPendingStep(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^defined step$`, func(_ []string) error { return nil })

	feature := gherkin.Feature{
		Name: "Test",
		Scenarios: []gherkin.Scenario{
			{
				Name: "Pending",
				Steps: []gherkin.Step{
					{Type: gherkin.Given, Text: "defined step"},
					{Type: gherkin.When, Text: "undefined step"},
					{Type: gherkin.Then, Text: "another step"},
				},
			},
		},
	}

	exec := runner.NewExecutor(reg)
	results := exec.Run(feature)

	sr := results[0]
	if sr.Status != runner.Pending {
		t.Errorf("scenario status = %v, want Pending", sr.Status)
	}
	if sr.Steps[1].Status != runner.Pending {
		t.Errorf("step 1 = %v, want Pending", sr.Steps[1].Status)
	}
	if sr.Steps[2].Status != runner.Skipped {
		t.Errorf("step 2 = %v, want Skipped", sr.Steps[2].Status)
	}
}

func TestExecutorMultipleScenarios(t *testing.T) {
	reg := runner.NewRegistry()
	_ = reg.Register(`^ok$`, func(_ []string) error { return nil })
	_ = reg.Register(`^bad$`, func(_ []string) error {
		return errors.New("fail")
	})

	feature := gherkin.Feature{
		Name: "Multi",
		Scenarios: []gherkin.Scenario{
			{Name: "Pass", Steps: []gherkin.Step{
				{Type: gherkin.Given, Text: "ok"},
			}},
			{Name: "Fail", Steps: []gherkin.Step{
				{Type: gherkin.Given, Text: "bad"},
			}},
		},
	}

	exec := runner.NewExecutor(reg)
	results := exec.Run(feature)

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Status != runner.Passed {
		t.Errorf("scenario 0 = %v, want Passed", results[0].Status)
	}
	if results[1].Status != runner.Failed {
		t.Errorf("scenario 1 = %v, want Failed", results[1].Status)
	}
}
