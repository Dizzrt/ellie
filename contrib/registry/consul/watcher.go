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
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if err := w.ctx.Err(); err != nil {
		return nil, err
	}

	select {
	case <-w.ctx.Done():
		return nil, nil
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
