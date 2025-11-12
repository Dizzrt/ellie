package discovery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dizzrt/ellie/registry"
	"github.com/google/uuid"
	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*builder)(nil)

type builder struct {
	discoverer registry.Discovery
	timeout    time.Duration
	subsetSize int
	insecure   bool
	debugLog   bool
}

func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		discoverer: d,
		timeout:    10 * time.Second,
		insecure:   false,
		debugLog:   true,
		subsetSize: 25,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	watchRes := &struct {
		err error
		w   registry.Watcher
	}{}

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		w, err := b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		watchRes.w = w
		watchRes.err = err
		close(done)
	}()

	var err error
	if b.timeout > 0 {
		select {
		case <-done:
			err = watchRes.err
		case <-time.After(b.timeout):
			err = fmt.Errorf("discovery create watcher overtime")
		}
	} else {
		<-done
		err = watchRes.err
	}

	if err != nil {
		cancel()
		return nil, err
	}

	r := &discoveryResolver{
		w:           watchRes.w,
		cc:          cc,
		ctx:         ctx,
		cancel:      cancel,
		insecure:    b.insecure,
		debugLog:    b.debugLog,
		subsetSize:  b.subsetSize,
		selectorKey: uuid.New().String(),
	}

	go r.watch()
	return r, nil
}

func (_ *builder) Scheme() string {
	return NAME
}
