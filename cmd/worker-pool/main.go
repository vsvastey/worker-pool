package main

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/worker"
	"sync"
	"time"
)

func main() {
	t1 := task.NewSleepTask(8*time.Second)
	t2 := task.NewSleepTask(10*time.Second)
	w1 := worker.NewSimpleWorker("worker1")
	w2 := worker.NewSimpleWorker("worker2")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		w1.RunTask(t1)
		wg.Done()
	}()
	go func() {
		w2.RunTask(t2)
		wg.Done()
	}()
	pb1 := progressbar.NewProgressBar()
	pb2 := progressbar.NewProgressBar()

	line1 := pb1.Draw()
	line2 := pb2.Draw()
	fmt.Println(line1)
	fmt.Println(line2)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case st1 := <-w1.Status():
				pb1.Set(fmt.Sprintf("%s - %s", st1.ID, st1.Task), st1.Progress)
			case st2 := <-w2.Status():
				pb2.Set(fmt.Sprintf("%s - %s", st2.ID, st2.Task), st2.Progress)
			case <-done:
				w1.Stop()
				w2.Stop()
				return
			}
			fmt.Print("\033[2F")
			line1 := pb1.Draw()
			line2 := pb2.Draw()
			fmt.Println(line1)
			fmt.Println(line2)
		}
	}()
	wg.Wait()
	done <- struct{}{}

}
