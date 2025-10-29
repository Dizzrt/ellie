package logid

import (
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	gen := NewGenerator()

	wg := sync.WaitGroup{}
	wg.Add(100)
	for range 100 {
		go func() {
			id := gen.Generate()
			fmt.Println(id.String())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGlobalGenerator(t *testing.T) {
	wg := sync.WaitGroup{}

	mu := sync.Mutex{}
	ids := make([]string, 0, 10000)

	wg.Add(10000)
	for range 10000 {
		go func() {
			id := Generate()
			mu.Lock()
			ids = append(ids, id.String())
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	slices.Sort(ids)
	fmt.Println(ids)
}

func TestID128Bits(t *testing.T) {
	fmt.Println(NewID128Bits(time.Now(), 0))
}

func BenchmarkGlobalGeneratorConcurrent(b *testing.B) {
	const concurrency = 100
	perGoroutine := b.N / concurrency

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for range concurrency {
		go func() {
			defer wg.Done()
			for range perGoroutine {
				Generate()
			}
		}()
	}

	wg.Wait()
}
