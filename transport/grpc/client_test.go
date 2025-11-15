package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dizzrt/ellie/contrib/registry/consul"
	"github.com/dizzrt/ellie/internal/mock/ping"
	"github.com/hashicorp/consul/api"
)

func TestDiscovery(t *testing.T) {
	var err error
	var client *api.Client
	if client, err = api.NewClient(&api.Config{Address: "dev.ellie.com:8500"}); err != nil {
		t.Fatal(err)
	}

	dis := consul.New(client)
	endpoint := "discovery:///test-app"

	conn, err := DialInsecure(
		WithEndpoint(endpoint),
		WithDiscovery(dis),
		WithPrintDiscoveryDebugLog(true),
		WithTimeout(5*time.Second),
	)

	if err != nil {
		t.Fatal(err)
	}

	pingClient := ping.NewPingServiceClient(conn)
	res, err := pingClient.Ping(context.Background(), &ping.PingRequest{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
