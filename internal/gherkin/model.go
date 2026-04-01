package gherkin

// StepType classifies a Gherkin step keyword.
type StepType int

const (
	Given StepType = iota
	When
	Then
)

// String returns the keyword for a StepType.
func (s StepType) String() string {
	names := [...]string{"Given", "When", "Then"}
	if int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}

// Step represents a single Given/When/Then step.
type Step struct {
	Type StepType
	Text string
}

// Scenario represents a single scenario within a feature.
type Scenario struct {
	Name  string
	Steps []Step
}

// Feature represents a parsed .feature file.
type Feature struct {
	Name      string
	Scenarios []Scenario
}
