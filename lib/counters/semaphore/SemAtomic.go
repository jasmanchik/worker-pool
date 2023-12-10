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

	// Я правильно понимаю, что нам тут мьютекс не нужен по сути.
	// Мы его используем только для корректной работы sync.Cond?
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.FreeRsr == 0 {
		s.cond.Wait() //зависает, но мьютекс освобождается
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
