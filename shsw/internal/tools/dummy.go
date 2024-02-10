package tools

import (
	"fmt"
	"os"
	"time"
)

type DummyToolConfig struct {
	Available bool
	Delay     int
	NewLog    bool
}

type DummyToolRunner struct {
	config DummyToolConfig
}

func NewDummyTool(config DummyToolConfig) *DummyToolRunner {
	return &DummyToolRunner{config: config}
}

func (runner *DummyToolRunner) Check() bool {
	return runner.config.Available
}

func (runner *DummyToolRunner) Run(tool Tool) error {
	fmt.Println("Running tool " + tool.Name)
	time.Sleep(time.Duration(runner.config.Delay) * time.Second)
	var log string
	if runner.config.NewLog {
		fmt.Println("Tool " + tool.Name + " is writing something new to log " + tool.LogFile)
		log = "Success " + time.Now().String()
	} else {
		fmt.Println("Tool " + tool.Name + " is writing the same thing to log " + tool.LogFile)
		log = "Nothing New"
	}
	err := os.WriteFile(tool.LogFile, []byte(log), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tool " + tool.Name + " finished")
	return nil
}
