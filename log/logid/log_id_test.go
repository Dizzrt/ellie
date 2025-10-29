package logid

import (
	"fmt"
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

	wg.Add(200)
	for range 200 {
		go func() {
			id := Generate()
			fmt.Println(id.String())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestID128Bits(t *testing.T) {
	fmt.Println(NewID128Bits(time.Now(), 0))
}
