package channel

import (
	"context"
	"sync"
	"testing"
	"time"
)

func BenchmarkChannelSemaphore(b *testing.B) {
	var wg sync.WaitGroup

	n := int64(10)
	m := int64(4)

	globalSem := NewSemaphore(&n)
	for i := 0; i < b.N; i++ {

		taskSem := NewSemaphore(&m)
		for j := 0; j < b.N; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
				defer cancel()

				for {
					if err := globalSem.Acquire(ctx, 1); err != nil {
						continue
					}
					if err := taskSem.Acquire(ctx, 1); err != nil {
						globalSem.Release(1)
						continue
					}
					break
				}

				globalSem.Release(1)
				taskSem.Release(1)
			}()
		}
	}

	wg.Wait()
}
