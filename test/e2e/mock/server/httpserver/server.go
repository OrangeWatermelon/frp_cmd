package httpserver

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	bindAddr string
	bindPort int
	handler  http.Handler

	l         net.Listener
	tlsConfig *tls.Config
	hs        *http.Server
}

type Option func(*Server) *Server

func New(options ...Option) *Server {
	s := &Server{
		bindAddr: "127.0.0.1",
	}

	for _, option := range options {
		s = option(s)
	}
	return s
}

func WithBindAddr(addr string) Option {
	return func(s *Server) *Server {
		s.bindAddr = addr
		return s
	}
}

func WithBindPort(port int) Option {
	return func(s *Server) *Server {
		s.bindPort = port
		return s
	}
}

func WithTlsConfig(tlsConfig *tls.Config) Option {
	return func(s *Server) *Server {
		s.tlsConfig = tlsConfig
		return s
	}
}

func WithHandler(h http.Handler) Option {
	return func(s *Server) *Server {
		s.handler = h
		return s
	}
}

func WithResponse(resp []byte) Option {
	return func(s *Server) *Server {
		s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(resp)
		})
		return s
	}
}

func (s *Server) Run() error {
	if err := s.initListener(); err != nil {
		return err
	}

	addr := net.JoinHostPort(s.bindAddr, strconv.Itoa(s.bindPort))
	hs := &http.Server{
		Addr:      addr,
		Handler:   s.handler,
		TLSConfig: s.tlsConfig,
	}

	s.hs = hs
	if s.tlsConfig == nil {
		go hs.Serve(s.l)
	} else {
		go hs.ServeTLS(s.l, "", "")
	}
	return nil
}

func (s *Server) Close() error {
	if s.hs != nil {
		return s.hs.Close()
	}
	return nil
}

func (s *Server) initListener() (err error) {
	s.l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.bindAddr, s.bindPort))
	return
}

func (s *Server) BindAddr() string {
	return s.bindAddr
}

func (s *Server) BindPort() int {
	return s.bindPort
}
