package tools

import (
	"fmt"
	"time"
)

type FakeTool struct {
	*DefaultTool
	exists bool
	sleep  int
}

func NewFakeTool(name string, exists bool, sleep int) *FakeTool {
	return &FakeTool{
		DefaultTool: &DefaultTool{name: name, logFile: "", stateFile: ""},
		exists:      exists,
		sleep:       sleep,
	}
}

func (f FakeTool) Check() bool {
	return f.exists
}

func (f FakeTool) Run() {
	fmt.Println("Running tool " + f.name)
	time.Sleep(time.Duration(f.sleep) * time.Second)
	fmt.Println("Tool " + f.name + " finished")
}
