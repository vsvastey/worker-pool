package task

import (
	"fmt"

	"github.com/Vastey/worker-pool/internal/yamlrawmessage"
)

const (
	sleep     = "sleep"
	copy_file = "copy_file"
	s3_upload = "s3_upload"
)

// Config represents universal Task Config
// It consists of Task type and variable Config
// For different Task types the structure of Config might be different
type Config struct {
	Type   string                        `yaml:"type"`
	Config yamlrawmessage.YAMLRawMessage `yaml:"config"`
}

// Factory is a factory of Tasks that builds a new Task based on Task Config
type Factory interface {
	// CreateTask returns a new Task constructed according to a config provided
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
		var payload S3UploadTaskConfig
		config.Config.Unmarshal(&payload)
		return NewS3UploadTask(&payload)
	default:
		return nil, fmt.Errorf("unknown Task type %v", config.Type)
	}
}
