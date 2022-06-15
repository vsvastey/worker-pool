package taskconfigqueue

import (
	"github.com/Vastey/worker-pool/internal/task"
	"sync"
)

type TaskConfigQueue struct {
	arr  []*task.Config
	mu   sync.Mutex
	head int
	tail int
	size int
}

func NewTaskConfigQueue(capacity int) *TaskConfigQueue {
	powTwo := 2
	for powTwo < capacity {
		powTwo *= 2
	}
	return &TaskConfigQueue{
		arr: make([]*task.Config, powTwo),
		mu:  sync.Mutex{},
	}
}

func (tq *TaskConfigQueue) Enqueue(taskConfig *task.Config) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.size == len(tq.arr) {
		oldLen := len(tq.arr)
		newArr := make([]*task.Config, tq.size*2)
		for i := 0; i < oldLen; i++ {
			newArr[i] = tq.arr[(tq.head+i)%oldLen]
		}
		tq.head = 0
		tq.tail = oldLen
		tq.arr = newArr
		tq.size = oldLen + 1
	}
	tq.arr[tq.tail] = taskConfig
	tq.tail = (tq.tail + 1) % len(tq.arr)
	tq.size += 1
}

func (tq *TaskConfigQueue) Dequeue() *task.Config {
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
