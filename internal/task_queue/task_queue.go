package task_queue

import (
	"fmt"
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

func (tq *TaskQueue) Enqueue(t task.Task) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.arr = append(tq.arr, t)
	return nil
}

func (tq *TaskQueue) Dequeue() (task.Task, error) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.arr) > 0 {
		res := tq.arr[0]
		tq.arr = tq.arr[1:]
		return res, nil
	}
	return nil, fmt.Errorf("empty queue")
}