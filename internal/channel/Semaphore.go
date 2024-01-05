package channel

import "context"

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

func (s *Semaphore) Acquire(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.weight:
			return nil
		default:
		}
	}
}

func (s *Semaphore) Release() {
	s.weight <- struct{}{}
}
