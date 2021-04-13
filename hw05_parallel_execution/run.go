package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		m = 0
	}
	var (
		err      error
		errCount int32
		wg       = new(sync.WaitGroup)
		taskCh   = make(chan Task, len(tasks))
	)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range taskCh {
				if atomic.LoadInt32(&errCount) >= int32(m) {
					return
				}

				if err := t(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)
	wg.Wait()
	if atomic.LoadInt32(&errCount) >= int32(m) {
		err = ErrErrorsLimitExceeded
	}
	return err
}
