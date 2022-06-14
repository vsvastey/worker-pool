package taskconfigqueue

import (
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskConfigQueue(t *testing.T) {
	t1 := &task.Config{Type: "t1"}
	t2 := &task.Config{Type: "t2"}
	t3 := &task.Config{Type: "t3"}
	t4 := &task.Config{Type: "t4"}

	tq := NewTaskConfigQueue(2)
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

func TestTaskConfigQueueCapacityIsPowerOfTwo(t *testing.T) {
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
		tq := NewTaskConfigQueue(k)
		assert.Equal(t, v, len(tq.arr))
	}
}
