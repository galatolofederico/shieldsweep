package messages

import (
	"github.com/galatolofederico/shieldsweep/shsw/internal/tools"
)

type RunReply struct {
	Started bool
	Message string
}

type StatusReply struct {
	Running   bool
	StartedAt string
	Tools     []tools.ToolState
}
