package main

import (
	"flag"
	"fmt"
	"github.com/Vastey/worker-pool/internal/dispatcher"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/task_queue"
	"github.com/Vastey/worker-pool/internal/worker"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var (
	tasksFilename string
)

func init() {
	flag.StringVar(&tasksFilename, "tasks", "", "file contains list of tasks")
}

func getTaskConfigsFromFile(filename string) ([]*task.Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("read file %s", filename))
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unmarshal file %s content", filename))
	}
	return config.Tasks, nil
}

func main() {
	flag.Parse()
	queue := task_queue.NewTaskQueue()

	taskConfigs, err := getTaskConfigsFromFile(tasksFilename)
	if err != nil {
		panic(err)
	}

	taskFactory := task.DefaultFactory{}
	for _, taskConfig := range taskConfigs {
		t, err := taskFactory.CreateTask(taskConfig)
		if err == nil {
			queue.Enqueue(t)
		}
	}

	w1 := worker.NewSimpleWorker("w1")
	w2 := worker.NewSimpleWorker("w2")

	workers := []worker.Worker{w1, w2}
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

	t := queue.Dequeue()
	for t != nil {
		wg.Add(1)
		tasksChan <- t
		t = queue.Dequeue()
	}
	wg.Wait()

	done <- struct{}{}

}
