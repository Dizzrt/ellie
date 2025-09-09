package registry

import (
	"context"
	"fmt"
	"sort"
)

type Registrar interface {
	Register(ctx context.Context, svc *ServiceInstance) error
	Deregister(ctx context.Context, svc *ServiceInstance) error
}

type Watcher interface {
	Next() ([]*ServiceInstance, error)
	Stop() error
}

type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type ServiceInstance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
}

func (si *ServiceInstance) String() string {
	return fmt.Sprintf("%s-%s", si.Name, si.ID)
}

func (si *ServiceInstance) Equal(other any) bool {
	if si == nil && other == nil {
		return true
	}

	if si == nil || other == nil {
		return false
	}

	osi, ok := other.(*ServiceInstance)
	if !ok {
		return false
	}

	if len(si.Endpoints) != len(osi.Endpoints) {
		return false
	}

	sort.Strings(si.Endpoints)
	sort.Strings(osi.Endpoints)
	for i := 0; i < len(si.Endpoints); i++ {
		if si.Endpoints[i] != osi.Endpoints[i] {
			return false
		}
	}

	if len(si.Metadata) != len(osi.Metadata) {
		return false
	}

	for k, v := range si.Metadata {
		if v != osi.Metadata[k] {
			return false
		}
	}

	return si.ID == osi.ID && si.Name == osi.Name && si.Version == osi.Version
}
