package semaphore

import (
	"sync"
	"sync/atomic"
)

type SemAtomic struct {
	max     int32
	FreeRsr int32
	mu      sync.Mutex
	cond    *sync.Cond
}

func NewSemAtomic(max int32) *SemAtomic {
	sem := &SemAtomic{
		max:     max,
		FreeRsr: max,
	}

	sem.cond = sync.NewCond(&sem.mu)

	return sem
}

func (s *SemAtomic) Acquire() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.FreeRsr == 0 {
		s.cond.Wait()
	}

	atomic.AddInt32(&s.FreeRsr, -1)
}

func (s *SemAtomic) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	atomic.AddInt32(&s.FreeRsr, 1)
	s.cond.Broadcast()
}

func (s *SemAtomic) GetFreeRsr() int {
	return int(s.FreeRsr)
}
