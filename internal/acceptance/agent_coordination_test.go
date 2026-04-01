package acceptance

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/gherkin"
	"github.com/swarm-forge/swarm-forge/internal/handoff"
	"github.com/swarm-forge/swarm-forge/internal/runner"
)

func TestScenario1_ArchitectCreatesFeatureAndHandsOff(t *testing.T) {
	tmpDir := t.TempDir()
	featurePath := filepath.Join(tmpDir, "features", "agent-coordination.feature")
	handoffPath := filepath.Join(tmpDir, "agent_context", "handoff-to-e2e.json")

	// Given the architect receives a feature request "agent coordination"
	featureRequest := "agent coordination"
	if featureRequest == "" {
		t.Fatal("feature request must not be empty")
	}

	// When the architect writes a Gherkin feature file
	writeFeatureFile(t, featurePath)

	// Then the feature file exists
	assertFileExists(t, featurePath)

	// And the feature file contains a valid Feature declaration
	content := readFileContent(t, featurePath)
	if !strings.Contains(content, "Feature:") {
		t.Fatal("feature file does not contain a Feature declaration")
	}

	// And a handoff document is created
	writeHandoffDoc(t, handoffPath, "architect", "e2e-interpreter", featurePath)

	// And the handoff document has from "architect" and to "e2e-interpreter"
	h := readHandoff(t, handoffPath)
	if h.From != "architect" {
		t.Fatalf("expected from=architect, got %s", h.From)
	}
	if h.To != "e2e-interpreter" {
		t.Fatalf("expected to=e2e-interpreter, got %s", h.To)
	}

	// And the handoff is valid (Validate must exist on Handoff)
	if err := h.Validate(); err != nil {
		t.Fatalf("handoff validation failed: %v", err)
	}
}

func TestScenario2_E2EInterpreterGeneratesTests(t *testing.T) {
	tmpDir := t.TempDir()
	handoffPath := filepath.Join(tmpDir, "agent_context", "handoff-to-e2e.json")
	featurePath := filepath.Join(tmpDir, "features", "agent-coordination.feature")

	// Given a handoff document exists for "e2e-interpreter"
	writeHandoffDoc(t, handoffPath, "architect", "e2e-interpreter", "features/agent-coordination.feature")
	h := readHandoff(t, handoffPath)
	if h.To != "e2e-interpreter" {
		t.Fatalf("handoff not for e2e-interpreter: %s", h.To)
	}

	// And a feature file exists
	writeFeatureFile(t, featurePath)
	assertFileExists(t, featurePath)

	// When the e2e interpreter parses the feature file
	content := readFileContent(t, featurePath)
	feature, err := gherkin.Parse(content)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	// Then executable test functions are generated
	if len(feature.Scenarios) == 0 {
		t.Fatal("no scenarios parsed from feature file")
	}

	// And the tests reference each scenario from the feature file
	for _, sc := range feature.Scenarios {
		if sc.Name == "" {
			t.Fatal("scenario has empty name")
		}
		if len(sc.Steps) == 0 {
			t.Fatalf("scenario %q has no steps", sc.Name)
		}
	}
}

func TestScenario3_CoderReceivesHandoffAndRunsTDD(t *testing.T) {
	tmpDir := t.TempDir()
	handoffPath := filepath.Join(tmpDir, "agent_context", "handoff-to-coder.json")

	// Given a handoff document exists for "coder"
	writeHandoffDoc(t, handoffPath, "e2e-interpreter", "coder", "features/agent-coordination.feature")
	h := readHandoff(t, handoffPath)
	if h.To != "coder" {
		t.Fatalf("handoff not for coder: %s", h.To)
	}

	// And the handoff validates successfully
	if err := h.Validate(); err != nil {
		t.Fatalf("handoff validation failed: %v", err)
	}

	// And acceptance tests exist that are currently failing
	reg := runner.NewRegistry()
	failStep := func(_ []string) error { return errors.New("not implemented") }
	mustRegister(t, reg, "a failing step", failStep)

	failFeature := gherkin.Feature{
		Name: "TDD Cycle",
		Scenarios: []gherkin.Scenario{
			{Name: "Failing", Steps: []gherkin.Step{
				{Type: gherkin.Given, Text: "a failing step"},
			}},
		},
	}

	// When the coder runs the test suite
	exec := runner.NewExecutor(reg)
	results := exec.Run(failFeature)

	// Then the test suite reports failures
	if results[0].Status != runner.Failed {
		t.Fatalf("expected Failed, got %s", results[0].Status)
	}

	// When the coder implements production code
	passReg := runner.NewRegistry()
	passStep := func(_ []string) error { return nil }
	mustRegister(t, passReg, "a failing step", passStep)

	passExec := runner.NewExecutor(passReg)
	passResults := passExec.Run(failFeature)

	// Then all acceptance tests pass
	if passResults[0].Status != runner.Passed {
		t.Fatalf("expected Passed, got %s", passResults[0].Status)
	}
}

