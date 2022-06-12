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
	res := tq.Dequeue()
	assert.Nil(t, res)

	tq.Enqueue(t1)
	res = tq.Dequeue()
	assert.Equal(t, t1, res)

	tq.Enqueue(t1)
	tq.Enqueue(t2)

	res = tq.Dequeue()
	assert.Equal(t, t1, res)

	res = tq.Dequeue()
	assert.Equal(t, t2, res)
}
