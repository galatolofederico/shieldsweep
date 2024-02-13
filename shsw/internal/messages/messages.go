package messages

type RunReply struct {
	Started bool
	Message string
}

type ToolStateReply struct {
	Name          string
	State         string
	LastRun       string
	LastLogChange string
	LastError     string
}

type StatusReply struct {
	Running   bool
	StartedAt string
	Tools     []ToolStateReply
}

type LogsReply struct {
	Tool          string
	LastLogChange string
	LastRun       string
	Logs          []string
}

type LogReply struct {
	Tool          string
	State         string
	LastLogChange string
	LastRun       string
	Log           string
	LastError     string
}
