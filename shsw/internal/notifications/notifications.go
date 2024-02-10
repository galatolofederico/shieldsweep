package notifications

import (
	"encoding/json"

	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type NotificationConfig struct {
	Type   string
	Config json.RawMessage
}

type NotificationRunner interface {
	Notify(results []tools.ToolResult) error
}

type Notification struct {
	Type   string
	Runner NotificationRunner
}

func (noitification *Notification) Notify(results []tools.ToolResult) error {
	return noitification.Runner.Notify(results)
}
