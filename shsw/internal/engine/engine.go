package engine

import (
	"encoding/json"
	"fmt"
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
	home   string
	tools  []tools.Tool
	config EngineConfig
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

func (engine *Engine) Run() {
	works := make(chan string)
	results := make(chan tools.ToolResult)

	go func() {
		for _, tool := range engine.tools {
			works <- tool.Name
		}
		close(works)
	}()

	var wg sync.WaitGroup
	shared := []int{1, 2, 3}

	for i := 0; i < engine.config.Parallelism; i++ {
		wg.Add(1)
		go func(Engine *Engine) {
			defer wg.Done()
			for config := range works {
				tool := engine.GetTool(config)
				tool.Run(results)
				shared[0] = 100

			}
		}(engine)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Printf("Engine is running %p", &engine)
	go func(engine *Engine) {
		for {
			for _, tool := range engine.tools {
				t := engine.GetTool(tool.Name)
				fmt.Println("SHARED", shared)
				fmt.Printf("(FROM GORUTINE) Tool %v is at addres %p\n", t.Name, t)
				fmt.Printf("(FROM GORUTINE) Tool %v is running at address %p and values is %v \n", t.Name, &t.State.Running, t.State.Running)
			}
			time.Sleep(time.Millisecond * 100)
		}

	}(engine)

	for res := range results {
		fmt.Println(res)
	}
}
