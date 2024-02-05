package tools

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func GetToolRunner(toolName string, config json.RawMessage) ToolRunner {
	switch toolName {
	case "dummy1", "dummy2", "dummy3":
		var dummyConfig DummyToolConfig
		err := json.Unmarshal(config, &dummyConfig)
		if err != nil {
			panic(err)
		}
		return NewDummyTool(dummyConfig)
	default:
		panic(errors.Errorf("Tool %v not found", toolName))
	}
}
