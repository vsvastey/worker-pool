package main

import (
	"context"
	"flag"

	"github.com/Vastey/worker-pool/internal/task"
)

var (
	configFilename string
)

func init() {
	flag.StringVar(&configFilename, "config", "", "configuration filename")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	config, err := getConfigFromFile(configFilename)
	if err != nil {
		panic(err)
	}

	manager := NewManager(&task.DefaultFactory{})

	for i := 0; i < config.WorkerCount; i++ {
		err := manager.AddWorker(ctx)
		if err != nil {
			// TODO: log
		}
	}

	for _, taskConfig := range config.Tasks {
		err = manager.AddTask(taskConfig)
		if err != nil {
			//TODO: log
		}
	}

	go manager.ShowProgress()

	manager.Run()
}
