package task

import (
	"fmt"
	"time"
)

// SleepTaskConfig is a configuration of SleepTask
type SleepTaskConfig struct {
	// Duration is a time.Duration interval that must be slept
	Duration time.Duration `yaml:"duration"`
}

// SleepTask is a task that does nothing (sleeps) provided amount of time
type SleepTask struct {
	// total is a time interval that must be slept
	total time.Duration
	// progress is how much percent of total has been already slept
	progress int
	// name is a task description
	name string
}

// NewSleepTask is a constructor of SleepTask
func NewSleepTask(config *SleepTaskConfig) (*SleepTask, error) {
	return &SleepTask{
		total:    config.Duration,
		progress: 0,
		name:     fmt.Sprintf("sleep %v", config.Duration),
	}, nil
}

func (st *SleepTask) Do() <-chan Status {
	res := make(chan Status)

	step := time.NewTicker(st.total / 100)
	timeout := time.After(st.total)
	go func() {
		for {
			select {
			case <-step.C:
				st.progress += 1
				res <- Status{Progress: st.progress}
			case <-timeout:
				res <- Status{Progress: 100}
				close(res)
				return
			}
		}
	}()
	return res
}

func (st SleepTask) Name() string {
	return st.name
}
