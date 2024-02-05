package engine

import (
	"fmt"
	"path/filepath"

	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type Engine struct {
	home         string
	toolsConfigs []tools.ToolConfig
}

func NewEngine(home string) *Engine {
	return &Engine{home: home, toolsConfigs: []tools.ToolConfig{}}
}

func (engine *Engine) AddToolConfig(config tools.ToolConfig) {
	if config.Runner.Check() {
		config.LogFile = filepath.Join(engine.home, config.Name, "logs", "log.txt")
		config.StateFile = filepath.Join(engine.home, config.Name, "state", "state.json")
		config.Load()
		engine.toolsConfigs = append(engine.toolsConfigs, config)
		fmt.Println("Added tool " + config.Name)
	} else {
		fmt.Println("Tool " + config.Name + " not added")
	}
}

func (engine *Engine) Run() {
	for _, config := range engine.toolsConfigs {
		config.Runner.Run(config)
	}
}
