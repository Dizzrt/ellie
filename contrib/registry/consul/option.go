package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
)

type Option func(*Registry)

func WithHealthCheck(enable bool) Option {
	return func(r *Registry) {
		r.enableHealthCheck = enable
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(r *Registry) {
		r.timeout = timeout
	}
}

func WithDatacenter(datacenter Datacenter) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.dc = datacenter
	}
}

func WithHeartbeat(enable bool) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.heartbeat = enable
	}
}

func WithServiceResolver(resolver ServiceResolver) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.resolver = resolver
	}
}

func WithHealthCheckInterval(interval int) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.healthcheckInterval = interval
	}
}

func WithDeregisterCriticalServiceAfter(deregisterCriticalServiceAfter int) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.deregisterCriticalServiceAfter = deregisterCriticalServiceAfter
	}
}

func WithServiceChecks(checks ...*api.AgentServiceCheck) Option {
	return func(r *Registry) {
		if r.cli == nil {
			return
		}

		r.cli.serviceChecks = checks
	}
}
