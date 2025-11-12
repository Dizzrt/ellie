package discovery

import (
	"time"
)

const NAME = "discovery"

type Option func(b *builder)

func WithTimeout(timeout time.Duration) Option {
	return func(b *builder) {
		b.timeout = timeout
	}
}

func WithInsecure(insecure bool) Option {
	return func(b *builder) {
		b.insecure = insecure
	}
}

func WithSubsetSize(size int) Option {
	return func(b *builder) {
		b.subsetSize = size
	}
}

func WithDebugLog(debugLog bool) Option {
	return func(b *builder) {
		b.debugLog = debugLog
	}
}
