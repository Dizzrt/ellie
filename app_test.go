package ellie

import (
	"context"
	"testing"

	"github.com/dizzrt/ellie/contrib/registry/consul"
	"github.com/dizzrt/ellie/internal/mock/ping"
	"github.com/dizzrt/ellie/transport/grpc"
	"github.com/dizzrt/ellie/transport/http"
	"github.com/hashicorp/consul/api"
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
	gsrv := grpc.NewServer(grpc.Address(":50051"))
	hsrv := http.NewServer(http.Address(":8081"))
	ping.RegisterPingServiceServer(gsrv, &pingServer{})
	ping.RegisterPingServiceHTTPServer(hsrv, &pingServer{})

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	reg := consul.New(client)
	app := New(
		ID("test-app"),
		Name("test-app"),
		Version("dev"),
		Server(gsrv, hsrv),
		Registrar(reg),
	)

	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
