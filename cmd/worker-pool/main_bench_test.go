package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"testing"
	"time"
	"worker-pool/internal/atomic"
	"worker-pool/internal/channel"
)

func BenchmarkXSemaphore(b *testing.B) {
	for k := 0; k < b.N; k++ {
		var wg sync.WaitGroup

		globalSem := semaphore.NewWeighted(int64(10))
		for i := 0; i < 10; i++ {
			taskSem := semaphore.NewWeighted(int64(4))

			for j := 0; j < 100; j++ {
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
}

func BenchmarkChannelSemaphore(b *testing.B) {
	for k := 0; k < b.N; k++ {
		var wg sync.WaitGroup

		n := int64(10)
		m := int64(4)

		globalSem := channel.NewSemaphore(&n)
		for i := 0; i < 10; i++ {

			taskSem := channel.NewSemaphore(&m)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					for {
						if err := globalSem.Acquire(ctx); err != nil {
							continue
						}
						if err := taskSem.Acquire(ctx); err != nil {
							globalSem.Release()
							continue
						}
						break
					}

					globalSem.Release()
					taskSem.Release()
				}()
			}
		}

		wg.Wait()
	}
}

func BenchmarkAtomicSemaphore(b *testing.B) {
	for k := 0; k < b.N; k++ {
		var wg sync.WaitGroup

		n := int64(10)
		m := int64(4)

		globalSem := atomic.NewSemaphore(&n)
		for i := 0; i < 10; i++ {
			taskSem := atomic.NewSemaphore(&m)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for {
						ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
						if err := globalSem.Acquire(ctx); err != nil {
							cancel()
							continue
						}

						if err := taskSem.Acquire(ctx); err != nil {
							cancel()
							globalSem.Release()
							continue
						}
						cancel()
						break
					}

					globalSem.Release()
					taskSem.Release()
				}()
			}
		}

		wg.Wait()
	}
}
