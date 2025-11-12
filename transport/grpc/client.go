package grpc

import (
	"fmt"
	"time"

	"github.com/dizzrt/ellie/transport/grpc/resolver/discovery"
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
		subsetSize:             25,
		timeout:                2000 * time.Millisecond,
		printDiscoveryDebugLog: true,
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

	if options.discovery != nil {
		resolvers := grpc.WithResolvers(discovery.NewBuilder(
			options.discovery,
			discovery.WithInsecure(isInsecure),
			discovery.WithTimeout(options.timeout),
			discovery.WithSubsetSize(options.subsetSize),
			discovery.WithDebugLog(options.printDiscoveryDebugLog),
		))

		grpcOpts = append(grpcOpts, resolvers)
	}

	if isInsecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else if options.tlsConf != nil {
		temp := credentials.NewTLS(options.tlsConf)
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(temp))
	} else {
		return nil, fmt.Errorf("tls config is nil")
	}

	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.NewClient(options.endpoint, grpcOpts...)
}
