package taskqueue

import (
	"github.com/Vastey/worker-pool/internal/task"
	"sync"
)

type TaskQueue struct {
	arr  []task.Task
	mu   sync.Mutex
	head int
	tail int
	size int
}

func NewTaskQueue(capacity int) *TaskQueue {
	powTwo := 2
	for powTwo < capacity {
		powTwo *= 2
	}
	return &TaskQueue{
		arr: make([]task.Task, powTwo),
		mu:  sync.Mutex{},
	}
}

func (tq *TaskQueue) Enqueue(t task.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.size == len(tq.arr) {
		oldLen := len(tq.arr)
		newArr := make([]task.Task, tq.size*2)
		for i := 0; i < oldLen; i++ {
			newArr[i] = tq.arr[(tq.head+i)%oldLen]
		}
		tq.head = 0
		tq.tail = oldLen
		tq.arr = newArr
		tq.size = oldLen + 1
	}
	tq.arr[tq.tail] = t
	tq.tail = (tq.tail + 1) % len(tq.arr)
	tq.size += 1
}

func (tq *TaskQueue) Dequeue() task.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.size <= 0 {
		return nil
	}

	item := tq.arr[tq.head]
	tq.size -= 1
	tq.head = (tq.head + 1) % len(tq.arr)
	return item
}
