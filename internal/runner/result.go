package runner

// Status represents the outcome of a step or scenario.
type Status int

const (
	Passed  Status = iota
	Failed
	Pending
	Skipped
)

// String returns a human-readable status name.
func (s Status) String() string {
	names := [...]string{"Passed", "Failed", "Pending", "Skipped"}
	if int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}

// StepResult holds the outcome of running a single step.
type StepResult struct {
	Text   string
	Status Status
	Error  string
}

// ScenarioResult holds results for an entire scenario.
type ScenarioResult struct {
	Name   string
	Steps  []StepResult
	Status Status
}
