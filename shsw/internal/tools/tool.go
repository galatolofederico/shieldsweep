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

type Tool interface {
	SetLogFile(string)
	GetLogFile() string
	SetStateFile(string)
	GetStateFile() string
	GetName() string
	Load()
	Save()
	Check() bool
	Run()
}

type DefaultTool struct {
	state     ToolState
	name      string
	logFile   string
	stateFile string
}

func (t *DefaultTool) Load() {
	utils.CheckPathForFile(t.stateFile)
	if _, err := os.Stat(t.stateFile); os.IsNotExist(err) {
		t.state = ToolState{LastRun: "never", LastLogHash: "none"}
		t.Save()
	}
	dat, err := os.ReadFile(t.stateFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &t.state)
	if err != nil {
		panic(err)
	}
}

func (t *DefaultTool) Save() {
	utils.CheckPathForFile(t.stateFile)
	encoded, err := json.Marshal(t.state)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(t.stateFile, encoded, 0644)
	if err != nil {
		panic(err)
	}
}

func (t *DefaultTool) GetName() string {
	return t.name
}

func (t *DefaultTool) SetLogFile(file string) {
	utils.CheckPathForFile(file)
	t.logFile = file
}

func (t *DefaultTool) GetLogFile() string {
	return t.logFile
}

func (t *DefaultTool) SetStateFile(file string) {
	utils.CheckPathForFile(file)
	t.stateFile = file
}

func (t *DefaultTool) GetStateFile() string {
	return t.stateFile
}
