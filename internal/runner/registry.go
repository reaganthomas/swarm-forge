package runner

import (
	"errors"
	"fmt"
	"regexp"
)

// StepFunc executes a step. Captures from the regex are passed in.
type StepFunc func(matches []string) error

type stepDef struct {
	pattern *regexp.Regexp
	fn      StepFunc
}

// Registry holds step definitions mapped to regex patterns.
type Registry struct {
	steps []stepDef
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a step definition with the given regex pattern.
func (r *Registry) Register(pattern string, fn StepFunc) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	r.steps = append(r.steps, stepDef{pattern: re, fn: fn})
	return nil
}

// Match finds the step definition matching the given text.
func (r *Registry) Match(text string) (StepFunc, []string, error) {
	matched := r.findMatches(text)
	return selectMatch(matched, text)
}

type matchResult struct {
	def     stepDef
	captures []string
}

func (r *Registry) findMatches(text string) []matchResult {
	var matched []matchResult
	for _, sd := range r.steps {
		subs := sd.pattern.FindStringSubmatch(text)
		if subs == nil {
			continue
		}
		matched = append(matched, matchResult{
			def:      sd,
			captures: subs[1:],
		})
	}
	return matched
}

func selectMatch(matched []matchResult, text string) (StepFunc, []string, error) {
	if len(matched) == 0 {
		return nil, nil, errors.New("no step definition matches: " + text)
	}
	if len(matched) > 1 {
		return nil, nil, errors.New("ambiguous step definitions for: " + text)
	}
	return matched[0].def.fn, matched[0].captures, nil
}
