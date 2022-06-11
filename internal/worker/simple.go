package worker

import (
	"github.com/Vastey/worker-pool/internal/task"
	"sync"
)

type SimpleWorker struct {
	ID string
	taskChan chan task.Task
	statusChan chan Status
	stopChan chan struct{}
}

func NewSimpleWorker(name string) *SimpleWorker {
	newStatus := Status{
		ID:       name,
		Task:     "idle",
		Progress: 0,
	}
	statusChan := make(chan Status)
	go func() {
		statusChan <- newStatus
	}()
	return &SimpleWorker{
		ID: name,
		taskChan: make(chan task.Task),
		statusChan: statusChan,
		stopChan: make(chan struct{}),
	}
}

func (sw SimpleWorker) Status() <-chan Status{
	return sw.statusChan
}

func (sw *SimpleWorker) Stop() {
	sw.stopChan <- struct{}{}
	close(sw.statusChan)
}

func (sw *SimpleWorker) runTask(task task.Task) {
	workerStatus := Status{
		ID: sw.ID,
		Task: task.Name(),
		Progress: 0,
	}
	sw.statusChan <- workerStatus
	ch := task.Do()
	for taskStatus := range ch{
		workerStatus.Progress = taskStatus.Progress
		sw.statusChan <- workerStatus
	}
	workerStatus.Task = "idle"
	workerStatus.Progress = 0
	sw.statusChan <- workerStatus
}

func (sw *SimpleWorker) Work(pool chan chan task.Task, wg *sync.WaitGroup) {
	pool <- sw.taskChan
	for {
		select {
		case t := <-sw.taskChan:
		    sw.runTask(t)
			pool <- sw.taskChan
			if wg != nil {
				wg.Done()
			}

		case <-sw.stopChan:
			return
		}
	}
}