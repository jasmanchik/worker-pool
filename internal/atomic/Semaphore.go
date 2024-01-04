package atomic

import (
	"sync"
	"sync/atomic"
)

type Semaphore struct {
	weight *int64
	cond   *sync.Cond
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for *s.weight == 0 {
		s.cond.Wait() //зависает, но мьютекс освобождается
	}

	atomic.AddInt64(s.weight, -1)
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	atomic.AddInt64(s.weight, 1)
	s.cond.Broadcast()
}

func NewSemaphore(max *int64) *Semaphore {
	sem := &Semaphore{
		weight: max,
		cond:   sync.NewCond(&sync.Mutex{}),
	}
	return sem
}
