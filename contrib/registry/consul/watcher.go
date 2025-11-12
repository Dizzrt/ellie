package consul

import (
	"context"
	"time"

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

	// svcs := make([]*registry.ServiceInstance, 0)
	// temp, ok := w.set.services.Load().([]*registry.ServiceInstance)
	// if ok {
	// 	svcs = append(svcs, temp...)
	// }

	// if len(svcs) > 0 {
	// 	return svcs, nil
	// }

	// avoid block on no event
	timeoutCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		return []*registry.ServiceInstance{}, nil
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
