package logid

import (
	"sync"
	"time"
)

var _ Generator = (*generator)(nil)

const _SEQUENCE_MASK = 0xFFF

type generator struct {
	mu            sync.Mutex
	seq           int
	millTimestamp int64
}

func (gen *generator) Generate() LogID {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	now := time.Now()
	mill := now.UnixMilli()
	if mill == gen.millTimestamp {
		gen.seq = (gen.seq + 1) & _SEQUENCE_MASK
	} else {
		gen.seq = 0
	}

	gen.millTimestamp = mill
	return NewID128Bits(now, gen.seq)
}

func NewGenerator() Generator {
	return &generator{}
}
