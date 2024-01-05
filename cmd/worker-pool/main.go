package main

import (
	"context"
	"fmt"
	"time"
)

type Semaphore interface {
	Acquire(ctx context.Context, weight int64)
	Release(weight int64)
}

func main() {

	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("read from closed channel")
			default:
				fmt.Println("default")
			}
		}
	}()

	close(ch)
	time.Sleep(time.Millisecond)

	//var T, N, M int
	//flag.IntVar(&T, "T", 3, "number of tasks")
	//flag.IntVar(&N, "N", 10, "number of threads")
	//flag.IntVar(&M, "M", 4, "max threads per task")
	//flag.Parse()
	//
	//var wg sync.WaitGroup
	//
	//n := int64(N)
	//m := int64(M)
	//globalSem := atomic.NewSemaphore(&n)
	//for i := 0; i < T; i++ {
	//	taskSem := atomic.NewSemaphore(&m)
	//	for j := 0; j < 100; j++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//
	//			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	//			defer cancel()
	//			for {
	//				if err := globalSem.Acquire(ctx); err != nil {
	//					continue
	//				}
	//
	//				if err := taskSem.Acquire(ctx); err != nil {
	//					globalSem.Release()
	//					continue
	//				}
	//
	//				break
	//			}
	//
	//			time.Sleep(50 * time.Millisecond)
	//			globalSem.Release()
	//			taskSem.Release()
	//		}()
	//	}
	//}
	//
	//wg.Wait()
}
