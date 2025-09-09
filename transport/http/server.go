package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Dizzrt/ellie/internal/endpoint"
	"github.com/Dizzrt/ellie/internal/host"
	"github.com/Dizzrt/ellie/transport"
	"github.com/gorilla/mux"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

type Server struct {
	*http.Server

	err    error
	lis    net.Listener
	router *mux.Router

	tlsConf  *tls.Config
	endpoint *url.URL
	network  string
	address  string
	timeout  time.Duration
	filters  []FilterFunc
	// TODO middleware
	pathParamsDecoder  HTTPCodecRequestDecoder
	queryParamsDecoder HTTPCodecRequestDecoder
	requestBodyDecoder HTTPCodecRequestDecoder
	responseEncoder    HTTPCodecResponseEncoder
	errorEncoder       HTTPCodecErrorEncoder
	strictSlash        bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:            "tcp",
		address:            ":0",
		timeout:            1 * time.Second,
		pathParamsDecoder:  DefaultPathParamsDecoder,
		queryParamsDecoder: DefaultQueryParamsDecoder,
		requestBodyDecoder: DefaultRequestBodyDecoder,
		responseEncoder:    DefaultResponseEncoder,
		errorEncoder:       DefaultErrorEncoder,
		strictSlash:        true,
		router:             mux.NewRouter(),
	}

	srv.router.NotFoundHandler = http.DefaultServeMux
	srv.router.MethodNotAllowedHandler = http.DefaultServeMux
	for _, opt := range opts {
		opt(srv)
	}

	srv.router.StrictSlash(srv.strictSlash)
	srv.Server = &http.Server{
		TLSConfig: srv.tlsConf,
		Handler:   FilterChain(srv.filters...)(srv.router),
	}

	return srv
}

func (s *Server) HandlePrefix(prefix string, h http.Handler) {
	s.router.PathPrefix(prefix).Handler(h)
}

func (s *Server) initializeListenerAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}

		s.lis = lis
	}

	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}

		s.endpoint = endpoint.New(endpoint.Scheme("http", s.tlsConf != nil), addr)
	}

	return s.err
}

// region interfaces impl

func (s *Server) Start(ctx context.Context) error {
	if err := s.initializeListenerAndEndpoint(); err != nil {
		return err
	}

	s.BaseContext = func(l net.Listener) context.Context {
		return ctx
	}

	// TODO log
	fmt.Printf("[][HTTP] server listening on %s\n", s.lis.Addr().String())

	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO log
	fmt.Println("[HTTP] server stopping")

	err := s.Shutdown(ctx)
	if err != nil {
		if ctx.Err() != nil {
			// TODO log
			fmt.Println("[HTTP] server couldn't stop gracefully in time, forcing stop")
			err = s.Close()
		}
	}

	return err
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.initializeListenerAndEndpoint(); err != nil {
		return nil, s.err
	}

	return s.endpoint, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

// endregion
