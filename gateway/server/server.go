package server

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	name string
	srv  *http.Server
	addr string
}

func New(name, addr string, tlsConfig *tls.Config, handler http.Handler) *Server {
	s := &Server{name: name}
	s.srv = &http.Server{Addr: addr, Handler: handler, TLSConfig: tlsConfig, ReadHeaderTimeout: 2 * time.Second}
	return s
}

func (s *Server) Start(errCh chan error) {
	l, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		errCh <- err
		return
	}
	s.addr = l.Addr().String()

	log.Printf("%s listening on %s", s.name, l.Addr())

	go func() {
		if s.srv.TLSConfig != nil {
			errCh <- s.srv.ServeTLS(l, "", "")
		} else {
			errCh <- s.srv.Serve(l)
		}
	}()
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
