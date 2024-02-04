package engine

import (
	"fmt"

	"github.com/galatolofederico/shieldsweep/internal/tools"
)

type Engine struct {
	tools []tools.Tool
}

func NewEngine() *Engine {
	return &Engine{tools: []tools.Tool{}}
}

func (e *Engine) AddTool(t tools.Tool) {
	if t.Check() {
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
