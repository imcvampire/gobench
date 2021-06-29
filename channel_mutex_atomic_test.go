package gobench

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/imcvampire/gobench
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkChannel-8       5964044               204.6 ns/op
// BenchmarkMutex-8        18504628                69.57 ns/op
// BenchmarkAtomic-8       51319978                21.21 ns/op

var sharedValue int64

var valueMutex sync.Mutex

func addAtomic(x int64) {
	atomic.AddInt64(&sharedValue, x)
}

func channelUpdater(b *testing.B, ch <-chan int64) {
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
	var n = int64(b.N)

	sharedValue = 0

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				if sharedValue >= n {
					wg.Done()
					return
				} else {
					fn()
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan int64)
	go channelUpdater(b, ch)

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
