package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/Dizzrt/ellie/internal/endpoint"
	"github.com/Dizzrt/ellie/internal/host"
	"github.com/Dizzrt/ellie/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	*grpc.Server
	baseCtx  context.Context
	tlsConf  *tls.Config
	lis      net.Listener
	err      error
	network  string
	address  string
	endpoint *url.URL
	timeout  time.Duration
	// middleware
	// streamMiddleware
	unaryInts    []grpc.UnaryServerInterceptor
	streamInts   []grpc.StreamServerInterceptor
	grpcOpts     []grpc.ServerOption
	health       *health.Server
	customHealth bool
	// metadata
	cleanup           func()
	disableReflection bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		health:  health.NewServer(),
	}

	for _, opt := range opts {
		opt(srv)
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		// TODO
	}

	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	streamInts := []grpc.StreamServerInterceptor{
		// TODO
	}

	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}

	if srv.tlsConf != nil {
		temp := credentials.NewTLS(srv.tlsConf)
		grpcOpts = append(grpcOpts, grpc.Creds(temp))
	}

	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)
	// TODO metadata

	if !srv.customHealth {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	}

	if !srv.disableReflection {
		reflection.Register(srv.Server)
	}

	srv.cleanup, _ = admin.Register(srv.Server)
	return srv
}

func (s *Server) initializeListenerAndEndpoint() error {
	// initialize listener
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}

		s.lis = lis
	}

	// initialize endpoint
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}

		s.endpoint = endpoint.New(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}

	return s.err
}

// region interfaces impl

func (s *Server) Start(ctx context.Context) error {
	if err := s.initializeListenerAndEndpoint(); err != nil {
		return s.err
	}

	s.baseCtx = ctx
	// TODO log
	fmt.Printf("[gRPC] server listening on %s\n", s.lis.Addr().String())

	if !s.customHealth {
		s.health.Resume()
	}

	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.cleanup != nil {
		s.cleanup()
	}

	if !s.customHealth {
		s.health.Shutdown()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// TODO log
		fmt.Println("[gRPC] server stopping")
		s.GracefulStop()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		// TODO log
		fmt.Println("[gRPC] server couldn't stop gracefully in time, forcing stop")
		s.Server.Stop()
	}

	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.initializeListenerAndEndpoint(); err != nil {
		return nil, s.err
	}

	return s.endpoint, nil
}

// endregion
