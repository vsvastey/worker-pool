package dispatcher

import (
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/worker"
	"sync"
)

type Dispatcher struct {
	tasks chan task.Task
	workersInput chan chan task.Task
}

func NewDispatcher(workers []worker.Worker, tasks chan task.Task, wg *sync.WaitGroup) *Dispatcher{
	workersInput := make(chan chan task.Task, len(workers))
	for _, w := range workers {
		go func(worker worker.Worker) {
			worker.Work(workersInput, wg)
		}(w)
	}
	return &Dispatcher{
		tasks:   tasks,
		workersInput: workersInput,
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case t := <-d.tasks:
			go func(t task.Task) {
				ch := <-d.workersInput
				ch <- t
			}(t)
		}
	}
}
