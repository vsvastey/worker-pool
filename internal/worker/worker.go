package worker

import (
	"sync"
)

type Status struct {
	ID       string
	Task     string
	Progress int
}

type Worker interface {
	Status() <-chan Status
	Work(wg *sync.WaitGroup)
	Stop()
}
