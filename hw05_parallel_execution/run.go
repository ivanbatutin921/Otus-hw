package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	var errorCount int32
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
			wg.Done()
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorCount) >= int32(m) {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if atomic.LoadInt32(&errorCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
