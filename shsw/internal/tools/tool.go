package tools

import (
	"encoding/json"
	"os"
	"time"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
)

const (
	Ready   = "ready"
	Running = "running"
	Queued  = "queued"
	Failed  = "failed"
)

type ToolState struct {
	LastRun     string
	LastLogHash string
	LastError   string
	State       string
}

type ToolRunner interface {
	Check() bool
	Run(config Tool) error
}

type Tool struct {
	State     ToolState
	Runner    ToolRunner
	Name      string
	LogFile   string
	StateFile string
}

type ToolResult struct {
	Name     string
	IsLogNew bool
	Error    error
}

func (config *Tool) Run(ch chan<- ToolResult) {
	result := ToolResult{Name: config.Name, IsLogNew: false, Error: nil}
	//TODO: implement log hash check
	utils.CheckPathForFile(config.LogFile)
	config.State.LastRun = time.Now().Format(time.RFC3339)
	config.State.State = Running
	err := config.Runner.Run(*config)

	if err != nil {
		config.State.LastError = err.Error()
		config.State.State = Failed
		result.Error = err
	} else {
		config.State.LastError = ""
		config.State.State = Ready
	}

	config.Save()
	ch <- result
}

func (config *Tool) Load() {
	utils.CheckPathForFile(config.StateFile)
	if _, err := os.Stat(config.StateFile); os.IsNotExist(err) {
		config.State = ToolState{
			LastRun:     "never",
			LastLogHash: "none",
			LastError:   "",
			State:       Ready,
		}
		config.Save()
	}
	dat, err := os.ReadFile(config.StateFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &config.State)
	if err != nil {
		panic(err)
	}
}

func (config *Tool) Save() {
	utils.CheckPathForFile(config.StateFile)
	encoded, err := json.Marshal(config.State)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(config.StateFile, encoded, 0644)
	if err != nil {
		panic(err)
	}
}
