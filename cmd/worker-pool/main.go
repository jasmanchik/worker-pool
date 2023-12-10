package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
	"worker-pool/lib/counters/semaphore"
)

type TaskLimiter struct {
	counter int32
}

func main() {
	var T, N, M int
	flag.IntVar(&T, "T", 3, "number of tasks")
	flag.IntVar(&N, "N", 10, "number of threads")
	flag.IntVar(&M, "M", 4, "max threads per task")
	flag.Parse()

	var wg sync.WaitGroup
	//sem := semaphore.NewSemCh(int32(N))
	sem := semaphore.NewSemAtomic(int32(N))

	start := time.Now()

	for i := 0; i < T; i++ { //запускаем T задач
		i := i
		for j := 0; j < 100; j++ {
			wg.Add(1)
			j := j
			go func() {

				fmt.Printf("Пытаюсь захватить поток на %d итерации %d задачи, свободных потоков %d \n", j, i, sem.GetFreeRsr())
				sem.Acquire()
				fmt.Printf("Захватил поток на %d итерации %d задачи, свободных потоков %d \n", j, i, sem.GetFreeRsr())

				time.Sleep(500 * time.Millisecond)

				sem.Release()

				wg.Done()
			}()
		}

	}

	wg.Wait()

	finalTime := time.Since(start).Seconds()

	fmt.Printf("Прошло времени %f", finalTime)
}
