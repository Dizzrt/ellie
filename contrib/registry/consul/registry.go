package consul

import (
	"context"
	"sync"
	"time"

	"github.com/dizzrt/ellie/registry"
	"github.com/hashicorp/consul/api"
)

var _ registry.Registrar = (*Registry)(nil)
var _ registry.Discovery = (*Registry)(nil)

type Registry struct {
	lock sync.RWMutex

	cli               *Client
	timeout           time.Duration
	registry          map[string]*serviceSet
	enableHealthCheck bool
}

func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{}

	return r
}

func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return nil, nil
}

func (r *Registry) tryDelete(ss *serviceSet) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	if ss.ref.Add(-1) != 0 {
		return false
	}
	ss.cancel()
	delete(r.registry, ss.serviceName)
	return true
}
