package engine

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
	"github.com/pkg/errors"
)

type EngineConfig struct {
	Tools []string
}

type Engine struct {
	home         string
	toolsConfigs []tools.ToolConfig
	config       EngineConfig
}

func NewEngine(home string) *Engine {
	config := EngineConfig{}
	configFile := filepath.Join(home, "shsw.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(errors.Wrapf(err, "Error reading config file: %v\n", configFile))
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	engine := &Engine{home: home, toolsConfigs: []tools.ToolConfig{}, config: config}
	for _, toolName := range config.Tools {
		runner := tools.GetToolRunner(toolName)
		if runner.Check() {
			color.Green("[+] Tool " + toolName + " found")
			toolConfig := tools.ToolConfig{
				State:     tools.ToolState{LastRun: "never", LastLogHash: "none"},
				Runner:    runner,
				Name:      toolName,
				LogFile:   filepath.Join(home, toolName, "logs", "log.txt"),
				StateFile: filepath.Join(home, toolName, "state", "state.json"),
			}
			toolConfig.Load()
			engine.toolsConfigs = append(engine.toolsConfigs, toolConfig)
		} else {
			color.Red("[!] Tool " + toolName + " not found")
		}
	}
	return engine
}

func (engine *Engine) Run() {
	for _, config := range engine.toolsConfigs {
		config.Run()
	}
}
