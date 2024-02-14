package messages

type RunReply struct {
	Started bool
	Message string
}

type ToolStateReply struct {
	Name            string
	State           string
	LatestRun       string
	LatestLogChange string
	LatestError     string
}

type StatusReply struct {
	Running   bool
	StartedAt string
	Tools     []ToolStateReply
}

type LogsReply struct {
	Tool            string
	LatestLogChange string
	LatestRun       string
	Logs            []string
}

type LogReply struct {
	Tool            string
	State           string
	LatestLogChange string
	LatestRun       string
	LatestError     string
	Log             string
	LogDate         string
}
