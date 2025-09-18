package http_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	nhttp "net/http"

	"github.com/Dizzrt/ellie/internal/mock/ping"
	"github.com/Dizzrt/ellie/transport/http"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pingServer struct {
	ping.UnimplementedPingServiceServer
}

func (s *pingServer) Ping(ctx context.Context, req *ping.PingRequest) (*ping.PingResponse, error) {
	status.New(codes.Unknown, "unknown error")
	return &ping.PingResponse{
		Message: "pong",
	}, nil
}

func (s *pingServer) Hello(ctx context.Context, req *ping.HelloRequest) (*ping.HelloResponse, error) {
	return &ping.HelloResponse{
		Message: fmt.Sprintf("hello %s, type is %s", req.GetName(), req.GetType()),
	}, nil
}

func TestHTTPServer(t *testing.T) {
	ctx := context.Background()

	var opts = []http.ServerOption{}
	srv := http.NewServer(opts...)

	ping.RegisterPingHTTPServer(srv, &pingServer{})
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	//
	e, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}

	url := e.String() + "/hello/ellie?type=mock"
	// resp, err := nhttp.Post(url, "application/json", strings.NewReader(``))
	resp, err := nhttp.Post(url, "application/json", strings.NewReader(`{"name": "ellieFromBody","type": "mockFromBody"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// if resp.statusCode != http.StatusOK {
	// 	t.Fatal(resp.statusCode)
	// }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(body))
	_ = srv.Stop(ctx)
}
