package start

import (
	"fmt"
	"io"

	"github.com/swarm-forge/swarm-forge/internal/banner"
	"github.com/swarm-forge/swarm-forge/internal/preflight"
	"github.com/swarm-forge/swarm-forge/internal/prompt"
	"github.com/swarm-forge/swarm-forge/internal/setup"
	"github.com/swarm-forge/swarm-forge/internal/tmux"
)

const window = "swarm"

// Config holds everything needed for the start sequence.
type Config struct {
	Commander        tmux.Commander
	Session          string
	ProjectRoot      string
	FS               setup.FS
	LookPath         preflight.LookPathFunc
	ConstitutionPath string
	Stdout           io.Writer
}

// Run performs the full startup sequence.
func Run(cfg Config) error {
	if err := runPreflight(cfg); err != nil {
		return err
	}
	if err := runSetup(cfg); err != nil {
		return err
	}
	constitution, err := readConstitution(cfg)
	if err != nil {
		return err
	}
	banner.Print(cfg.Stdout)
	if err := createSession(cfg); err != nil {
		return err
	}
	if err := writeAndLaunchAgents(cfg, constitution); err != nil {
		return err
	}
	return initMetricsPane(cfg)
}

func runPreflight(cfg Config) error {
	return preflight.Check(cfg.LookPath, "tmux", "claude", "watch")
}

func runSetup(cfg Config) error {
	if err := setup.EnsureDirs(cfg.FS, cfg.ProjectRoot); err != nil {
		return err
	}
	return setup.WriteHelperScripts(cfg.FS, cfg.ProjectRoot)
}

func readConstitution(cfg Config) (string, error) {
	path := cfg.ProjectRoot + "/" + cfg.ConstitutionPath
	data, err := cfg.FS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("constitution: %w", err)
	}
	return string(data), nil
}

func createSession(cfg Config) error {
	if cfg.Commander.HasSession(cfg.Session) {
		if err := tmux.KillSession(cfg.Commander, cfg.Session); err != nil {
			return err
		}
	}
	if err := tmux.CreateSession(cfg.Commander, cfg.Session, window); err != nil {
		return err
	}
	if err := tmux.SplitGrid(cfg.Commander, cfg.Session, window); err != nil {
		return err
	}
	return tmux.SetPaneTitles(cfg.Commander, cfg.Session, window)
}

func writeAndLaunchAgents(cfg Config, constitution string) error {
	promptsDir := cfg.ProjectRoot + "/.swarmforge/prompts"
	if err := cfg.FS.MkdirAll(promptsDir, 0o755); err != nil {
		return err
	}
	agents := []struct {
		pane         int
		name         string
		instructions string
	}{
		{0, "Architect", prompt.ArchitectInstructions},
		{1, "E2E-Interpreter", prompt.E2EInterpreterInstructions},
		{2, "Coder", prompt.CoderInstructions},
	}
	for _, a := range agents {
		acfg := prompt.AgentConfig{
			Role:         a.name,
			Instructions: a.instructions,
			Session:      cfg.Session,
			ProjectRoot:  cfg.ProjectRoot,
		}
		content := prompt.Build(acfg, constitution)
		promptFile := promptsDir + "/" + a.name + ".md"
		if err := cfg.FS.WriteFile(promptFile, []byte(content), 0o644); err != nil {
			return err
		}
		if err := tmux.LaunchAgent(cfg.Commander, cfg.Session, a.pane, a.name, promptFile, cfg.ProjectRoot); err != nil {
			return err
		}
	}
	return nil
}

func initMetricsPane(cfg Config) error {
	metricsCmd := "cd '" + cfg.ProjectRoot + "' && touch logs/agent_messages.log && tail -f logs/agent_messages.log"
	return tmux.SendKeys(cfg.Commander, cfg.Session, window, 3, metricsCmd)
}
