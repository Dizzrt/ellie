package consul

import (
	"context"
	"reflect"
	"testing"

	"github.com/dizzrt/ellie/registry"
	"github.com/hashicorp/consul/api"
)

func TestRegistry_Register(t *testing.T) {
	opts := []Option{
		WithHealthCheck(false),
	}

	type args struct {
		ctx        context.Context
		serverName string
		server     []*registry.ServiceInstance
	}

	test := []struct {
		name    string
		args    args
		want    []*registry.ServiceInstance
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:        context.Background(),
				serverName: "server-1",
				server: []*registry.ServiceInstance{
					{
						ID:        "1",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "1",
					Name:      "server-1",
					Version:   "v0.0.1",
					Metadata:  nil,
					Endpoints: []string{"http://127.0.0.1:8000"},
				},
			},
			wantErr: false,
		},
		{
			name: "registry new service replace old service",
			args: args{
				ctx:        context.Background(),
				serverName: "server-1",
				server: []*registry.ServiceInstance{
					{
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
					{
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.2",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "2",
					Name:      "server-1",
					Version:   "v0.0.2",
					Metadata:  nil,
					Endpoints: []string{"http://127.0.0.1:8000"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := api.NewClient(&api.Config{Address: "192.168.124.10:8500"})
			if err != nil {
				t.Fatalf("create consul client failed: %v", err)
			}

			r := New(cli, opts...)

			for _, instance := range tt.args.server {
				instance := instance
				err = r.Register(tt.args.ctx, instance)
				if err != nil {
					t.Error(err)
				}
				defer func() {
					err = r.Deregister(tt.args.ctx, instance)
					if err != nil {
						t.Error(err)
					}
				}()
			}

			watchCtx, watchCancel := context.WithCancel(context.Background())
			watch, err := r.Watch(watchCtx, tt.args.serverName)
			if err != nil {
				t.Error(err)
				watchCancel()
				return
			}

			got, err := watch.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", got)
				watchCancel()
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService() got = %v, want %v", got, tt.want)
			}

			err = watch.Stop()
			if err != nil {
				t.Error(err)
			}
			watchCancel()
		})
	}
}
