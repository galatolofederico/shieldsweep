package notifications

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func GetNotificationRunner(typename string, config json.RawMessage) NotificationRunner {
	switch typename {
	case "command":
		var commandConfig CommandConfig
		err := json.Unmarshal(config, &commandConfig)
		if err != nil {
			panic(err)
		}
		return NewCommandRunner(commandConfig)
	default:
		panic(errors.Errorf("Notification %v does not exists", typename))
	}
}
