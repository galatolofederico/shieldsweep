package tools

func GetToolRunner(toolName string) ToolRunner {
	switch toolName {
	case "dummy1":
		return NewDummyTool(true, 1)
	case "dummy2":
		return NewDummyTool(true, 2)
	case "dummy3":
		return NewDummyTool(false, 0)
	default:
		return nil
	}
}
