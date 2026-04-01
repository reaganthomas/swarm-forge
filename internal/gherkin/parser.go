package gherkin

import (
	"errors"
	"strings"
)

// Parse reads Gherkin text and returns a Feature AST.
func Parse(input string) (Feature, error) {
	lines := strings.Split(input, "\n")
	state := newParserState()

	for _, raw := range lines {
		if err := state.processLine(raw); err != nil {
			return Feature{}, err
		}
	}

	return state.build()
}

type parserState struct {
	feature      Feature
	current      *Scenario
	lastStepType *StepType
	hasFeature   bool
	inDocString  bool
}

func newParserState() *parserState {
	return &parserState{}
}

func (p *parserState) processLine(raw string) error {
	line := strings.TrimSpace(raw)
	if p.checkDocString(line) {
		return nil
	}
	if p.inDocString || isSkippable(line) {
		return nil
	}
	return p.parseLine(line)
}

func (p *parserState) checkDocString(line string) bool {
	if strings.HasPrefix(line, `"""`) {
		p.inDocString = !p.inDocString
		return true
	}
	return false
}

func isSkippable(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func (p *parserState) parseLine(line string) error {
	switch {
	case strings.HasPrefix(line, "Feature:"):
		return p.handleFeature(line)
	case strings.HasPrefix(line, "Scenario:"):
		return p.handleScenario(line)
	default:
		return p.handleStep(line)
	}
}

func (p *parserState) handleFeature(line string) error {
	p.feature.Name = strings.TrimSpace(strings.TrimPrefix(line, "Feature:"))
	p.hasFeature = true
	return nil
}

func (p *parserState) handleScenario(line string) error {
	p.finishScenario()
	name := strings.TrimSpace(strings.TrimPrefix(line, "Scenario:"))
	p.current = &Scenario{Name: name}
	p.lastStepType = nil
	return nil
}

func (p *parserState) handleStep(line string) error {
	if !isStepLine(line) || p.current == nil {
		return nil // description text or step-like text before first scenario
	}
	step, err := parseStep(line, p.lastStepType)
	if err != nil {
		return err
	}
	p.current.Steps = append(p.current.Steps, step)
	p.lastStepType = &step.Type
	return nil
}

func isStepLine(line string) bool {
	prefixes := []string{"Given ", "When ", "Then ", "And ", "But "}
	for _, p := range prefixes {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

func (p *parserState) finishScenario() {
	if p.current != nil {
		p.feature.Scenarios = append(p.feature.Scenarios, *p.current)
	}
}

func (p *parserState) build() (Feature, error) {
	if !p.hasFeature {
		return Feature{}, errors.New("no Feature found in input")
	}
	p.finishScenario()
	return p.feature, nil
}

func parseStep(line string, lastType *StepType) (Step, error) {
	stepType, text, ok := extractKeyword(line)
	if ok {
		return Step{Type: stepType, Text: text}, nil
	}
	return resolveInherited(line, lastType)
}

func extractKeyword(line string) (StepType, string, bool) {
	prefixes := []struct {
		keyword string
		st      StepType
	}{
		{"Given ", Given},
		{"When ", When},
		{"Then ", Then},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(line, p.keyword) {
			return p.st, strings.TrimPrefix(line, p.keyword), true
		}
	}
	return 0, "", false
}

func resolveInherited(line string, lastType *StepType) (Step, error) {
	prefix, found := inheritedPrefix(line)
	if !found {
		return Step{}, errors.New("unrecognized line: " + line)
	}
	if lastType == nil {
		return Step{}, errors.New("And/But without preceding step")
	}
	text := strings.TrimPrefix(line, prefix)
	return Step{Type: *lastType, Text: text}, nil
}

func inheritedPrefix(line string) (string, bool) {
	for _, p := range []string{"And ", "But "} {
		if strings.HasPrefix(line, p) {
			return p, true
		}
	}
	return "", false
}
