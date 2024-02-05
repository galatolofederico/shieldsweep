package tools

import (
	"encoding/json"
	"os"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
)

type ToolState struct {
	LastRun     string
	LastLogHash string
}

type ToolRunner interface {
	Check() bool
	Run(config ToolConfig)
}

type ToolConfig struct {
	State     ToolState
	Runner    ToolRunner
	Name      string
	LogFile   string
	StateFile string
}

func (config *ToolConfig) Load() {
	utils.CheckPathForFile(config.StateFile)
	if _, err := os.Stat(config.StateFile); os.IsNotExist(err) {
		config.State = ToolState{LastRun: "never", LastLogHash: "none"}
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

func (config *ToolConfig) Save() {
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
