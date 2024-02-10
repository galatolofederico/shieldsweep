package tools

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func GetToolRunner(toolName string, config json.RawMessage) ToolRunner {
	// TODO: implement regex to match toolName
	// name-[0-9] so that one can use the same tool with different configurations
	// but with different names (name must be unique)
	switch toolName {
	case "dummy1", "dummy2", "dummy3":
		var dummyConfig DummyToolConfig
		err := json.Unmarshal(config, &dummyConfig)
		if err != nil {
			panic(err)
		}
		return NewDummyTool(dummyConfig)
	case "rkhunter":
		var rkhunterConfig RKHunterConfig
		if config != nil {
			err := json.Unmarshal(config, &rkhunterConfig)
			if err != nil {
				panic(err)
			}
		}
		return NewRKHunter(rkhunterConfig)
	default:
		panic(errors.Errorf("Tool %v not found", toolName))
	}
}
