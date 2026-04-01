package main

import (
	"fmt"
	"os"

	"github.com/swarm-forge/swarm-forge/internal/gherkin"
	"github.com/swarm-forge/swarm-forge/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gherkin-runner <feature-file>")
		os.Exit(1)
	}
	os.Exit(run(os.Args[1]))
}

func run(path string) int {
	feature, err := parseFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		return 1
	}

	registry := runner.NewRegistry()
	exec := runner.NewExecutor(registry)
	results := exec.Run(feature)
	return printResults(results)
}

func parseFile(path string) (gherkin.Feature, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return gherkin.Feature{}, err
	}
	return gherkin.Parse(string(data))
}

func printResults(results []runner.ScenarioResult) int {
	exitCode := 0
	for _, sr := range results {
		fmt.Printf("Scenario: %s [%s]\n", sr.Name, sr.Status)
		for _, step := range sr.Steps {
			printStep(step)
		}
		if sr.Status != runner.Passed {
			exitCode = 1
		}
	}
	return exitCode
}

func printStep(step runner.StepResult) {
	fmt.Printf("  %s [%s]", step.Text, step.Status)
	if step.Error != "" {
		fmt.Printf(" - %s", step.Error)
	}
	fmt.Println()
}
