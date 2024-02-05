package tools

import (
	"fmt"
	"time"
)

type DummyToolRunner struct {
	exists bool
	sleep  int
}

func NewDummyTool(exists bool, sleep int) *DummyToolRunner {
	return &DummyToolRunner{
		exists: exists,
		sleep:  sleep,
	}
}

func (runner *DummyToolRunner) Check() bool {
	return runner.exists
}

func (runner *DummyToolRunner) Run(config ToolConfig) {
	fmt.Println("Running tool " + config.Name)
	time.Sleep(time.Duration(runner.sleep) * time.Second)
	fmt.Println("Tool " + config.Name + " finished")
}
