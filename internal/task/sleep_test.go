package task

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	name := task.Name()
	assert.Equal(t, "sleep 100ms", name)
}
