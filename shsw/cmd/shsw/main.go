package main

import (
	"github.com/galatolofederico/shieldsweep/shsw/internal/engine"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

func main() {
	e := engine.NewEngine()
	e.AddTool(tools.NewFakeTool("fake1", true, 1))
	e.AddTool(tools.NewFakeTool("fake2", false, 0))
	e.Run()
}
