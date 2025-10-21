package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Dizzrt/ellie/internal/mock/ping"
	"github.com/Dizzrt/ellie/middleware/tracing"
	"google.golang.org/grpc"
)

type pingServer struct {
	ping.UnimplementedPingServiceServer
}

func (s *pingServer) Ping(ctx context.Context, req *ping.PingRequest) (*ping.PingResponse, error) {
	return &ping.PingResponse{
		Message: "pong",
	}, nil
}

type testKey struct{}

func getPingServer(t *testing.T, opts ...ServerOption) *Server {
	opts = append(opts,
		UnaryInterceptor(
			tracing.UnaryServerInterceptor(),
		),
		Options(grpc.InitialConnWindowSize(0)),
	)

	srv := NewServer(opts...)
	ping.RegisterPingServiceServer(srv, &pingServer{})
	if e, err := srv.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
		t.Fatal(e, err)
	}

	return srv
}

func TestPing(t *testing.T) {
	// start server
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")

	srv := getPingServer(t)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	// client
	e, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}

	conn, err := DialInsecure(
		WithEndpoint(e.Host),
		// WithOptions(grpc.WithBlock()),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)

	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		t.Fatal(err)
	}

	client := ping.NewPingServiceClient(conn)
	resp, err := client.Ping(ctx, &ping.PingRequest{})
	if err != nil {
		t.Errorf("failed to call with error: %v", err)
	}

	t.Log(resp)

	_ = srv.Stop(ctx)
}

func TestPingWithTLS(t *testing.T) {
	// start server
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")

	srvPemF := "../../internal/mock/certs/server.pem"
	srvKeyF := "../../internal/mock/certs/server.key"
	cert, err := tls.LoadX509KeyPair(srvPemF, srvKeyF)
	if err != nil {
		t.Fatal(err)
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	srv := getPingServer(t, TLSConfig(tlsConf))
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	// client
	caPemF := "../../internal/mock/certs/ca.pem"
	caCert, err := os.ReadFile(caPemF)
	if err != nil {
		t.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		t.Fatal("unable to append CA certs")
	}

	clientTlsConf := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: "localhost",
		MinVersion: tls.VersionTLS12,
	}

	e, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}

	conn, err := Dial(
		WithEndpoint(e.Host),
		// WithOptions(grpc.WithBlock()),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
		WithTLSConfig(clientTlsConf),
	)

	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		t.Fatal(err)
	}

	client := ping.NewPingServiceClient(conn)
	resp, err := client.Ping(ctx, &ping.PingRequest{})
	if err != nil {
		t.Errorf("failed to call with error: %v", err)
	}

	t.Log(resp)

	_ = srv.Stop(ctx)
}

func TestPingWithTracing(t *testing.T) {
	ctx := context.Background()

	// init tracing provider
	tp, err := tracing.Initialize(ctx,
		tracing.ServiceName("transport-test"),
		tracing.ServiceVersion("v0.0.1-dev"),
		tracing.Endpoint("localhost:4317"),
		tracing.EndpointType(tracing.EndpointType_GRPC),
		tracing.Insecure(true),
		tracing.Metadata(map[string]string{
			"ip":  "127.0.0.1",
			"env": "dev",
		}),
	)
	if err != nil {
		log.Fatalf("init tracing provider failed: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}()

	// start server
	srv := getPingServer(t)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	// client
	e, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}

	conn, err := DialInsecure(
		WithEndpoint(e.Host),
		WithUnaryClientInterceptor(
			tracing.UnaryClientInterceptor(),
		),
	)

	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		t.Fatal(err)
	}

	client := ping.NewPingServiceClient(conn)
	resp, err := client.Ping(ctx, &ping.PingRequest{})
	if err != nil {
		t.Errorf("failed to call with error: %v", err)
	}

	t.Log(resp)
	_ = srv.Stop(ctx)
}
