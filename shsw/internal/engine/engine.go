package engine

import (
	"fmt"
	"path/filepath"

	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type Engine struct {
	home  string
	tools []tools.Tool
}

func NewEngine(home string) *Engine {
	return &Engine{home: home, tools: []tools.Tool{}}
}

func (e *Engine) AddTool(t tools.Tool) {
	if t.Check() {
		t.SetLogFile(filepath.Join(e.home, t.GetName(), "logs", "log.txt"))
		t.SetStateFile(filepath.Join(e.home, t.GetName(), "state", "state.json"))
		t.Load()
		e.tools = append(e.tools, t)
		fmt.Println("Added tool " + t.GetName())
	} else {
		fmt.Println("Tool " + t.GetName() + " not added")
	}
}

func (e *Engine) Run() {
	for _, t := range e.tools {
		t.Run()
	}
}
