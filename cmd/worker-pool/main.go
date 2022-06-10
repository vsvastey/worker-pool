package main

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"time"
)

func main() {
	pb := progressbar.NewProgressBar("task1:")
	t := task.NewSleepTask(10*time.Second)
	ch := t.Do()
	for status := range ch{
		pb.Set(status.Progress)
		fmt.Printf("\r%s", pb.Draw())
	}
}
