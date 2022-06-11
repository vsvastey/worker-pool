package worker

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestSimpleWorkerRunsTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	task1 := task.NewMockTask(ctrl)
	task1.EXPECT().Do().DoAndReturn(func() <-chan task.Status{
		ch := make(chan task.Status)
		go func() {
			ch <- task.Status{State: "done", Progress: 100}
			close(ch)
		}()
		return ch
	}).Times(1)

	w := NewSimpleWorker("test")
	w.RunTask(task1)
	// TODO: remove printing from worker
	fmt.Println("")
}