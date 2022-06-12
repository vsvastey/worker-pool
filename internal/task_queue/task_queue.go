package task_queue

import (
	"github.com/Vastey/worker-pool/internal/task"
	"sync"
)

type TaskQueue struct {
	arr []task.Task
	mu sync.Mutex
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		arr: []task.Task{},
		mu:  sync.Mutex{},
	}
}

func (tq *TaskQueue) Enqueue(t task.Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.arr = append(tq.arr, t)
}

func (tq *TaskQueue) Dequeue() task.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.arr) > 0 {
		res := tq.arr[0]
		tq.arr = tq.arr[1:]
		return res
	}
	return nil
}