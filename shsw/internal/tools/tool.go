package tools

import (
	"encoding/json"
	"os"
	"time"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
)

const (
	Ready    = "ready"
	Running  = "running"
	Queued   = "queued"
	Finished = "finished"
	Failed   = "failed"
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

func (tool *Tool) Run(ch chan<- ToolResult) {
	result := ToolResult{Name: tool.Name, IsLogNew: false, Error: nil}
	//TODO: implement log hash check
	utils.CheckPathForFile(tool.LogFile)
	tool.State.LastRun = time.Now().Format(time.RFC3339)
	tool.State.State = Running
	err := tool.Runner.Run(*tool)

	if err != nil {
		tool.State.LastError = err.Error()
		tool.State.State = Failed
		result.Error = err
	} else {
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
			LastRun:     "never",
			LastLogHash: "none",
			LastError:   "",
			State:       Ready,
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
