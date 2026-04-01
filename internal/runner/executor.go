package runner

import "github.com/swarm-forge/swarm-forge/internal/gherkin"

// Executor runs Feature scenarios against a Registry.
type Executor struct {
	registry *Registry
}

// NewExecutor creates an Executor with the given Registry.
func NewExecutor(registry *Registry) *Executor {
	return &Executor{registry: registry}
}

// Run executes all scenarios in a Feature and returns results.
func (e *Executor) Run(feature gherkin.Feature) []ScenarioResult {
	results := make([]ScenarioResult, len(feature.Scenarios))
	for i, sc := range feature.Scenarios {
		results[i] = e.runScenario(sc)
	}
	return results
}

func (e *Executor) runScenario(sc gherkin.Scenario) ScenarioResult {
	result := ScenarioResult{Name: sc.Name, Status: Passed}
	failed := false

	for _, step := range sc.Steps {
		sr := e.runStep(step, failed)
		result.Steps = append(result.Steps, sr)
		failed = failed || sr.Status != Passed
		result.Status = worstStatus(result.Status, sr.Status)
	}

	return result
}

func (e *Executor) runStep(step gherkin.Step, skipRemaining bool) StepResult {
	if skipRemaining {
		return StepResult{Text: step.Text, Status: Skipped}
	}
	return e.executeStep(step)
}

func (e *Executor) executeStep(step gherkin.Step) StepResult {
	fn, matches, err := e.registry.Match(step.Text)
	if err != nil {
		return StepResult{Text: step.Text, Status: Pending}
	}
	if err := fn(matches); err != nil {
		return StepResult{Text: step.Text, Status: Failed, Error: err.Error()}
	}
	return StepResult{Text: step.Text, Status: Passed}
}

func worstStatus(a, b Status) Status {
	pa := statusPriority(a)
	pb := statusPriority(b)
	if pa >= pb {
		return a
	}
	return b
}

func statusPriority(s Status) int {
	priorities := [...]int{0, 3, 2, 1}
	if int(s) < len(priorities) {
		return priorities[s]
	}
	return 0
}
