package task

import (
	"fmt"
	"time"
)

type SleepTaskConfig struct {
	Duration time.Duration `yaml:"duration"`
}

type SleepTask struct {
	total    time.Duration
	progress int
}

func NewSleepTask(config *SleepTaskConfig) (*SleepTask, error) {
	return &SleepTask{
		total:    config.Duration,
		progress: 0,
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

func (st SleepTask) Caption() string {
	return fmt.Sprintf("sleep %v", st.total)
}
