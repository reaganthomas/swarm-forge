Feature: Agent Coordination Workflow

  The swarm coordinates work through a structured handoff protocol.
  When a user makes a feature request, the Architect writes a Gherkin
  feature file and hands off to the E2E Interpreter, who writes
  acceptance tests, then hands off to the Coder for TDD.

  Scenario: Architect creates feature file and hands off to E2E Interpreter
    Given the architect receives a feature request "agent coordination"
    When the architect writes a Gherkin feature file to "features/agent-coordination.feature"
    Then the feature file exists at "features/agent-coordination.feature"
    And the feature file contains a valid Feature declaration
    And a handoff document is created at "agent_context/handoff-to-e2e.json"
    And the handoff document has from "architect" and to "e2e-interpreter"

  Scenario: E2E Interpreter generates acceptance tests from feature file
    Given a handoff document exists at "agent_context/handoff-to-e2e.json" for "e2e-interpreter"
    And a feature file exists at "features/agent-coordination.feature"
    When the e2e interpreter parses the feature file
    Then executable test functions are generated
    And the tests reference each scenario from the feature file

  Scenario: Coder receives handoff and runs TDD cycle
    Given a handoff document exists at "agent_context/handoff-to-coder.json" for "coder"
    And acceptance tests exist that are currently failing
    When the coder runs the test suite
    Then the test suite reports failures
    When the coder implements production code
    Then all acceptance tests pass

  Scenario: Gherkin parser correctly parses a feature file
    Given a feature file with the content:
      """
      Feature: Sample
        Scenario: Basic
          Given a precondition
          When an action occurs
          Then an outcome is observed
      """
    When the parser processes the content
    Then the parsed feature name is "Sample"
    And the parsed scenario count is 1
    And the first scenario has 3 steps

  Scenario: Step registry matches step text to definitions
    Given a step definition registered with pattern "a user named (.*)"
    When matching the text "a user named Alice"
    Then the match succeeds with capture "Alice"

  Scenario: Executor runs scenarios and reports results
    Given a feature with one passing and one failing scenario
    When the executor runs the feature
    Then the first scenario result is "Passed"
    And the second scenario result is "Failed"
