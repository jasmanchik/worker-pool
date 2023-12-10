package semaphore

type Semaphore interface {
	Acquire()
	Release()
	GetFreeRsr() int
}
