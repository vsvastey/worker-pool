package main

import (
	"io/ioutil"

	"github.com/Vastey/worker-pool/internal/task"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	WorkerCount int            `yaml:"worker_count"`
	Tasks       []*task.Config `yaml:"tasks"`
}

func getConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", filename)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal file %s content", filename)
	}
	return &config, nil
}
