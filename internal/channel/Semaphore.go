package channel

type Semaphore struct {
	weight chan struct{}
}

func NewSemaphore(maxWeight *int64) *Semaphore {
	sem := &Semaphore{
		weight: make(chan struct{}, *maxWeight),
	}

	for i := 0; i < int(*maxWeight); i++ {
		sem.weight <- struct{}{}
	}

	return sem
}

func (s *Semaphore) Acquire() {
	<-s.weight
}

func (s *Semaphore) Release() {
	s.weight <- struct{}{}
}
