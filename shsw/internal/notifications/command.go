package notifications

import (
	"fmt"

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
	fmt.Println("Notifying with command: " + runner.config.Command)
	return nil
}
