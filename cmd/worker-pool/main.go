package main

import (
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/worker"
	"time"
)

func main() {
	t := task.NewSleepTask(10*time.Second)
	w := worker.NewSimpleWorker("worker1")
	w.RunTask(t)
}
