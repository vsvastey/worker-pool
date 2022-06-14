package taskqueue

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
	t3 := task.NewMockTask(ctrl)
	t4 := task.NewMockTask(ctrl)

	tq := NewTaskQueue(2)
	assert.Equal(t, 2, len(tq.arr))
	// dequeue from empty queue
	res := tq.Dequeue()
	assert.Nil(t, res)

	tq.Enqueue(t1)
	res = tq.Dequeue()
	assert.Equal(t, t1, res)

	tq.Enqueue(t1)
	tq.Enqueue(t2)
	assert.Equal(t, 2, len(tq.arr))
	res = tq.Dequeue()
	assert.Equal(t, t1, res)
	assert.Equal(t, 2, len(tq.arr))

	tq.Enqueue(t3)
	assert.Equal(t, 2, len(tq.arr))
	tq.Enqueue(t4)
	assert.Equal(t, 4, len(tq.arr))
	res = tq.Dequeue()
	assert.Equal(t, t2, res)
	res = tq.Dequeue()
	assert.Equal(t, t3, res)
	res = tq.Dequeue()
	assert.Equal(t, t4, res)
	assert.Equal(t, 4, len(tq.arr))
}

func TestTaskQueueCapacityIsPowerOfTwo(t *testing.T) {
	cases := map[int]int{
		0:  2,
		2:  2,
		3:  4,
		4:  4,
		7:  8,
		8:  8,
		9:  16,
		15: 16,
	}
	for k, v := range cases {
		tq := NewTaskQueue(k)
		assert.Equal(t, v, len(tq.arr))
	}
}
