package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
	"github.com/pkg/errors"
)

type EngineToolConfig struct {
	Name   string
	Config json.RawMessage
}

type EngineConfig struct {
	Parallelism int
	Tools       []EngineToolConfig
}

type Engine struct {
	home   string
	tools  []tools.Tool
	config EngineConfig
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

	engine := &Engine{home: home, tools: []tools.Tool{}, config: config}
	for _, config := range config.Tools {
		runner := tools.GetToolRunner(config.Name, config.Config)
		if runner.Check() {
			color.Green("[+] Tool " + config.Name + " found")
			toolConfig := tools.Tool{
				State:     tools.ToolState{LastRun: "never", LastLogHash: "none"},
				Runner:    runner,
				Name:      config.Name,
				LogFile:   filepath.Join(home, config.Name, "logs", "log.txt"),
				StateFile: filepath.Join(home, config.Name, "state", "state.json"),
			}
			toolConfig.Load()
			//TODO: check if tool already exists
			engine.tools = append(engine.tools, toolConfig)
		} else {
			color.Red("[!] Tool " + config.Name + " not found")
		}
	}
	return engine
}

func (engine *Engine) Run() {
	works := make(chan tools.Tool)
	results := make(chan tools.ToolResult)

	go func() {
		for _, tool := range engine.tools {
			works <- tool
		}
		close(works)
	}()

	var wg sync.WaitGroup

	for i := 0; i < engine.config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for config := range works {
				config.Run(results)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}
}
