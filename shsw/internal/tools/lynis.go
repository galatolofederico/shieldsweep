package tools

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type LynisConfig struct {
	Path string
}

type LynisRunner struct {
	config LynisConfig
}

func sanitizeLynisLog(log string) string {
	lynisRegex := regexp.MustCompile(`Lynis .*`)
	programVersionRegex := regexp.MustCompile(`Program version:.*`)
	kernelVersionRegex := regexp.MustCompile(`Kernel version:.*`)
	lastTimeSyncRegex := regexp.MustCompile(`.*Last time synchronization.*\n?`)
	longExecutionRegex := regexp.MustCompile(`.*had a long execution.*\n?`)

	log = lynisRegex.ReplaceAllString(log, "Lynis x.x.x")
	log = programVersionRegex.ReplaceAllString(log, "Program version: x.x.x")
	log = kernelVersionRegex.ReplaceAllString(log, "Kernel version: x.x.x")
	log = lastTimeSyncRegex.ReplaceAllString(log, "")
	log = longExecutionRegex.ReplaceAllString(log, "")

	return strings.Replace(log, "\n\n", "\n", -1)
}

func NewLynis(config LynisConfig) *LynisRunner {
	if config.Path == "" {
		config.Path = "/usr/bin/lynis"
	}
	return &LynisRunner{config: config}
}

func (runner *LynisRunner) Check() bool {
	_, err := os.Stat(runner.config.Path)
	return !os.IsNotExist(err)
}

func (runner *LynisRunner) Run(tool Tool) error {
	if _, err := os.Stat(tool.TempLogFile); !os.IsNotExist(err) {
		os.Remove(tool.TempLogFile)
	}

	logFile, err := os.OpenFile(tool.TempLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "error creating temp log file")
	}

	cmd := exec.Command(
		runner.config.Path,
		"audit",
		"system",
		"--no-colors",
	)
	cmd.Dir = "/tmp"
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error running Lynis")
	}
	output, err := os.ReadFile(tool.TempLogFile)
	if err != nil {
		return errors.Wrap(err, "error reading temp log file")
	}
	output = []byte(sanitizeLynisLog(string(output)))
	return os.WriteFile(tool.LogFile, output, 0644)
}