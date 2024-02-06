package messages

type RunReply struct {
	Started bool
	Message string
}

type ToolStateReply struct {
	Name      string
	State     string
	LastRun   string
	LastError string
}

type StatusReply struct {
	Running   bool
	StartedAt string
	Tools     []ToolStateReply
}
