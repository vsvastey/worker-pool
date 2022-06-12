package task

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestCreateSleepTask(t *testing.T) {
	df := DefaultFactory{}
	yamlCfg := `---
type: sleep
config:
  duration: 25s`
	var taskConfig Config
	err := yaml.Unmarshal([]byte(yamlCfg), &taskConfig)
	assert.Nil(t, err)

	task, err := df.CreateTask(&taskConfig)
	assert.Nil(t, err)
	assert.Equal(t, "sleep 25s", task.Name())
}

func TestCreateCopyFileTask(t *testing.T) {
	df := DefaultFactory{}
	yamlCfg := `---
type: copy_file
config:
  source: /home/path/to/file/file.dat
  destination: /another/path/new_name.file
`
	var taskConfig Config
	err := yaml.Unmarshal([]byte(yamlCfg), &taskConfig)
	assert.Nil(t, err)

	task, err := df.CreateTask(&taskConfig)
	assert.Nil(t, err)
	assert.Equal(t, "copy file.dat", task.Name())
}
