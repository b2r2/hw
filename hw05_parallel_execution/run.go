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
	var err error
	var errCount rune
	taskCh := make(chan struct{}, n)
	wg := sync.WaitGroup{}
	defer wg.Wait()
	if m <= 0 {
		m = 0
	}
	for _, t := range tasks {
		if atomic.LoadInt32(&errCount) >= rune(m) {
			err = ErrErrorsLimitExceeded
			break
		}
		t := t
		taskCh <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-taskCh
				wg.Done()
			}()
			if atomic.LoadInt32(&errCount) >= rune(m) {
				return
			}
			if err := t(); err != nil {
				atomic.AddInt32(&errCount, 1)
			}
		}()
	}
	return err
}
