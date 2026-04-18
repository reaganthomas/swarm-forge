Feature: SwarmForge shell workspace setup

  The shell launcher prepares a project-local workspace, initializes git
  when needed, writes helper state, creates worktrees, and starts one
  tmux session per configured role.

  Scenario: Startup initializes a git repository when the working directory is not a repo
    Given a working directory without a ".git" directory
    When "swarmforge.sh" starts
    Then a git repository is initialized
    And the default branch is renamed to "master"
    And the first commit message is "Initial swarmforge repository"

  Scenario: Startup ensures the SwarmForge paths are ignored by git
    Given a working directory without a ".gitignore" file
    When "swarmforge.sh" initializes the repository
    Then ".gitignore" contains ".swarmforge/"
    And ".gitignore" contains ".worktrees/"
    And ".gitignore" contains "swarmtools/"
    And ".gitignore" contains "logs/"
    And ".gitignore" contains "agent_context/"

  Scenario: Startup prepares the local workspace directories
    Given a valid swarm configuration
    When "swarmforge.sh" prepares the workspace
    Then the directory "features" exists under the project root
    And the directory "logs" exists under the project root
    And the directory "agent_context" exists under the project root
    And the directory ".swarmforge" exists under the project root
    And the directory ".swarmforge/prompts" exists under the project root
    And the directory "swarmtools" exists under the project root
    And the directory ".worktrees" exists under the project root

  Scenario: Startup writes sessions metadata and a notify helper
    Given a valid swarm configuration
    When "swarmforge.sh" prepares the workspace
    Then the file ".swarmforge/sessions.tsv" exists
    And the file "swarmtools/notify-agent.sh" exists
    And "swarmtools/notify-agent.sh" is executable

  Scenario: Startup creates one git worktree per non-master role
    Given "swarmforge/swarmforge.conf" contains:
      """
      window architect claude master
      window coder codex coder
      window reviewer codex reviewer
      window logger none none
      """
    And the matching prompt files exist
    When "swarmforge.sh" prepares worktrees
    Then the worktree ".worktrees/coder" is created from "HEAD"
    And the worktree ".worktrees/reviewer" is created from "HEAD"
    And no worktree is created for "master"
    And no worktree is created for "none"

  Scenario: Existing worktrees are reused
    Given a valid swarm configuration with a "coder" worktree
    And ".worktrees/coder/.git" already exists
    When "swarmforge.sh" prepares worktrees
    Then the existing "coder" worktree is left in place

  Scenario: Existing swarm sessions are killed before startup continues
    Given a valid swarm configuration
    And tmux already has a session for one configured role
    When "swarmforge.sh" starts tmux sessions
    Then the existing session is killed before a replacement session is created

  Scenario: Startup creates one tmux session per configured role
    Given a valid swarm configuration
    When "swarmforge.sh" launches the swarm
    Then a tmux session named "swarmforge-architect" is created
    And a tmux session named "swarmforge-coder" is created
    And a tmux session named "swarmforge-reviewer" is created
    And a tmux session named "swarmforge-logger" is created
    And each session uses the window name "swarm"

  Scenario: Startup opens Terminal windows when osascript is available
    Given a valid swarm configuration
    And "osascript" is installed
    When "swarmforge.sh" finishes launching the swarm
    Then one Terminal window is opened for each session
    And the window ids are written to ".swarmforge/window-ids"

  Scenario: Startup attaches to the cleanup-owner session when osascript is unavailable
    Given a valid swarm configuration
    And "osascript" is not installed
    When "swarmforge.sh" finishes launching the swarm
    Then the current shell attaches to one running tmux session
