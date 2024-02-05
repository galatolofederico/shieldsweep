package tools

import (
	"fmt"
	"math/rand"
	"os"
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

func (runner *DummyToolRunner) Run(config ToolConfig) error {
	fmt.Println("Running tool " + config.Name)
	time.Sleep(time.Duration(runner.sleep) * time.Second)
	if rand.Float64() < 0.75 {
		fmt.Println("Tool " + config.Name + " is writing log " + config.LogFile)
		log := "Success " + time.Now().String()
		err := os.WriteFile(config.LogFile, []byte(log), 0644)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Tool " + config.Name + " finished")
	return nil
}
