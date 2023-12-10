package semaphore

// Semaphore Размещать интерфейс по месту использования?
type Semaphore interface {
	Acquire()
	Release()
	GetFreeRsrGlob() int
	GetFreeRsrLoc() int
}
