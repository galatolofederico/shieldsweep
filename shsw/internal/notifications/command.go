package notifications

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type CommandConfig struct {
	Uid       int
	Gid       int
	Shell     string
	ShellFlag string
	Command   string
}

type CommandRunner struct {
	config CommandConfig
}

func NewCommandRunner(config CommandConfig) *CommandRunner {
	if config.Shell == "" {
		config.Shell = "/bin/sh"
	}
	if config.ShellFlag == "" {
		config.ShellFlag = "-c"
	}
	return &CommandRunner{config: config}
}

func (runner *CommandRunner) Notify(results []tools.ToolResult) error {
	color.Green(fmt.Sprintf("Running command: %v %v (uid: %d gid: %d)\n", runner.config.Shell, runner.config.Command, runner.config.Uid, runner.config.Gid))
	cmd := exec.Command(runner.config.Shell, runner.config.ShellFlag, runner.config.Command)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(runner.config.Uid),
		Gid: uint32(runner.config.Gid),
	}

	output, err := cmd.Output()
	fmt.Println(err)
	fmt.Println(string(output))
	return err
}
