package http_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	nhttp "net/http"

	"github.com/Dizzrt/ellie/internal/mock/ping"
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

	url := e.String() + "/ping"
	resp, err := nhttp.Get(url)
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
