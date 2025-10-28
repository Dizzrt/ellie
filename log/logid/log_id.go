package logid

import "sync"

type LogID interface {
	String() string
}

type Generator interface {
	Generate() LogID
}

type globalGenerator struct {
	mu  sync.Mutex
	gen Generator
}

var global globalGenerator

func init() {
	global = globalGenerator{
		gen: NewGenerator(),
	}
}

func SetLogIDGenerator(gen Generator) {
	global.mu.Lock()
	defer global.mu.Unlock()

	global.gen = gen
}

func Generate() LogID {
	return global.gen.Generate()
}
