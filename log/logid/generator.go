package logid

import (
	"sync"
	"time"
)

var _ Generator = (*generator)(nil)

type generator struct {
	mu            sync.Mutex
	seq           int
	millTimestamp int64
}

func (gen *generator) Generate() LogID {
	gen.mu.Lock()

	now := time.Now()
	mill := now.UnixMilli()
	if mill == gen.millTimestamp {
		gen.seq = gen.seq + 1
	} else {
		gen.seq = 0
	}

	gen.millTimestamp = mill
	seq := gen.seq
	gen.mu.Unlock()

	return NewID128Bits(now, seq)
}

func NewGenerator() Generator {
	return &generator{}
}
