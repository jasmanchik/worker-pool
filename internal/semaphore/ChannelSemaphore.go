package semaphore

type ChannelSemaphore struct {
	gMaxThreads chan struct{}
	lMaxThreads chan struct{}
}

func NewChannelSemaphore(gMaxThreads chan struct{}, lMaxThreads int32) *ChannelSemaphore {
	sem := &ChannelSemaphore{
		gMaxThreads: gMaxThreads,
		lMaxThreads: make(chan struct{}, lMaxThreads),
	}

	for i := 0; i < int(lMaxThreads); i++ {
		sem.lMaxThreads <- struct{}{}
	}

	return sem
}

func (s *ChannelSemaphore) Acquire() {
	<-s.gMaxThreads
	<-s.lMaxThreads
}

func (s *ChannelSemaphore) Release() {
	s.gMaxThreads <- struct{}{}
	s.lMaxThreads <- struct{}{}
}

func (s *ChannelSemaphore) GetFreeRsrGlob() int {
	return len(s.gMaxThreads)
}

func (s *ChannelSemaphore) GetFreeRsrLoc() int {
	return len(s.lMaxThreads)
}
