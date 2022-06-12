package main

import (
	"flag"
	"fmt"
	"github.com/Vastey/worker-pool/internal/dispatcher"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/task_queue"
	"github.com/Vastey/worker-pool/internal/worker"
	"sync"
)

var (
	configFilename string
)

func init() {
	flag.StringVar(&configFilename, "config", "", "configuration")
}

type WorkerAndProgress struct {
	worker worker.Worker
	pb     *progressbar.ProgressBar
}

func main() {
	flag.Parse()

	config, err := getConfigFromFile(configFilename)
	if err != nil {
		panic(err)
	}

	queue := task_queue.NewTaskQueue()

	taskFactory := task.DefaultFactory{}
	for _, taskConfig := range config.Tasks {
		t, err := taskFactory.CreateTask(taskConfig)
		if err == nil {
			queue.Enqueue(t)
		}
	}

	wps := make([]*WorkerAndProgress, config.WorkerCount)
	workers := make([]worker.Worker, config.WorkerCount)
	for i := 0; i < config.WorkerCount; i++ {
		w, err := worker.NewSimpleWorker(fmt.Sprintf("w%d", i))
		if err != nil {
			// TODO: log
			continue
		}
		workers[i] = w
		pb := progressbar.NewProgressBar()
		wps[i] = &WorkerAndProgress{
			worker: w,
			pb:     pb,
		}
	}

	tasksChan := make(chan task.Task)
	wg := sync.WaitGroup{}
	d := dispatcher.NewDispatcher(workers, tasksChan, &wg)
	go d.Dispatch()

	done := make(chan struct{})

	go showProgress(wps, done)

	t := queue.Dequeue()
	for t != nil {
		wg.Add(1)
		tasksChan <- t
		t = queue.Dequeue()
	}
	wg.Wait()

	done <- struct{}{}
}
