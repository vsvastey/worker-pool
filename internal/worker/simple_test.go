package worker

import (
	"context"
	"sync"
	"testing"

	"github.com/Vastey/worker-pool/internal/task"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockTaskFactory struct {
	taskToReturn task.Task
}

func (mtf MockTaskFactory) CreateTask(taskConfig *task.Config) (task.Task, error) {
	return mtf.taskToReturn, nil
}

func TestSimpleWorkerRunsTask(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	task1 := task.NewMockTask(ctrl)
	task1.EXPECT().Name().Return("task1").Times(1)
	task1.EXPECT().Do().DoAndReturn(func() <-chan task.Status {
		ch := make(chan task.Status)
		go func() {
			ch <- task.Status{Progress: 10}
			ch <- task.Status{Progress: 50}
			ch <- task.Status{Progress: 100}
			close(ch)
		}()
		return ch
	}).Times(1)

	taskFactory := MockTaskFactory{
		taskToReturn: task1,
	}
	taskConfigChan := make(chan *task.Config)

	wg := sync.WaitGroup{}
	w, err := NewSimpleWorker("test", taskFactory, taskConfigChan)
	assert.Nil(t, err)

	wg.Add(1)
	go w.Work(ctx, &wg)

	go func() {
		taskConfigChan <- &task.Config{}
	}()

	go func() {
		for range w.Status() {
		}
	}()

	wg.Wait()
	w.Stop()
}
