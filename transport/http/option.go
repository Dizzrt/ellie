package http

import (
	"crypto/tls"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerOption func(*Server)

func TLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = tlsConfig
	}
}

func Endpoint(endpoint *url.URL) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
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

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func DefaultSuccessCode(code int) ServerOption {
	return func(s *Server) {
		s.defaultSuccessCode = code
	}
}

func DefaultSuccessMessage(message string) ServerOption {
	return func(s *Server) {
		s.defaultSuccessMessage = message
	}
}

func ResponseEncoder(encoder HTTPResponseEncoder) ServerOption {
	return func(s *Server) {
		s.responseEncoder = encoder
	}
}

func RedirectTrailingSlash(isStrict bool) ServerOption {
	return func(s *Server) {
		s.redirectTrailingSlash = isStrict
	}
}

func NoRouteHandlers(handler gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		s.noRouteHandlers = append(s.noRouteHandlers, handler)
	}
}

func NoMethodHandlers(handler gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		s.NoMethodHandler = append(s.NoMethodHandler, handler)
	}
}
