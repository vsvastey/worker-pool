package task

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/yaml_rawmessage"
)

const (
	sleep     = "sleep"
	copy_file = "copy_file"
	s3_upload = "s3_upload"
)

type Config struct {
	Type   string                         `yaml:"type"`
	Config yaml_rawmessage.YAMLRawMessage `yaml:"config"`
}

type Factory interface {
	CreateTask(config *Config) (Task, error)
}

type DefaultFactory struct{}

func (df DefaultFactory) CreateTask(config *Config) (Task, error) {
	switch config.Type {
	case sleep:
		var payload SleepTaskConfig
		config.Config.Unmarshal(&payload)
		return NewSleepTask(&payload)
	case copy_file:
		var payload CopyFileTaskConfig
		config.Config.Unmarshal(&payload)
		return NewCopyFileTask(&payload)
	case s3_upload:
		var payload S3UploadConfig
		config.Config.Unmarshal(&payload)
		return NewS3Upload(&payload)
	default:
		return nil, fmt.Errorf("unknown Task type %v", config.Type)
	}
}
