Feature: SwarmForge shell agent launch

  The shell launcher generates startup instructions for each role,
  launches agent backends in their assigned worktrees, and provides a
  helper script for agent-to-agent messaging.

  Scenario: Startup writes one generated instruction file per agent-backed role
    Given "swarmforge/swarmforge.conf" contains:
      """
      window architect claude master
      window coder codex coder
      window reviewer codex reviewer
      window logger none none
      """
    And the matching prompt files exist
    When "swarmforge.sh" launches the roles
    Then the file ".swarmforge/prompts/architect.md" exists
    And the file ".swarmforge/prompts/coder.md" exists
    And the file ".swarmforge/prompts/reviewer.md" exists
    And no generated prompt file is required for "logger"

  Scenario: Generated instructions tell agents to read constitution and role prompts recursively
    Given "swarmforge/swarmforge.conf" contains:
      """
      window architect claude master
      """
    And "swarmforge/architect.prompt" exists
    When "swarmforge.sh" writes the generated instruction file for "architect"
    Then ".swarmforge/prompts/architect.md" contains "Read swarmforge/constitution.prompt"
    And ".swarmforge/prompts/architect.md" contains "read every file it refers to recursively"
    And ".swarmforge/prompts/architect.md" contains "Read swarmforge/architect.prompt"

  Scenario: Claude roles launch in their assigned worktrees with the generated prompt
    Given "swarmforge/swarmforge.conf" contains:
      """
      window architect claude master
      """
    And "swarmforge/architect.prompt" exists
    When "swarmforge.sh" launches the "architect" role
    Then tmux sends a command containing "cd '" to the architect pane
    And the command contains "claude --append-system-prompt-file"
    And the command contains "--permission-mode acceptEdits"
    And the command contains "\"$(cat '.swarmforge/prompts/architect.md')\""

  Scenario: Codex roles launch in their assigned worktrees with the generated prompt
    Given "swarmforge/swarmforge.conf" contains:
      """
      window coder codex coder
      """
    And "swarmforge/coder.prompt" exists
    When "swarmforge.sh" launches the "coder" role
    Then tmux sends a command containing "codex -C" to the coder pane
    And the command contains "\"$(cat '.swarmforge/prompts/coder.md')\""

  Scenario: The logger role tails the shared agent log without an agent backend
    Given "swarmforge/swarmforge.conf" contains:
      """
      window logger none none
      """
    When "swarmforge.sh" launches the "logger" role
    Then tmux sends a command containing "touch logs/agent_messages.log"
    And the command contains "tail -f logs/agent_messages.log"

  Scenario: The notify helper routes messages by role or index
    Given ".swarmforge/sessions.tsv" contains session rows for "architect" and "coder"
    When "swarmtools/notify-agent.sh" sends a message to "architect"
    Then the message is sent to the architect session
    When "swarmtools/notify-agent.sh" sends a message to "2"
    Then the message is sent to the coder session

  Scenario: The notify helper logs the message before sending it to tmux
    Given ".swarmforge/sessions.tsv" contains a session row for "architect"
    When "swarmtools/notify-agent.sh" sends the message "hello architect"
    Then "logs/agent_messages.log" receives a timestamped entry for that session
    And tmux sends the literal message text to pane "0.0"

  Scenario: The cleanup owner appends shutdown cleanup to its launch command
    Given a valid swarm configuration
    When "swarmforge.sh" chooses the cleanup owner
    Then the first configured window becomes the cleanup owner
    And its launch command includes "swarm-cleanup.sh"
