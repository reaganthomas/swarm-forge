package gherkin_test

import (
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/gherkin"
)

func TestStepTypeString(t *testing.T) {
	tests := []struct {
		stepType gherkin.StepType
		want     string
	}{
		{gherkin.Given, "Given"},
		{gherkin.When, "When"},
		{gherkin.Then, "Then"},
	}

	for _, tt := range tests {
		if got := tt.stepType.String(); got != tt.want {
			t.Errorf("StepType.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestStepConstruction(t *testing.T) {
	step := gherkin.Step{Type: gherkin.Given, Text: "a user exists"}

	if step.Type != gherkin.Given {
		t.Errorf("step.Type = %v, want Given", step.Type)
	}
	if step.Text != "a user exists" {
		t.Errorf("step.Text = %q, want %q", step.Text, "a user exists")
	}
}

func TestScenarioConstruction(t *testing.T) {
	scenario := gherkin.Scenario{
		Name: "User logs in",
		Steps: []gherkin.Step{
			{Type: gherkin.Given, Text: "a user exists"},
			{Type: gherkin.When, Text: "they log in"},
			{Type: gherkin.Then, Text: "they see the dashboard"},
		},
	}

	if scenario.Name != "User logs in" {
		t.Errorf("scenario.Name = %q, want %q", scenario.Name, "User logs in")
	}
	if len(scenario.Steps) != 3 {
		t.Fatalf("len(scenario.Steps) = %d, want 3", len(scenario.Steps))
	}
}

func TestFeatureConstruction(t *testing.T) {
	feature := gherkin.Feature{
		Name: "Authentication",
		Scenarios: []gherkin.Scenario{
			{Name: "Login", Steps: []gherkin.Step{}},
		},
	}

	if feature.Name != "Authentication" {
		t.Errorf("feature.Name = %q, want %q", feature.Name, "Authentication")
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("len(feature.Scenarios) = %d, want 1", len(feature.Scenarios))
	}
}
