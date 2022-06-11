package worker

import "github.com/Vastey/worker-pool/internal/task"

type Worker interface {
	RunTask(task task.Task)
}
