package task_queue

import (
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t1 := task.NewMockTask(ctrl)
	t2 := task.NewMockTask(ctrl)

	tq := NewTaskQueue()
	// dequeue from empty queue
	_, err := tq.Dequeue()
	assert.NotNil(t, err)

	err = tq.Enqueue(t1)
	assert.Nil(t, err)
	tmp, err := tq.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, t1, tmp)

	err = tq.Enqueue(t1)
	assert.Nil(t, err)
	err = tq.Enqueue(t2)
	assert.Nil(t, err)

	tmp, err = tq.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, t1, tmp)

	tmp, err = tq.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, t2, tmp)
}
