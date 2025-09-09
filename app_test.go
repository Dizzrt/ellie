package ellie

import (
	"context"
	"testing"
	"time"

	"github.com/Dizzrt/ellie/internal/mock/ping"
	"github.com/Dizzrt/ellie/transport/grpc"
	"github.com/Dizzrt/ellie/transport/http"
)

type pingServer struct {
	ping.UnimplementedPingServiceServer
}

func (s *pingServer) Ping(ctx context.Context, req *ping.PingRequest) (*ping.PingResponse, error) {
	return &ping.PingResponse{
		Message: "pong",
	}, nil
}

func TestApp(t *testing.T) {
	gsrv := grpc.NewServer()
	hsrv := http.NewServer()

	ping.RegisterPingServiceServer(gsrv, &pingServer{})
	ping.RegisterPingHTTPServer(hsrv, &pingServer{})

	opts := []Option{
		Server(gsrv, hsrv),
	}

	app := New(opts...)
	time.AfterFunc(15*time.Second, func() {
		app.Stop()
	})

	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
