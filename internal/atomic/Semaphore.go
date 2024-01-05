package atomic

import (
	"context"
	"sync"
	"sync/atomic"
)

type Semaphore struct {
	weight *int64
	cond   *sync.Cond
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for atomic.LoadInt64(s.weight) == 0 {
		waitCh := make(chan struct{})
		go func() {
			s.cond.Wait()
			close(waitCh)
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-waitCh:
		}
	}

	atomic.AddInt64(s.weight, -1)

	return nil
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	atomic.AddInt64(s.weight, 1)
	s.cond.Signal()
}

func (s *Semaphore) Weight() int64 {
	return atomic.LoadInt64(s.weight)
}

func NewSemaphore(max *int64) *Semaphore {
	sem := &Semaphore{
		weight: max,
		cond:   sync.NewCond(&sync.Mutex{}),
	}
	return sem
}
