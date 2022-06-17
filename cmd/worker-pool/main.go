package main

import (
	"context"
	"flag"

	"github.com/Vastey/worker-pool/internal/task"

	log "github.com/sirupsen/logrus"
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
		log.Fatalf("Error reading config file %s: %v", configFilename, err)
	}

	manager := NewManager(&task.DefaultFactory{})

	for i := 0; i < config.WorkerCount; i++ {
		err := manager.AddWorker(ctx)
		if err != nil {
			log.Errorf("Error adding a worker: %v", err)
		}
	}

	for _, taskConfig := range config.Tasks {
		err = manager.AddTask(taskConfig)
		if err != nil {
			log.Errorf("Error adding a task: %v", err)
		}
	}

	go manager.ShowProgress()

	manager.Run()
}
