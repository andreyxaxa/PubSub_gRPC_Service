package grpcserver

import (
	"fmt"
	"net"
	"time"

	pbgrpc "google.golang.org/grpc"
)

const (
	_defaultAddr = ":80"
)

type Server struct {
	App     *pbgrpc.Server
	notify  chan error
	address string
}

func New(opts ...Option) *Server {
	s := &Server{
		App:     pbgrpc.NewServer(),
		notify:  make(chan error, 1),
		address: _defaultAddr,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		s.notify <- fmt.Errorf("failed to listen: %w", err)
		close(s.notify)

		return
	}

	s.notify <- s.App.Serve(ln)
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	done := make(chan struct{})

	go func() {
		s.App.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(15 * time.Second):
		s.App.Stop()
	}

	return nil
}
