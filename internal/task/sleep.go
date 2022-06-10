package task

import (
	"time"
)

type SleepTask struct {
	total    time.Duration
	progress int
}

func NewSleepTask(interval time.Duration) *SleepTask {
	return &SleepTask{total: interval, progress: 0}
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
				res <- Status{State: INPROGRESS_STATE, Progress: st.progress}
			case <-timeout:
				res <- Status{State: DONE_STATE, Progress: 100}
				close(res)
				return
			}
		}
	}()
	return res
}
