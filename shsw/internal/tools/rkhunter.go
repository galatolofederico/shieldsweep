package tools

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
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

// TODO: create a temp file and use it as log file
// then check if the temp file has content
// if not means that there was an error
// rkhunter returns 1 if there is a warning but the scan is successful
// than remove from the logfile the lines containing dates and system info like kernel version
// finally save the content of the temp file in the actual log file
func (runner *RKHunterRunner) Run(tool Tool) error {
	cmd := exec.Command(
		runner.config.Path,
		"-c",
		"--sk",
		"--nocolors",
		"-l",
		tool.TempLogFile,
	)
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrapf(err, "Error running rkhunter: %v\n", string(output))
	} else {
		return nil
	}
}
