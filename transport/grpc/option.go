package grpc

import (
	"net"
	"net/url"
	"time"

	"google.golang.org/grpc"
)

type ServerOption func(*Server)

// TODO tls

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

// TODO health

func DisableReflection() ServerOption {
	return func(s *Server) {
		s.disableReflection = true
	}
}