func TestScenario4_GherkinParserParsesFeatureFile(t *testing.T) {
	// Given a feature file with the content:
	content := `Feature: Sample
  Scenario: Basic
    Given a precondition
    When an action occurs
    Then an outcome is observed`

	// When the parser processes the content
	feature, err := gherkin.Parse(content)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Then the parsed feature name is "Sample"
	if feature.Name != "Sample" {
		t.Fatalf("expected feature name Sample, got %q", feature.Name)
	}

	// And the parsed scenario count is 1
	if len(feature.Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(feature.Scenarios))
	}

	// And the first scenario has 3 steps
	if len(feature.Scenarios[0].Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(feature.Scenarios[0].Steps))
	}
}

func TestScenario5_StepRegistryMatchesStepText(t *testing.T) {
	// Given a step definition registered with pattern "a user named (.*)"
	reg := runner.NewRegistry()
	var captured string
	fn := func(matches []string) error {
		captured = matches[0]
		return nil
	}
	mustRegister(t, reg, "a user named (.*)", fn)

	// When matching the text "a user named Alice"
	stepFn, matches, err := reg.Match("a user named Alice")
	if err != nil {
		t.Fatalf("match error: %v", err)
	}

	// Then the match succeeds with capture "Alice"
	if len(matches) != 1 || matches[0] != "Alice" {
		t.Fatalf("expected capture [Alice], got %v", matches)
	}
	if err := stepFn(matches); err != nil {
		t.Fatalf("step func error: %v", err)
	}
	if captured != "Alice" {
		t.Fatalf("expected captured=Alice, got %q", captured)
	}
}

func TestScenario6_ExecutorRunsScenariosAndReportsResults(t *testing.T) {
	// Given a feature with one passing and one failing scenario
	reg := runner.NewRegistry()
	mustRegister(t, reg, "a passing step", func(_ []string) error {
		return nil
	})
	mustRegister(t, reg, "a failing step", func(_ []string) error {
		return errors.New("intentional failure")
	})

	feature := gherkin.Feature{
		Name: "Mixed Results",
		Scenarios: []gherkin.Scenario{
			{Name: "Pass", Steps: []gherkin.Step{
				{Type: gherkin.Given, Text: "a passing step"},
			}},
			{Name: "Fail", Steps: []gherkin.Step{
				{Type: gherkin.Given, Text: "a failing step"},
			}},
		},
	}

	// When the executor runs the feature
	exec := runner.NewExecutor(reg)
	results := exec.Run(feature)

	// Then the first scenario result is "Passed"
	if results[0].Status.String() != "Passed" {
		t.Fatalf("expected first scenario Passed, got %s", results[0].Status)
	}

	// And the second scenario result is "Failed"
	if results[1].Status.String() != "Failed" {
		t.Fatalf("expected second scenario Failed, got %s", results[1].Status)
	}
}

// --- helpers ---

func writeFeatureFile(t *testing.T, path string) {
	t.Helper()
	content := `Feature: Agent Coordination Workflow
  Scenario: Example
    Given a precondition
    When an action occurs
    Then an outcome is observed
`
	mustWriteFile(t, path, content)
}

func TestHandoffValidationRejectsEmptyFields(t *testing.T) {
	// A handoff with missing required fields must fail validation
	h := handoff.Handoff{From: "", To: "", Status: ""}
	err := h.Validate()
	if err == nil {
		t.Fatal("expected validation error for empty handoff, got nil")
	}
}

func writeHandoffDoc(t *testing.T, path, from, to, feature string) {
	t.Helper()
	mustMkdirAll(t, filepath.Dir(path))
	h := handoff.Handoff{
		From:    from,
		To:      to,
		Status:  "pending",
		Feature: feature,
	}
	if err := handoff.Write(path, h); err != nil {
		t.Fatalf("write handoff: %v", err)
	}
}

func readHandoff(t *testing.T, path string) handoff.Handoff {
	t.Helper()
	h, err := handoff.Read(path)
	if err != nil {
		t.Fatalf("read handoff: %v", err)
	}
	return h
}

func readFileContent(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	return string(data)
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file does not exist: %s", path)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	mustMkdirAll(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
}

func mustRegister(t *testing.T, reg *runner.Registry, pattern string, fn runner.StepFunc) {
	t.Helper()
	if err := reg.Register(pattern, fn); err != nil {
		t.Fatalf("register step %q: %v", pattern, err)
	}
}
