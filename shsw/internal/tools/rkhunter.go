package tools

import (
	"os"
	"os/exec"
)

type RKHunterConfig struct {
	Path string
}

type RKHunterRunner struct {
	config RKHunterConfig
}

func NewRKHunter(config RKHunterConfig) *RKHunterRunner {
	if config.Path == "" {
		config.Path = "/usr/bin/rkhunter"
	}
	return &RKHunterRunner{config: config}
}

func (runner *RKHunterRunner) Check() bool {
	_, err := os.Stat(runner.config.Path)
	return !os.IsNotExist(err)
}

func (runner *RKHunterRunner) Run(config Tool) error {
	cmd := exec.Command(
		runner.config.Path,
		"-sk",
		"-l",
		config.LogFile,
	)
	_, err := cmd.Output()
	if err != nil {
		return err
	} else {
		return nil
	}
}
