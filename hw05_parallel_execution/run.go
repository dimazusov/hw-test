package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	doneCh := make(chan error, len(tasks))
	defer close(doneCh)

	taskCh := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	var isErr = false
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(n)
	defer wg.Wait()

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				mutex.Lock()
				if isErr {
					mutex.Unlock()
					break
				}
				mutex.Unlock()

				doneCh <- task()
			}
		}()
	}

	var totalErrors int
	var err error

	for i := 0; i < len(tasks); i++ {
		err = <-doneCh
		if err != nil {
			totalErrors++
		}

		if totalErrors == m {
			mutex.Lock()
			isErr = true
			mutex.Unlock()

			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
