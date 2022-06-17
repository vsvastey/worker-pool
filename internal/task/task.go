package task

// Status represent a current status of a running Task
type Status struct {
	// Progress in a number of percents completed
	Progress int
}

// Task is an interface of a Task that might be executed
type Task interface {
	// Do starts a task in a new goroutine and returns a channel of Status
	// that will be updated during the task's execution
	Do() <-chan Status
	// Caption returns a text description of a task
	Caption() string
}
