package main

import (
	"fmt"
	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/taskqueue"
	"github.com/Vastey/worker-pool/internal/util"
	"github.com/Vastey/worker-pool/internal/worker"
	"sync"
)

type WorkerAndProgress struct {
	worker worker.Worker
	pb     *progressbar.ProgressBar
}

type Manager struct {
	wps         []*WorkerAndProgress
	wg          *sync.WaitGroup
	taskQueue   *taskqueue.TaskQueue
	done        chan struct{}
	taskChan    chan task.Task
	taskFactory task.Factory
}

func NewManager(taskFactory task.Factory) *Manager {
	return &Manager{
		wps:         []*WorkerAndProgress{},
		wg:          &sync.WaitGroup{},
		taskQueue:   taskqueue.NewTaskQueue(8),
		done:        make(chan struct{}),
		taskChan:    make(chan task.Task),
		taskFactory: taskFactory,
	}
}

func (m *Manager) AddTask(taskConfig *task.Config) error {
	t, err := m.taskFactory.CreateTask(taskConfig)
	if err != nil {
		return err
	}
	m.taskQueue.Enqueue(t)
	return nil
}

func (m *Manager) AddWorker() error {
	w, err := worker.NewSimpleWorker(util.RandomString(5), m.taskChan)
	if err != nil {
		// TODO: log
		return err
	}

	pb := progressbar.NewProgressBar()
	m.wps = append(m.wps, &WorkerAndProgress{
		worker: w,
		pb:     pb,
	})
	go w.Work(m.wg)
	return nil
}

func (m *Manager) Run() {
	t := m.taskQueue.Dequeue()
	for t != nil {
		m.wg.Add(1)
		m.taskChan <- t
		t = m.taskQueue.Dequeue()
	}
	m.wg.Wait()
	m.done <- struct{}{}
}

func (m *Manager) ShowProgress() {
	for i := 0; i < len(m.wps); i++ {
		fmt.Println(m.wps[i].pb.Draw())
	}

	updated := make(chan struct{})
	for i := 0; i < len(m.wps); i++ {
		go func(wp *WorkerAndProgress) {
			for st := range wp.worker.Status() {
				wp.pb.Set(fmt.Sprintf("%s - %s", st.ID, st.Task), st.Progress)
				updated <- struct{}{}
			}
		}(m.wps[i])
	}

	for {
		select {
		case <-updated:
			fmt.Print("\033[", len(m.wps), "F")
			for i := 0; i < len(m.wps); i++ {
				fmt.Println(m.wps[i].pb.Draw())
			}
		case <-m.done:
			for i := 0; i < len(m.wps); i++ {
				m.wps[i].worker.Stop()
			}
			return
		}
	}
}
