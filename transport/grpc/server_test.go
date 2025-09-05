package grpc

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Dizzrt/ellie/internal/mock/ping"
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

func getPingServer(t *testing.T) *Server {
	srv := NewServer(
		UnaryInterceptor(func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return handler(ctx, req)
		}),
		Options(grpc.InitialConnWindowSize(0)),
	)

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
