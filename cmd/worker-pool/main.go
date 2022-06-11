package main

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/dispatcher"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/worker"
	"sync"
	"time"
)

func main() {
	t1 := task.NewSleepTask(2*time.Second)
	t2 := task.NewSleepTask(10*time.Second)
	t3 := task.NewSleepTask(5*time.Second)
	t4 := task.NewSleepTask(4*time.Second)
	w1 := worker.NewSimpleWorker("worker1")
	w2 := worker.NewSimpleWorker("worker2")

	workers := []worker.Worker{w1, w2}
	tasks := []task.Task{t1, t2, t3, t4}
	tasksChan := make(chan task.Task)

	wg := sync.WaitGroup{}
	d := dispatcher.NewDispatcher(workers, tasksChan, &wg)
	go d.Dispatch()

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



	for _, t := range tasks {
		wg.Add(1)
		tasksChan <- t
	}
	wg.Wait()

	done <- struct{}{}

}
