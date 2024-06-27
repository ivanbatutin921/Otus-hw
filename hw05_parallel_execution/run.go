package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan Task)
	errorCh := make(chan error)
	var errorCount int
	var wg sync.WaitGroup

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				err := task()
				if err != nil {
					errorCh <- err
					//return nil
				}
			}
		}()
	}

	go func() {
		for err := range errorCh {
			errorCount++
			if errorCount >= m {
				close(errorCh)
				break
			}
			_ = err
		}
	}()

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
