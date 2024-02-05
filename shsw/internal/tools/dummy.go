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

func (runner *DummyToolRunner) Run(config Tool) error {
	fmt.Println("Running tool " + config.Name)
	time.Sleep(time.Duration(runner.config.Delay) * time.Second)
	var log string
	if runner.config.NewLog {
		fmt.Println("Tool " + config.Name + " is writing something new to log " + config.LogFile)
		log = "Success " + time.Now().String()
	} else {
		fmt.Println("Tool " + config.Name + " is writing the same thing to log " + config.LogFile)
		log = "Nothing New"
	}
	err := os.WriteFile(config.LogFile, []byte(log), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tool " + config.Name + " finished")
	return nil
}
