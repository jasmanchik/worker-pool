package internal

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"testing"
)

func BenchmarkXSemaphore(b *testing.B) {
	var wg sync.WaitGroup

	globalSem := semaphore.NewWeighted(int64(10))
	for i := 0; i < b.N; i++ {
		taskSem := semaphore.NewWeighted(int64(4))

		for j := 0; j < b.N; j++ {
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
			}()
		}
	}

	wg.Wait()
}
