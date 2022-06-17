package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSleepTaskWorksAtLeastAsLongAsInterval(t *testing.T) {
	interval := 100 * time.Millisecond
	config := SleepTaskConfig{Duration: interval}
	task, err := NewSleepTask(&config)
	assert.Nil(t, err)

	start := time.Now()
	ch := task.Do()

	// wait until the chan is closed
	for range ch {
	}

	duration := time.Since(start)
	assert.True(t, duration >= interval)
}

func TestSleepTaskUpdatesProgress(t *testing.T) {
	interval := 100 * time.Millisecond
	config := SleepTaskConfig{Duration: interval}
	task, err := NewSleepTask(&config)
	assert.Nil(t, err)

	ch := task.Do()

	prev := 0
	var status Status
	for status = range ch {
		assert.True(t, status.Progress >= prev)
		prev = status.Progress
	}
	assert.Equal(t, 100, status.Progress)
}

func TestSleepTaskHasName(t *testing.T) {
	interval := 100 * time.Millisecond
	config := SleepTaskConfig{Duration: interval}
	task, err := NewSleepTask(&config)
	assert.Nil(t, err)

	name := task.Caption()
	assert.Equal(t, "sleep 100ms", name)
}
