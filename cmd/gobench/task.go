package main

const (
	// HTTP Protocol
	HTTP = iota
	// WebSocket Protocol
	WebSocket
)

// Task struct
type Task struct {
	URL             string
	ProtocolVersion int
	HTTPMethod      string
}

// CreateTask create a task from given cmd args
func CreateTask(cmdArgs CmdArgs) (task *Task) {
	task = &Task{
		URL:             cmdArgs.URL,
		ProtocolVersion: cmdArgs.HTTPVersion,
		HTTPMethod:      cmdArgs.HTTPMethod,
	}
	return
}
