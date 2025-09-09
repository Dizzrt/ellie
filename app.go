package ellie

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dizzrt/ellie/log"
	"github.com/Dizzrt/ellie/registry"
	"github.com/Dizzrt/ellie/transport"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoints() []string
}

type App struct {
	opts     options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 15 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.logger != nil {
		log.SetLogger(o.logger)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (app *App) ID() string {
	return app.opts.id
}

func (app *App) Name() string {
	return app.opts.name
}

func (app *App) Version() string {
	return app.opts.version
}

func (app *App) Metadata() map[string]string {
	return app.opts.metadata
}

func (app *App) Endpoints() []string {
	if app.instance != nil {
		return app.instance.Endpoints
	}

	return nil
}

func (app *App) Run() error {
	instance, err := app.buildInstance()
	if err != nil {
		return err
	}

	app.mu.Lock()
	app.instance = instance
	app.mu.Unlock()

	sctx := NewContext(app.ctx, app)
	for _, fn := range app.opts.beforeStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	wg := sync.WaitGroup{}
	octx := NewContext(app.opts.ctx, app)
	eg, ctx := errgroup.WithContext(sctx)

	for _, srv := range app.opts.servers {
		server := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx := context.WithoutCancel(octx)
			if app.opts.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(stopCtx, app.opts.stopTimeout)
				defer cancel()
			}

			return server.Stop(stopCtx)
		})

		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return server.Start(octx)
		})
	}

	wg.Wait()
	if app.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(ctx, app.opts.registrarTimeout)
		defer rcancel()

		if err = app.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}

	for _, fn := range app.opts.afterStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, app.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return app.Stop()
		}
	})

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	err = nil
	for _, fn := range app.opts.afterStop {
		err = fn(sctx)
	}

	return err
}

func (app *App) Stop() error {
	var err error = nil

	sctx := NewContext(app.ctx, app)
	for _, fn := range app.opts.beforeStop {
		err = fn(sctx)
	}

	app.mu.Lock()
	instance := app.instance
	app.mu.Unlock()

	if app.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(app.ctx, app), app.opts.registrarTimeout)
		defer cancel()

		if err = app.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}

	if app.cancel != nil {
		app.cancel()
	}

	return err
}

func (app *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(app.opts.endpoints))
	for _, e := range app.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	if len(endpoints) == 0 {
		for _, srv := range app.opts.servers {
			if temp, ok := srv.(transport.Endpointer); ok {
				e, err := temp.Endpoint()
				if err != nil {
					return nil, err
				}

				endpoints = append(endpoints, e.String())
			}
		}
	}

	return &registry.ServiceInstance{
		ID:        app.opts.id,
		Name:      app.opts.name,
		Version:   app.opts.version,
		Metadata:  app.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

type appKey struct{}

func NewContext(ctx context.Context, info AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, info)
}

func FromContext(ctx context.Context) (AppInfo, bool) {
	info, ok := ctx.Value(appKey{}).(AppInfo)
	return info, ok
}
