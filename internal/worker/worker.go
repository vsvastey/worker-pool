package worker

import (
	"github.com/Vastey/worker-pool/internal/task"
	"sync"
)

type Status struct {
	ID string
	Task string
	Progress int
}

type Worker interface {
	Status() <-chan Status
	Work(pool chan chan task.Task, wg *sync.WaitGroup)
	Stop()
}
