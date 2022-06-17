package worker

import (
	"context"
	"sync"
)

// Status represents current status of running worker
type Status struct {
	// ID is a string identifier of a worker
	ID string
	// Task is a text description of current running Task
	Task string
	// Progress in a number of percents of current Task progress
	Progress int
}

// Worker is an interface of a Worker that can run a Task
type Worker interface {
	// Status returns a channel of Worker Status
	// that will be updated during the Worker work
	Status() <-chan Status
	// Work starts monitoring for new tasks
	Work(ctx context.Context, wg *sync.WaitGroup)
	// Stop stops monitoring for new tasks
	Stop()
}
