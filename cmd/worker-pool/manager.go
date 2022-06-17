package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/Vastey/worker-pool/internal/progressbar"
	"github.com/Vastey/worker-pool/internal/task"
	"github.com/Vastey/worker-pool/internal/taskconfigqueue"
	"github.com/Vastey/worker-pool/internal/util"
	"github.com/Vastey/worker-pool/internal/worker"
)

type WorkerAndProgress struct {
	worker worker.Worker
	pb     *progressbar.ProgressBar
}

type Manager struct {
	wps             []*WorkerAndProgress
	wg              *sync.WaitGroup
	taskConfigQueue *taskconfigqueue.TaskConfigQueue
	done            chan struct{}
	taskConfigChan  chan *task.Config
	taskFactory     task.Factory
}

func NewManager(taskFactory task.Factory) *Manager {
	return &Manager{
		wps:             []*WorkerAndProgress{},
		wg:              &sync.WaitGroup{},
		taskConfigQueue: taskconfigqueue.NewTaskConfigQueue(8),
		done:            make(chan struct{}),
		taskConfigChan:  make(chan *task.Config),
		taskFactory:     taskFactory,
	}
}

func (m *Manager) AddTask(taskConfig *task.Config) error {
	m.taskConfigQueue.Enqueue(taskConfig)
	return nil
}

func (m *Manager) AddWorker(ctx context.Context) error {
	w, err := worker.NewSimpleWorker(util.RandomString(5), m.taskFactory, m.taskConfigChan)
	if err != nil {
		return errors.Wrap(err, "simple worker constructor")
	}

	pb := progressbar.NewProgressBar()
	m.wps = append(m.wps, &WorkerAndProgress{
		worker: w,
		pb:     pb,
	})
	go w.Work(ctx, m.wg)
	return nil
}

func (m *Manager) Run() {
	taskConfig := m.taskConfigQueue.Dequeue()
	for taskConfig != nil {
		m.wg.Add(1)
		m.taskConfigChan <- taskConfig
		taskConfig = m.taskConfigQueue.Dequeue()
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
