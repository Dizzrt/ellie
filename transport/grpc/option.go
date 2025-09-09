package grpc

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"google.golang.org/grpc"
)

// region ServerOption

type ServerOption func(*Server)

func TLSConfig(tlsConf *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = tlsConf
	}
}

func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func Address(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func Endpoint(endpoint *url.URL) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// TODO logger
// TODO middleware

func UnaryInterceptor(ints ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = ints
	}
}

func StreamInterceptor(ints ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = ints
	}
}

func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

func CustomHealth() ServerOption {
	return func(s *Server) {
		s.customHealth = true
	}
}

func DisableReflection() ServerOption {
	return func(s *Server) {
		s.disableReflection = true
	}
}

// endregion

// region ClientOption

type ClientOption func(*clientOptions)

type clientOptions struct {
	endpoint   string
	subsetSize int
	tlsConf    *tls.Config
	timeout    time.Duration
	// discovery  registry.Discovery
	// middleware
	// streamMiddleware
	unaryClientInts  []grpc.UnaryClientInterceptor
	streamClientInts []grpc.StreamClientInterceptor
	grpcOpts         []grpc.DialOption
	// balancerName
	// filters
	// healthCheckConfig
	// printDiscoveryDebugLog
}

func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithSubsetSize(size int) ClientOption {
	return func(o *clientOptions) {
		o.subsetSize = size
	}
}

func WithTLSConfig(tlsConf *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = tlsConf
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// func WithDiscovery(discovery registry.Discovery) ClientOption {
// 	return func(o *clientOptions) {
// 		o.discovery = discovery
// 	}
// }

func WithUnaryClientInterceptor(ints ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryClientInts = ints
	}
}

func WithStreamClientInterceptor(ints ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamClientInts = ints
	}
}

func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// endregion
