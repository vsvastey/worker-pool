package worker

import "github.com/Vastey/worker-pool/internal/task"

type Status struct {
	ID string
	Task string
	Progress int
}

type Worker interface {
	RunTask(task task.Task)
	Status() <-chan Status
}
