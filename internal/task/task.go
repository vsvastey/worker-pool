package task

// Status of a task
type Status struct {
	Progress int
}

// Task describes a single job
//
// Method Do runs the task in a goroutine and returns a channel of Status
// that provides information about task's status
type Task interface {
	// Do runs a task and returns a Status channel
	Do() <-chan Status
	// Name return text description of a task
	Name() string
}
