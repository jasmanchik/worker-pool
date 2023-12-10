package semaphore

import (
	"sync"
	"sync/atomic"
)

type AtomicSemaphore struct {
	gMaxThreads *int32
	lMaxThreads int32
	mu          sync.Mutex
	cond        *sync.Cond
}

func (s *AtomicSemaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for *s.gMaxThreads == 0 || s.lMaxThreads == 0 {
		s.cond.Wait() //зависает, но мьютекс освобождается
	}

	atomic.AddInt32(s.gMaxThreads, -1)
	atomic.AddInt32(&s.lMaxThreads, -1)
}

func (s *AtomicSemaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	atomic.AddInt32(s.gMaxThreads, 1)
	atomic.AddInt32(&s.lMaxThreads, 1)
	s.cond.Broadcast()
}

func NewAtomicSemaphore(gMaxThreads *int32, lMaxThreads int32) *AtomicSemaphore {
	sem := &AtomicSemaphore{
		gMaxThreads: gMaxThreads,
		lMaxThreads: lMaxThreads,
		mu:          sync.Mutex{},
	}
	sem.cond = sync.NewCond(&sem.mu)
	return sem
}

func (s *AtomicSemaphore) GetFreeRsrGlob() int {
	return int(*s.gMaxThreads)
}

func (s *AtomicSemaphore) GetFreeRsrLoc() int {
	return int(s.lMaxThreads)
}
