package tools

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

type RKHunterConfig struct {
	Path string
}

type RKHunterRunner struct {
	config RKHunterConfig
}

func sanitizeLog(log string) string {
	pattern := `(?m)^The system checks took:.*$`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(log, "")
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

func (runner *RKHunterRunner) Run(tool Tool) error {
	if _, err := os.Stat(tool.TempLogFile); !os.IsNotExist(err) {
		os.Remove(tool.TempLogFile)
	}
	cmd := exec.Command(
		runner.config.Path,
		"-c",
		"--sk",
		"--nocolors",
		"--logfile",
		tool.TempLogFile,
	)
	output, _ := cmd.Output()
	if _, err := os.Stat(tool.TempLogFile); os.IsNotExist(err) {
		return errors.Errorf("%v", output)
	}
	output = []byte(sanitizeLog(string(output)))
	return os.WriteFile(tool.LogFile, output, 0644)
}
