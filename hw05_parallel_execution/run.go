package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	taskCh := make(chan Task, len(tasks))
	doneCh := make(chan error, len(tasks))
	defer close(taskCh)
	defer close(doneCh)

	for _, task := range tasks {
		taskCh <- task
	}

	var isErr = false

	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if !ok {
					break
				}

				m.Lock()
				err := isErr
				m.Unlock()
				if err == true {
					break
				}
				doneCh <- task()
			}
		}()
	}

	var totalErrors int
	var err error

	for i := 0; i < len(tasks); i++ {
		select {
		case err = <-doneCh:
			if err != nil {
				totalErrors++
			}

			if totalErrors == M {
				defer wg.Wait()

				m.Lock()
				isErr = true
				m.Unlock()
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}
