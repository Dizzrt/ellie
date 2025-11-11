package consul

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	r := &Registry{
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
		timeout:           10 * time.Second,
		cli: &Client{
			dc:                             SingleDatacenter,
			cli:                            apiClient,
			resolver:                       defaultResolver,
			healthcheckInterval:            10,
			heartbeat:                      true,
			deregisterCriticalServiceAfter: 600,
			cancelers:                      make(map[string]*canceler),
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	r.lock.RLock()
	set := r.registry[serviceName]
	r.lock.RUnlock()

	getRemote := func() []*registry.ServiceInstance {
		services, _, err := r.cli.Service(ctx, serviceName, 0, true)
		if err == nil && len(services) > 0 {
			return services
		}

		return nil
	}

	if set == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}

		return nil, fmt.Errorf("service %s not resolved in registry", serviceName)
	}

	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if ss == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}

		return nil, fmt.Errorf("service %s not found in registry", serviceName)
	}

	return ss, nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.lock.Lock()
	set, ok := r.registry[serviceName]
	if !ok {
		cancelCtx, cancel := context.WithCancel(context.Background())
		set = &serviceSet{
			registry:    r,
			watcher:     make(map[*watcher]struct{}),
			services:    &atomic.Value{},
			serviceName: serviceName,
			ctx:         cancelCtx,
			cancel:      cancel,
		}

		r.registry[serviceName] = set
	}

	set.ref.Add(1)
	r.lock.Unlock()

	w := &watcher{
		event: make(chan struct{}, 1),
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{}
	set.lock.Unlock()

	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if len(ss) > 0 {
		select {
		case w.event <- struct{}{}:
		default:
		}
	}

	if !ok {
		if err := r.resolve(ctx, set); err != nil {
			return nil, err
		}
	}

	return w, nil
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

func (r *Registry) resolve(ctx context.Context, ss *serviceSet) error {
	listServices := r.cli.Service
	if r.timeout > 0 {
		listServices = func(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
			timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
			defer cancel()

			return r.cli.Service(timeoutCtx, service, index, passingOnly)
		}
	}

	services, idx, err := listServices(ctx, ss.serviceName, 0, true)
	if err != nil {
		return err
	}

	if len(services) > 0 {
		ss.broadcast(services)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				tempService, tempIndex, err := listServices(ss.ctx, ss.serviceName, idx, true)
				if err != nil {
					if err := sleepCtx(ss.ctx, time.Second); err != nil {
						return
					}

					continue
				}

				if len(tempService) != 0 && tempIndex != idx {
					services = tempService
					ss.broadcast(services)
				}

				idx = tempIndex

			case <-ss.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		var services []*registry.ServiceInstance
		ss, _ := set.services.Load().([]*registry.ServiceInstance)
		if ss == nil {
			continue
		}

		services = append(services, ss...)
		allServices[name] = services
	}

	return
}
