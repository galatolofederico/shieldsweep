package tools

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/pkg/errors"
)

type RKHunterConfig struct {
	Path string
}

type RKHunterRunner struct {
	config RKHunterConfig
}

func sanitizeRKHunterLog(log string) string {
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
	if utils.FileExists(tool.CurrentLogFile) {
		os.Remove(tool.CurrentLogFile)
	}
	logFile, err := os.OpenFile(tool.CurrentLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "error creating temp log file")
	}

	tmpFile, err := os.CreateTemp("/tmp", "rkhunter")
	if err != nil {
		return errors.Wrap(err, "error creating temp file")
	}
	defer os.Remove(tmpFile.Name())
	os.Remove(tmpFile.Name())
	// its weird, i know but rkhunter always returns 1
	// the only way to know if it failed is to check if the log file was created
	// so i need to remove the temp file befor running it
	// and also defer the removal for actual cleanup

	cmd := exec.Command(
		runner.config.Path,
		"-c",
		"--sk",
		"--nocolors",
		"--logfile",
		tmpFile.Name(),
	)
	cmd.Dir = "/tmp"
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err = cmd.Run()
	output, errRead := os.ReadFile(tool.CurrentLogFile)
	if errRead != nil {
		return errors.Wrap(errRead, "error reading log file")
	}
	if !utils.FileExists(tmpFile.Name()) {
		return errors.Wrapf(err, "error running rkhunter. output: %v", string(output))
	}
	logFile.Close()
	output = []byte(sanitizeRKHunterLog(string(output)))
	return os.WriteFile(tool.CurrentLogFile, output, 0644)
}
