package main

import (
	"github.com/galatolofederico/shieldsweep/internal/engine"
	"github.com/galatolofederico/shieldsweep/internal/tools"
)

func main() {
	e := engine.NewEngine()
	e.AddTool(tools.NewFakeTool("fake1", true, 1))
	e.AddTool(tools.NewFakeTool("fake2", false, 0))
	e.Run()
}
