package notifications

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type CommandConfig struct {
	Uid     int
	Gid     int
	Command []string
}

type CommandRunner struct {
	config CommandConfig
}

func NewCommandRunner(config CommandConfig) *CommandRunner {
	return &CommandRunner{config: config}
}

func (runner *CommandRunner) Notify(results []tools.ToolResult) error {
	color.Green(fmt.Sprintf("Running command: %v (uid: %d gid: %d)\n", runner.config.Command, runner.config.Uid, runner.config.Gid))
	cmd := exec.Command(runner.config.Command[0], runner.config.Command[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(runner.config.Uid),
		Gid: uint32(runner.config.Gid),
	}

	_, err := cmd.Output()
	return err
}
