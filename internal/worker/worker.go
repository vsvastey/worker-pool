package worker

import (
	"context"
	"sync"
)

type Status struct {
	ID       string
	Task     string
	Progress int
}

type Worker interface {
	Status() <-chan Status
	Work(ctx context.Context, wg *sync.WaitGroup)
	Stop()
}
