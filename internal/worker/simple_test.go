package worker

import (
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSimpleWorkerRunsTask(t *testing.T) {
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

	taskChan := make(chan task.Task)

	wg := sync.WaitGroup{}
	w, err := NewSimpleWorker("test", taskChan)
	assert.Nil(t, err)

	wg.Add(1)
	go w.Work(&wg)

	go func() {
		taskChan <- task1
	}()

	go func() {
		for range w.Status() {
		}
	}()

	wg.Wait()
	w.Stop()
}
