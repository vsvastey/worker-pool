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

	pool := make(chan chan task.Task)

	w, err := NewSimpleWorker("test")
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go w.Work(pool, &wg)

	go func() {
		ch := <-pool
		ch <- task1
		ch = <-pool
	}()

	go func() {
		for range w.Status() {
		}
	}()

	wg.Wait()
	w.Stop()
}
