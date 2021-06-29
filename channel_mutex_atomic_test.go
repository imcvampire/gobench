package gobench

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// goos: linux
// goarch: amd64
// pkg: github.com/imcvampire/gobench
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkChannel-8       5230195               241.8 ns/op
// BenchmarkMutex-8         8927854               114.9 ns/op
// BenchmarkAtomic-8       14916924                81.40 ns/op
// PASS
// ok      github.com/imcvampire/gobench   3.958s

var sharedValue int64

var valueMutex sync.Mutex

func addAtomic(x int64) {
	atomic.AddInt64(&sharedValue, x)
}

func channelUpdater(ch <-chan int64) {
	for x := range ch {
		sharedValue += x
	}
}

func addWithLocking(x int64) {
	valueMutex.Lock()
	sharedValue += x
	valueMutex.Unlock()
}

func runBenchmark(b *testing.B, fn func()) {
	sharedValue = 0
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					fn()
				}
			}
		}()
	}

	for int(sharedValue) <= b.N {
		time.Sleep(100)
	}
	cancel()
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan int64)
	go channelUpdater(ch)

	var sendOnly chan<- int64 = ch
	runBenchmark(b, func() {
		sendOnly <- 1
	})

	close(ch)
}

func BenchmarkMutex(b *testing.B) {
	runBenchmark(b, func() {
		addWithLocking(1)
	})
}

func BenchmarkAtomic(b *testing.B) {
	runBenchmark(b, func() {
		addAtomic(1)
	})
}
