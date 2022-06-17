package worker

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Vastey/worker-pool/internal/task"
	"github.com/pkg/errors"
)

type SimpleWorker struct {
	ID          string
	taskConfigs <-chan *task.Config
	statusChan  chan Status
	stopChan    chan struct{}
	taskFactory task.Factory
}

func NewSimpleWorker(name string, taskFactory task.Factory, taskConfigs <-chan *task.Config) (*SimpleWorker, error) {
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
		ID:          name,
		taskConfigs: taskConfigs,
		statusChan:  statusChan,
		stopChan:    make(chan struct{}),
		taskFactory: taskFactory,
	}, nil
}

func (sw SimpleWorker) Status() <-chan Status {
	return sw.statusChan
}

func (sw *SimpleWorker) Stop() {
	sw.stopChan <- struct{}{}
	close(sw.statusChan)
}

func (sw *SimpleWorker) runTask(taskConfig *task.Config) error {
	t, err := sw.taskFactory.CreateTask(taskConfig)
	if err != nil {
		return errors.Wrap(err, "create task")
	}
	workerStatus := Status{
		ID:       sw.ID,
		Task:     t.Caption(),
		Progress: 0,
	}
	sw.statusChan <- workerStatus
	ch := t.Do()
	for taskStatus := range ch {
		workerStatus.Progress = taskStatus.Progress
		sw.statusChan <- workerStatus
	}
	workerStatus.Task = "idle"
	workerStatus.Progress = 0
	sw.statusChan <- workerStatus
	return nil
}

func (sw *SimpleWorker) Work(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case t := <-sw.taskConfigs:
			err := sw.runTask(t)
			if err != nil {
				log.Errorf("Error running task: %v", err)
			}
			if wg != nil {
				wg.Done()
			}
		case <-sw.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
