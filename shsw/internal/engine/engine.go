package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	home      string
	running   bool
	startedAt string
	tools     []tools.Tool
	config    EngineConfig
}

func (engine *Engine) GetTool(name string) *tools.Tool {
	for i, tool := range engine.tools {
		if tool.Name == name {
			return &engine.tools[i]
		}
	}
	return nil
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
			if engine.GetTool(config.Name) != nil {
				panic(errors.Errorf("Duplicate tool name: %v\n", config.Name))
			}
			engine.tools = append(engine.tools, toolConfig)
		} else {
			color.Red("[!] Tool " + config.Name + " not found")
		}
	}
	return engine
}

func (engine *Engine) Run() []tools.ToolResult {
	engine.running = true
	engine.startedAt = time.Now().Format(time.RFC3339)

	runResults := []tools.ToolResult{}
	works := make(chan string)
	results := make(chan tools.ToolResult)

	go func() {
		for i, tool := range engine.tools {
			works <- tool.Name
			engine.tools[i].State.Queued = true
		}
		close(works)
	}()

	var wg sync.WaitGroup

	for i := 0; i < engine.config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for config := range works {
				tool := engine.GetTool(config)
				tool.Run(results)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		runResults = append(runResults, res)
	}

	engine.startedAt = ""
	engine.running = false
	return runResults
}

func (engine *Engine) GetToolStates() []tools.ToolState {
	ret := []tools.ToolState{}
	for _, tool := range engine.tools {
		ret = append(ret, tool.State)
	}
	return ret
}

func (engine *Engine) IsRunning() bool {
	return engine.running
}

func (engine *Engine) GetStartedAt() string {
	return engine.startedAt
}
