package registry

import "context"

type Registrar interface {
	Register(ctx context.Context, svc *ServiceInstance) error
	Deregister(ctx context.Context, svc *ServiceInstance) error
}

type ServiceInstance struct {
	ID        string
	Name      string
	Version   string
	Metadata  map[string]string
	Endpoints []string
}
