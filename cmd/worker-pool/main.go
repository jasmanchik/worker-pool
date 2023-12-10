package main

import (
	"flag"
	"sync"
	"time"
	"worker-pool/internal/semaphore"
)

func main() {
	var T, N, M int
	flag.IntVar(&T, "T", 3, "number of tasks")
	flag.IntVar(&N, "N", 10, "number of threads")
	flag.IntVar(&M, "M", 4, "max threads per task")
	flag.Parse()

	var wg sync.WaitGroup

	//Для реализации через атомики
	gMaxTh := int32(N)

	//для реализации через каналы
	//gMaxTh := make(chan struct{}, N)
	//for i := 0; i < N; i++ {
	//	gMaxTh <- struct{}{}
	//}

	for i := 0; i < T; i++ {
		sem := semaphore.NewAtomicSemaphore(&gMaxTh, int32(M))
		//sem := semaphore.NewChannelSemaphore(gMaxTh, int32(M))

		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				sem.Acquire()
				defer sem.Release()

				time.Sleep(500 * time.Millisecond)
			}()
		}

	}

	wg.Wait()
}
