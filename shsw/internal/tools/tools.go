package tools

import (
	"encoding/json"
	"os"
	"time"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/pkg/errors"
)

const (
	Ready    = "ready"
	Running  = "running"
	Queued   = "queued"
	Finished = "finished"
	Failed   = "failed"
)

type ToolState struct {
	LastRun       string
	LastLogChange string
	LastLogHash   string
	LastError     string
	State         string
}

type ToolRunner interface {
	Check() bool
	Run(config Tool) error
}

type ToolConfig struct {
	Name    string
	Enabled bool
	Config  json.RawMessage
}

type Tool struct {
	State       ToolState
	Runner      ToolRunner
	Name        string
	LogFile     string
	TempLogFile string
	StateFile   string
}

type ToolResult struct {
	Name     string
	IsLogNew bool
	Error    error
}

func (tool *Tool) Run(ch chan<- ToolResult) {
	result := ToolResult{Name: tool.Name, IsLogNew: false, Error: nil}

	utils.CheckPathForFile(tool.LogFile)
	utils.CheckPathForFile(tool.TempLogFile)
	tool.State.LastRun = time.Now().Format(time.RFC3339)
	tool.State.State = Running
	err := tool.Runner.Run(*tool)

	newLogHash := "none"
	if err == nil {
		exists := utils.FileExists(tool.LogFile)
		if exists {
			newLogHash, err = utils.SHA256File(tool.LogFile)
		} else {
			err = errors.Errorf("Log file not found: %v", tool.LogFile)
		}
	}

	if err != nil {
		tool.State.LastError = err.Error()
		tool.State.State = Failed
		result.Error = err
	} else {
		if newLogHash != tool.State.LastLogHash {
			tool.State.LastLogChange = time.Now().Format(time.RFC3339)
			tool.State.LastLogHash = newLogHash
			result.IsLogNew = true
		}
		tool.State.LastError = ""
		tool.State.State = Finished
	}

	tool.Save()
	ch <- result
}

func (tool *Tool) Load() {
	utils.CheckPathForFile(tool.StateFile)
	if _, err := os.Stat(tool.StateFile); os.IsNotExist(err) {
		tool.State = ToolState{
			LastRun:       "never",
			LastLogChange: "never",
			LastLogHash:   "none",
			LastError:     "",
			State:         Ready,
		}
		tool.Save()
	}
	dat, err := os.ReadFile(tool.StateFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &tool.State)
	if err != nil {
		panic(err)
	}
}

// TODO: se non riesce a salvare il file di stato deve davvero panicare?
func (tool *Tool) Save() {
	utils.CheckPathForFile(tool.StateFile)
	encoded, err := json.Marshal(tool.State)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(tool.StateFile, encoded, 0644)
	if err != nil {
		panic(err)
	}
}

func (tool *Tool) GetLog() string {
	logFile := tool.LogFile
	if tool.State.State == Running {
		logFile = tool.TempLogFile
	}
	utils.CheckPathForFile(logFile)
	dat, err := os.ReadFile(logFile)
	if err != nil {
		return "Log file not found"
	}
	return string(dat)
}

func (tool *Tool) GetLastError() string {
	return tool.State.LastError
}
