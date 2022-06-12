package main

import "github.com/Vastey/worker-pool/internal/task"

type Config struct {
	Tasks []*task.Config `yaml:"tasks"`
}
