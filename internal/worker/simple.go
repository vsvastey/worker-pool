package worker

import (
	"github.com/Vastey/worker-pool/internal/task"
)

type SimpleWorker struct {
	ID string
	statusChan chan Status
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
		statusChan: statusChan,
	}
}

func (sw SimpleWorker) Status() <-chan Status{
	return sw.statusChan
}

func (sw SimpleWorker) Stop() {
	close(sw.statusChan)
}

func (sw *SimpleWorker) RunTask(task task.Task) {
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
	workerStatus.ID = "idle"
	workerStatus.Progress = 0
	sw.statusChan <- workerStatus
}
