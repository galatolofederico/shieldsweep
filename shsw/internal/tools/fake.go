package tools

import (
	"fmt"
	"time"
)

type FakeToolRunner struct {
	exists bool
	sleep  int
}

func NewFakeTool(exists bool, sleep int) *FakeToolRunner {
	return &FakeToolRunner{
		exists: exists,
		sleep:  sleep,
	}
}

func (runner *FakeToolRunner) Check() bool {
	return runner.exists
}

func (runner *FakeToolRunner) Run(config ToolConfig) {
	fmt.Println("Running tool " + config.Name)
	time.Sleep(time.Duration(runner.sleep) * time.Second)
	fmt.Println("Tool " + config.Name + " finished")
}
