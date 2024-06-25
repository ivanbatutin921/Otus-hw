package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	errorCh := make(chan error)
	var errorCount int
	var wg sync.WaitGroup

	// запускаем n горутин для выполнения задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				err := task()
				if err!= nil {
					errorCh <- err
					return
				}
			}
		}()
	}

	// запускаем горутину для принятия ошибок и подсчета их количества
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range errorCh {
			_ = err
			errorCount++
			if errorCount >= m {
				close(errorCh)
				close(taskCh)
			}
		}
	}()

	// отправляем задачи в канал
	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}