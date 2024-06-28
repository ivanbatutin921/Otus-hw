package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("if all tasks are successful and take varying amounts of time, then all tasks are executed in the correct order", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
	
		var runTasksCount int32
		var taskResults []int
	
		for i := 0; i < tasksCount; i++ {
			taskNum := i
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				taskResults = append(taskResults, taskNum)
				return nil
			})
		}
	
		workersCount := 10
		maxErrorsCount := 50
		err := Run(tasks, workersCount, maxErrorsCount)
	
		require.Nil(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were executed")
		require.Equal(t, tasksCount, len(taskResults), "task results slice has incorrect length")
		for i := 0; i < tasksCount; i++ {
			require.Equal(t, i, taskResults[i], "tasks were not executed in the correct order")
		}
	})

	t.Run("if there are more errors than the error limit, then return an error", func(t *testing.T) {
		tasks := []Task{
			func() error { return nil }, // task 1: success
			func() error { return errors.New("error 1") }, // task 2: error
			func() error { return nil }, // task 3: success
			func() error { return errors.New("error 2") }, // task 4: error
			func() error { return errors.New("error 3") }, // task 5: error
		}

		n := 2 // 2 worker goroutines
		m := 2 // error limit: 2 errors allowed

		err := Run(tasks, n, m)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual error %q", err)
	})
}
