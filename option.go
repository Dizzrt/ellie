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
