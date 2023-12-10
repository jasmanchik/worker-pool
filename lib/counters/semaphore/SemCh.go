package semaphore

type SemCh struct {
	FreeRsr chan struct{}
}

func NewSemCh(max int32) *SemCh {
	sem := &SemCh{
		FreeRsr: make(chan struct{}, max),
	}

	for i := 0; i < int(max); i++ {
		sem.FreeRsr <- struct{}{}
	}

	return sem
}

func (s *SemCh) Acquire() {
	<-s.FreeRsr
}

func (s *SemCh) Release() {
	s.FreeRsr <- struct{}{}
}

func (s *SemCh) GetFreeRsr() int {
	return len(s.FreeRsr)
}
