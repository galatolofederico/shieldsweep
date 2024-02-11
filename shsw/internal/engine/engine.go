package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/galatolofederico/shieldsweep/shsw/internal/notifications"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
	"github.com/pkg/errors"
)

type EngineConfig struct {
	Parallelism   int
	Notifications []notifications.NotificationConfig
	Tools         []tools.ToolConfig
}

type Engine struct {
	home          string
	running       bool
	startedAt     string
	tools         []tools.Tool
	notifications []notifications.Notification
	config        EngineConfig
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
	color.Green(fmt.Sprintf("[+] Parallelism level: %v\n", config.Parallelism))

	engine := &Engine{
		home:          home,
		tools:         []tools.Tool{},
		notifications: []notifications.Notification{},
		config:        config,
	}
	for _, config := range config.Tools {
		if !config.Enabled {
			color.Yellow("[!] Tool " + config.Name + " disabled")
			continue
		}
		runner := tools.GetToolRunner(config.Name, config.Config)
		if runner.Check() {
			color.Green("[+] Tool " + config.Name + " found")
			toolConfig := tools.Tool{
				State:       tools.ToolState{LastRun: "never", LastLogHash: "none"},
				Runner:      runner,
				Name:        config.Name,
				LogFile:     filepath.Join(home, config.Name, "logs", "log.txt"),
				OldLogFile:  filepath.Join(home, config.Name, "logs", "old.txt"),
				TempLogFile: filepath.Join(home, config.Name, "logs", "tmp.txt"),
				StateFile:   filepath.Join(home, config.Name, "state", "state.json"),
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
	for _, config := range config.Notifications {
		runner := notifications.GetNotificationRunner(config.Type, config.Config)
		color.Green("[+] Notification '" + config.Type + "' loaded")
		notificationConfig := notifications.Notification{
			Type:   config.Type,
			Runner: runner,
		}
		engine.notifications = append(engine.notifications, notificationConfig)
	}

	if len(engine.tools) == 0 {
		panic(errors.New("No tools available"))
	}
	return engine
}

func (engine *Engine) Run() []tools.ToolResult {
	color.Green("[!] Running scan")
	engine.running = true
	engine.startedAt = time.Now().Format(time.RFC3339)

	runResults := []tools.ToolResult{}
	works := make(chan string)
	results := make(chan tools.ToolResult)

	for i := range engine.tools {
		engine.tools[i].State.State = tools.Queued
	}

	go func() {
		for _, tool := range engine.tools {
			works <- tool.Name
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

	for i := range engine.tools {
		if engine.tools[i].State.State != tools.Failed {
			engine.tools[i].State.State = tools.Ready
			engine.tools[i].Save()
		}
	}

	shouldNotify := false
	for _, result := range runResults {
		if result.IsLogNew {
			shouldNotify = true
			break
		}
	}

	if shouldNotify {
		for _, notification := range engine.notifications {
			color.Green("[+] Notifying with " + notification.Type)
			notification.Notify(runResults)
		}
	}

	color.Green("[!] Scan finished")
	return runResults
}

func (engine *Engine) GetToolStates() []messages.ToolStateReply {
	ret := []messages.ToolStateReply{}
	for _, tool := range engine.tools {
		ret = append(ret, messages.ToolStateReply{
			Name:          tool.Name,
			State:         tool.State.State,
			LastRun:       tool.State.LastRun,
			LastLogChange: tool.State.LastLogChange,
			LastError:     tool.State.LastError,
		})
	}
	return ret
}

func (engine *Engine) IsRunning() bool {
	return engine.running
}

func (engine *Engine) GetStartedAt() string {
	return engine.startedAt
}
