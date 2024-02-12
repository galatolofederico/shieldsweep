package tools

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ChkRootkitConfig struct {
	Path string
}

type ChkRootkitRunner struct {
	config ChkRootkitConfig
}

func sanitizeChkRootkitLog(log string) string {
	pattern := `.*\/tmp\/.*`
	re := regexp.MustCompile(pattern)
	log = re.ReplaceAllString(log, "")
	return strings.Replace(log, "\n\n", "\n", -1)
}

func NewChkRootkit(config ChkRootkitConfig) *ChkRootkitRunner {
	if config.Path == "" {
		config.Path = "/usr/bin/chkrootkit"
	}
	return &ChkRootkitRunner{config: config}
}

func (runner *ChkRootkitRunner) Check() bool {
	_, err := os.Stat(runner.config.Path)
	return !os.IsNotExist(err)
}

func (runner *ChkRootkitRunner) Run(tool Tool) error {
	if _, err := os.Stat(tool.CurrentLogFile); !os.IsNotExist(err) {
		os.Remove(tool.CurrentLogFile)
	}

	logFile, err := os.OpenFile(tool.CurrentLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "error creating temp log file")
	}

	cmd := exec.Command(
		runner.config.Path,
	)
	cmd.Dir = "/tmp"
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error running chkrootkit")
	}
	output, err := os.ReadFile(tool.CurrentLogFile)
	if err != nil {
		return errors.Wrap(err, "error reading temp log file")
	}
	output = []byte(sanitizeChkRootkitLog(string(output)))
	return os.WriteFile(tool.CurrentLogFile, output, 0644)
}
