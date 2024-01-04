package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
	"worker-pool/internal/channel"
)

type Semaphore interface {
	Acquire()
	Release()
}

func main() {
	var T, N, M int
	flag.IntVar(&T, "T", 3, "number of tasks")
	flag.IntVar(&N, "N", 10, "number of threads")
	flag.IntVar(&M, "M", 4, "max threads per task")
	flag.Parse()

	var wg sync.WaitGroup
	m := int32(M)
	n := int32(N)

	globalSem := channel.NewSemaphore(&n)
	for i := 0; i < T; i++ {
		taskSem := channel.NewSemaphore(&m)

		for j := 0; j < 100; j++ {
			j := j
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()

				globalSem.Acquire()
				defer globalSem.Release()
				taskSem.Acquire()
				defer taskSem.Release()

				time.Sleep(100 * time.Millisecond)
				fmt.Printf("task %d, thread %d\n", i, j)
			}()
		}
	}

	wg.Wait()
}
