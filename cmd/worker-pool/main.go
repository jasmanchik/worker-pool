package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Semaphore interface {
}

type GlobalSem struct {
	maxTh        int32
	maxThPerT    int32
	maxThRunning int32
	wg           sync.WaitGroup
}

type TaskSem struct {
	running int32
}

func main() {
	var T, N, M int
	flag.IntVar(&T, "T", 3, "number of tasks")
	flag.IntVar(&N, "N", 10, "number of threads")
	flag.IntVar(&M, "M", 4, "max threads per task")
	flag.Parse()

	gSem := GlobalSem{
		int32(N),
		int32(M),
		0,
		sync.WaitGroup{},
	}

	for i := 0; i < T; i++ { //запускаем T задач
		tSem := TaskSem{}
		gSem.wg.Add(1)
		i := i
		go func(gSem *GlobalSem, tSem *TaskSem) {
			fmt.Printf("running go routine T %d \n", i)
			defer gSem.wg.Done()
			// каждая задача запускается в несколько потоков
			// если у задачи есть ресурсы для запуска и глобально есть свободные потоки, то запускаем
			for tSem.running < gSem.maxThPerT && gSem.maxTh > gSem.maxThRunning {
				atomic.AddInt32(&tSem.running, 1)
				atomic.AddInt32(&gSem.maxThRunning, 1)
				fmt.Printf("Start: local count %d, global count %d, task %d \n", tSem.running, gSem.maxThRunning, i)
				gSem.wg.Add(1)
				go func(gSem *GlobalSem, tSem *TaskSem) {
					defer gSem.wg.Done()
					defer atomic.AddInt32(&gSem.maxThRunning, -1)
					defer atomic.AddInt32(&tSem.running, -1)
					time.Sleep(time.Duration(rand.Intn(5)+3) * time.Second)
					fmt.Printf("End: local count %d, global count %d, task %d \n", tSem.running, gSem.maxThRunning, i)
				}(gSem, tSem)
			}
		}(&gSem, &tSem)
	}

	gSem.wg.Wait()
}
