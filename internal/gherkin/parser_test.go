package gherkin_test

import (
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/gherkin"
)

func TestParseFeatureName(t *testing.T) {
	input := "Feature: User Authentication\n"
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if feature.Name != "User Authentication" {
		t.Errorf("feature.Name = %q, want %q", feature.Name, "User Authentication")
	}
}

func TestParseSingleScenario(t *testing.T) {
	input := `Feature: Login

  Scenario: Successful login
    Given a registered user
    When they enter valid credentials
    Then they see the dashboard
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("got %d scenarios, want 1", len(feature.Scenarios))
	}
	sc := feature.Scenarios[0]
	if sc.Name != "Successful login" {
		t.Errorf("scenario.Name = %q, want %q", sc.Name, "Successful login")
	}
	if len(sc.Steps) != 3 {
		t.Fatalf("got %d steps, want 3", len(sc.Steps))
	}
	assertStep(t, sc.Steps[0], gherkin.Given, "a registered user")
	assertStep(t, sc.Steps[1], gherkin.When, "they enter valid credentials")
	assertStep(t, sc.Steps[2], gherkin.Then, "they see the dashboard")
}

func TestParseAndInheritsPreviousType(t *testing.T) {
	input := `Feature: Multi-step
  Scenario: And keyword
    Given a user exists
    And they are logged in
    When they click logout
    Then they see the login page
    And they cannot access the dashboard
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	steps := feature.Scenarios[0].Steps
	if len(steps) != 5 {
		t.Fatalf("got %d steps, want 5", len(steps))
	}
	assertStep(t, steps[0], gherkin.Given, "a user exists")
	assertStep(t, steps[1], gherkin.Given, "they are logged in")
	assertStep(t, steps[2], gherkin.When, "they click logout")
	assertStep(t, steps[3], gherkin.Then, "they see the login page")
	assertStep(t, steps[4], gherkin.Then, "they cannot access the dashboard")
}

func TestParseButInheritsPreviousType(t *testing.T) {
	input := `Feature: But keyword
  Scenario: But step
    Given a user exists
    But they are not admin
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	steps := feature.Scenarios[0].Steps
	if len(steps) != 2 {
		t.Fatalf("got %d steps, want 2", len(steps))
	}
	assertStep(t, steps[1], gherkin.Given, "they are not admin")
}

func TestParseMultipleScenarios(t *testing.T) {
	input := `Feature: Multiple

  Scenario: First
    Given step one

  Scenario: Second
    Given step two
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(feature.Scenarios) != 2 {
		t.Fatalf("got %d scenarios, want 2", len(feature.Scenarios))
	}
	if feature.Scenarios[0].Name != "First" {
		t.Errorf("first scenario name = %q", feature.Scenarios[0].Name)
	}
	if feature.Scenarios[1].Name != "Second" {
		t.Errorf("second scenario name = %q", feature.Scenarios[1].Name)
	}
}

func TestParseIgnoresCommentsAndBlankLines(t *testing.T) {
	input := `# This is a comment
Feature: Comments

  # Another comment
  Scenario: With comments
    Given something

`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if feature.Name != "Comments" {
		t.Errorf("feature.Name = %q, want %q", feature.Name, "Comments")
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("got %d scenarios, want 1", len(feature.Scenarios))
	}
}

func TestParseIgnoresDescriptionLines(t *testing.T) {
	input := `Feature: With description

  This is a description paragraph that should be ignored.
  It can span multiple lines.

  Scenario: After description
    Given something
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("got %d scenarios, want 1", len(feature.Scenarios))
	}
	if feature.Scenarios[0].Name != "After description" {
		t.Errorf("name = %q", feature.Scenarios[0].Name)
	}
}

func TestParseSkipsDocStringBlocks(t *testing.T) {
	input := "Feature: Docstrings\n" +
		"  Scenario: With docstring\n" +
		"    Given a precondition\n" +
		"      \"\"\"\n" +
		"      Feature: Nested\n" +
		"        Scenario: Inner\n" +
		"          Given inside docstring\n" +
		"      \"\"\"\n" +
		"    Then an outcome is observed\n"
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	steps := feature.Scenarios[0].Steps
	if len(steps) != 2 {
		t.Fatalf("got %d steps, want 2", len(steps))
	}
	assertStep(t, steps[0], gherkin.Given, "a precondition")
	assertStep(t, steps[1], gherkin.Then, "an outcome is observed")
}

func TestParseDescriptionWithKeywordLikeLines(t *testing.T) {
	input := `Feature: Workflow

  When a user makes a request, the system handles it.
  Then the system responds accordingly.

  Scenario: Actual scenario
    Given something happens
`
	feature, err := gherkin.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("got %d scenarios, want 1", len(feature.Scenarios))
	}
	if len(feature.Scenarios[0].Steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(feature.Scenarios[0].Steps))
	}
}

func TestParseErrorAndWithoutPreviousStep(t *testing.T) {
	input := `Feature: Bad
  Scenario: No prior step
    And orphan step
`
	_, err := gherkin.Parse(input)
	if err == nil {
		t.Fatal("expected error for And without previous step")
	}
}

func TestParseEmptyInput(t *testing.T) {
	_, err := gherkin.Parse("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func assertStep(t *testing.T, step gherkin.Step, wantType gherkin.StepType, wantText string) {
	t.Helper()
	if step.Type != wantType {
		t.Errorf("step.Type = %v, want %v", step.Type, wantType)
	}
	if step.Text != wantText {
		t.Errorf("step.Text = %q, want %q", step.Text, wantText)
	}
}
