package ellie

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/Dizzrt/ellie/registry"
	"github.com/Dizzrt/ellie/transport"
)

type Option func(opts *options)

type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	// TODO logger
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server

	// hooks
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

func ID(id string) Option {
	return func(opts *options) {
		opts.id = id
	}
}

func Name(name string) Option {
	return func(opts *options) {
		opts.name = name
	}
}

func Version(version string) Option {
	return func(opts *options) {
		opts.version = version
	}
}

func Metadata(metadata map[string]string) Option {
	return func(opts *options) {
		opts.metadata = metadata
	}
}

func Endpoints(endpoints ...*url.URL) Option {
	return func(opts *options) {
		opts.endpoints = endpoints
	}
}

func Context(ctx context.Context) Option {
	return func(opts *options) {
		opts.ctx = ctx
	}
}

func Signal(sigs ...os.Signal) Option {
	return func(opts *options) {
		opts.sigs = sigs
	}
}

// TODO logger

func Registrar(r registry.Registrar) Option {
	return func(opts *options) {
		opts.registrar = r
	}
}

func RegistrarTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.registrarTimeout = timeout
	}
}

func StopTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.stopTimeout = timeout
	}
}

func Server(servers ...transport.Server) Option {
	return func(opts *options) {
		opts.servers = servers
	}
}

func BeforeStart(fn func(context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStart = append(opts.beforeStart, fn)
	}
}

func BeforeStop(fn func(context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStop = append(opts.beforeStop, fn)
	}
}

func AfterStart(fn func(context.Context) error) Option {
	return func(opts *options) {
		opts.afterStart = append(opts.afterStart, fn)
	}
}

func AfterStop(fn func(context.Context) error) Option {
	return func(opts *options) {
		opts.afterStop = append(opts.afterStop, fn)
	}
}
