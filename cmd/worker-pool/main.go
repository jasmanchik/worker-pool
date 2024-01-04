package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
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

	globalSem := semaphore.NewWeighted(int64(M))
	for i := 0; i < T; i++ {
		taskSem := semaphore.NewWeighted(int64(N))

		for j := 0; j < 100; j++ {
			j := j
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := globalSem.Acquire(context.Background(), 1)
				if err != nil {
					fmt.Errorf("cant't acquire semaphore: %v", err)
					return
				}
				defer globalSem.Release(1)
				err = taskSem.Acquire(context.Background(), 1)
				if err != nil {
					fmt.Errorf("cant't acquire semaphore: %v", err)
					return
				}
				defer taskSem.Release(1)

				time.Sleep(100 * time.Millisecond)
				fmt.Printf("task %d, thread %d\n", i, j)
			}()
		}
	}

	wg.Wait()
}
