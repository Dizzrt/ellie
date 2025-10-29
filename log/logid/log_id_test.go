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
	wg.Add(10)

	mu := sync.Mutex{}
	ids := make([]string, 0, 1000)
	for range 10 {
		go func() {
			subIds := make([]string, 0, 1000)
			for range 1000 {
				id := Generate()
				subIds = append(subIds, id.String())
			}

			mu.Lock()
			ids = append(ids, subIds...)
			mu.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()
	slices.Sort(ids)
	fmt.Println(ids)
}

func TestGlobalGeneratorRepeatCheck(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	mu := sync.Mutex{}
	ids := make([]string, 0, 10000)
	for range 10 {
		go func() {
			subIds := make([]string, 0, 1000)
			for range 1000 {
				id := Generate()
				subIds = append(subIds, id.String())
			}

			mu.Lock()
			ids = append(ids, subIds...)
			mu.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	mp := make(map[string]int)
	for _, id := range ids {
		temp := id[:21]
		mp[temp]++
	}

	for k, v := range mp {
		if v > 1 {
			fmt.Println(k, v)
		}
	}

	fmt.Println(len(ids), len(mp))
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
