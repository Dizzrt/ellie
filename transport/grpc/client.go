package grpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(false, opts...)
}

func DialInsecure(opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(true, opts...)
}

func dial(isInsecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		subsetSize: 25,
		timeout:    2000 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(&options)
	}

	ints := []grpc.UnaryClientInterceptor{}

	if len(options.unaryClientInts) > 0 {
		ints = append(ints, options.unaryClientInts...)
	}

	sints := []grpc.StreamClientInterceptor{}

	if len(options.streamClientInts) > 0 {
		sints = append(sints, options.streamClientInts...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(sints...),
	}

	if isInsecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}

	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.NewClient(options.endpoint, grpcOpts...)
}
