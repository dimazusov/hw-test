package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	taskCh := make(chan Task, len(tasks))
	doneCh := make(chan error, len(tasks))
	defer close(taskCh)
	defer close(doneCh)

	for _, task := range tasks {
		taskCh <- task
	}

	var isErr = false

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if !ok {
					break
				}

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
			defer wg.Wait()

			mutex.Lock()
			isErr = true
			mutex.Unlock()
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
