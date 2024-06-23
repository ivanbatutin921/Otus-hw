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
	resultCh := make(chan error)
	errCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go worker(taskCh, resultCh, &wg)
	}

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	go func() {
		errCount := 0
		for err := range resultCh {
			if err!= nil {
				errCount++
				if errCount >= m {
					errCh <- struct{}{}
					return
				}
			}
		}
		wg.Done()
	}()

	select {
	case <-errCh:
		return ErrErrorsLimitExceeded
	case <-wait(&wg):
		return nil
	}
}

func worker(taskCh <-chan Task, resultCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskCh {
		resultCh <- task()
	}
}

func wait(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}