package consul

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/dizzrt/ellie/registry"
)

type serviceSet struct {
	serviceName string
	ref         atomic.Int32
	lock        sync.RWMutex
	registry    *Registry
	services    *atomic.Value
	watcher     map[*watcher]struct{}

	ctx    context.Context
	cancel context.CancelFunc
}

func (set *serviceSet) broadcast(svcs []*registry.ServiceInstance) {
	set.services.Store(svcs)
	set.lock.RLock()
	defer set.lock.RUnlock()
	for k := range set.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}

func (set *serviceSet) delete(w *watcher) {
	set.lock.Lock()
	delete(set.watcher, w)
	set.lock.Unlock()
	set.registry.tryDelete(set)
}
