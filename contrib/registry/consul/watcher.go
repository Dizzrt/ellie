package consul

import (
	"context"

	"github.com/dizzrt/ellie/registry"
)

var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	ctx    context.Context
	cancel context.CancelFunc

	initialized bool
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if err := w.ctx.Err(); err != nil {
		return nil, err
	}

	if !w.initialized {
		svcs := make([]*registry.ServiceInstance, 0)
		if temp, ok := w.set.services.Load().([]*registry.ServiceInstance); ok {
			svcs = append(svcs, temp...)
		}

		w.initialized = true
		return svcs, nil
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.event:
	}

	svcs := make([]*registry.ServiceInstance, 0)
	temp, ok := w.set.services.Load().([]*registry.ServiceInstance)
	if ok {
		svcs = append(svcs, temp...)
	}

	return svcs, nil
}

func (w *watcher) Stop() error {
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
		w.set.delete(w)
	}

	return nil
}
