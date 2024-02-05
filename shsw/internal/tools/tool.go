package tools

import (
	"encoding/json"
	"os"
	"time"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
)

type ToolState struct {
	LastRun     string
	LastLogHash string
	LastError   string
	Running     bool
	Failing     bool
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

func (config *Tool) Run() {
	utils.CheckPathForFile(config.LogFile)
	config.State.LastRun = time.Now().String()
	config.State.Running = true
	err := config.Runner.Run(*config)
	if err != nil {
		config.State.LastError = err.Error()
		config.State.Failing = true
	} else {
		config.State.LastError = ""
		config.State.Failing = false
	}
	config.State.Running = false
	config.Save()
}

func (config *Tool) Load() {
	utils.CheckPathForFile(config.StateFile)
	if _, err := os.Stat(config.StateFile); os.IsNotExist(err) {
		config.State = ToolState{
			LastRun:     "never",
			LastLogHash: "none",
			LastError:   "",
			Running:     false,
			Failing:     false,
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
