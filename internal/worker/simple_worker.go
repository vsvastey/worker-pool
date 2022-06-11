package worker

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
)

type SimpleWorker struct {
	pb *progressbar.ProgressBar
}

func NewSimpleWorker(name string) *SimpleWorker {
	return &SimpleWorker{
		pb: progressbar.NewProgressBar(name),
	}
}

func (sw *SimpleWorker) RunTask(task task.Task) {
	ch := task.Do()
	for status := range ch{
		sw.pb.Set(status.Progress)
		fmt.Printf("\r%s", sw.pb.Draw())
	}
}
