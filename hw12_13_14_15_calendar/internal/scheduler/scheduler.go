package scheduler

import (
	"sync"
	"time"
)

type scheduler struct {
	tasks   []Task
	taskErr error
	mu      sync.Mutex
}

type Scheduler interface {
	AddTask(f func() error, t time.Duration)
	Run() error
}

func New() Scheduler {
	return &scheduler{
		mu: sync.Mutex{},
	}
}

func (m *scheduler) AddTask(f func() error, t time.Duration) {
	m.tasks = append(m.tasks, newTask(f, t))
}

func (m *scheduler) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(len(m.tasks))

	for _, t := range m.tasks {
		go func(task Task) {
			defer wg.Done()

			for {
				m.mu.Lock()
				if m.taskErr != nil {
					m.mu.Unlock()
					return
				}
				m.mu.Unlock()

				err := task.Run()
				if err != nil {
					m.mu.Lock()
					m.taskErr = err
					m.mu.Unlock()
					return
				}

				task.Sleep()
			}
		}(t)
	}

	wg.Wait()

	return m.taskErr
}
