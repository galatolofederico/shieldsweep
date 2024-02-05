package main

import (
	"github.com/galatolofederico/shieldsweep/shsw/internal/engine"
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

func main() {
	engine := engine.NewEngine("./tmp")
	fake1 := tools.ToolConfig{
		Runner: tools.NewFakeTool(true, 1),
		Name:   "fake1",
	}
	fake2 := tools.ToolConfig{
		Runner: tools.NewFakeTool(false, 2),
		Name:   "fake2",
	}
	engine.AddToolConfig(fake1)
	engine.AddToolConfig(fake2)
	engine.Run()
}
