package notifications

import (
	"os/exec"

	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type CommandConfig struct {
	Command string
}

type CommandRunner struct {
	config CommandConfig
}

func NewCommandRunner(config CommandConfig) *CommandRunner {
	return &CommandRunner{config: config}
}

func (runner *CommandRunner) Notify(results []tools.ToolResult) error {
	cmd := exec.Command(runner.config.Command)
	err := cmd.Run()
	return err
}
