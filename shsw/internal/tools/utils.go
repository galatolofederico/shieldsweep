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
	case "rkhunter":
		var rkhunterConfig RKHunterConfig
		if config != nil {
			err := json.Unmarshal(config, &rkhunterConfig)
			if err != nil {
				panic(err)
			}
		}
		return NewRKHunter(rkhunterConfig)
	case "chkrootkit":
		var chkrootkitConfig ChkRootkitConfig
		if config != nil {
			err := json.Unmarshal(config, &chkrootkitConfig)
			if err != nil {
				panic(err)
			}
		}
		return NewChkRootkit(chkrootkitConfig)
	case "lynis":
		var lynisConfig LynisConfig
		if config != nil {
			err := json.Unmarshal(config, &lynisConfig)
			if err != nil {
				panic(err)
			}
		}
		return NewLynis(lynisConfig)
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
