package ofptest

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/netrack/openflow"
)

type Server struct {
	Listener net.Listener

	Config *of.Server

	closed bool
	conns  map[of.Conn]struct{}

	mu sync.Mutex
}

func NewServer(handler of.Handler) *Server {
	srv := NewUnstartedServer(handler, nil)
	srv.Start()
	return srv
}

func newLocalListener() net.Listener {
	urls := []*url.URL{
		{Scheme: "tcp", Host: "127.0.0.1:0"},
		{Scheme: "tcp6", Host: "[::1]:0"},
	}

	for pos, u := range urls {
		ln, err := net.Listen(u.Scheme, u.Host)
		if err == nil {
			return ln
		}

		if pos == len(urls)-1 {
			text := "ofptest: failed to listen on a port: %v"
			panic(fmt.Errorf(text, err))
		}
	}

	return nil
}

func NewUnstartedServer(handler of.Handler, listener net.Listener) *Server {
	if listener == nil {
		listener = newLocalListener()
	}

	return &Server{
		Listener: listener,
		Config:   &of.Server{Handler: handler},
	}
}

func (s *Server) wrap() {
	oldCb := s.Config.ConnState

	s.Config.ConnState = func(conn of.Conn, state of.ConnState) {
		s.mu.Lock()
		if state == of.StateNew {
			// Persist the new connections, so they
			// could be closed gracefully.
			if s.conns == nil {
				s.conns = make(map[of.Conn]struct{})
			}

			s.conns[conn] = struct{}{}
		}

		s.mu.Unlock()
		if oldCb != nil {
			oldCb(conn, state)
		}
	}

}

func (s *Server) Start() {
	s.wrap()

	go func() {
		s.Config.Serve(s.Listener)
	}()
}

func (s *Server) closeChanConn(conn of.Conn, ch chan struct{}) {
	conn.Close()
	ch <- struct{}{}
}

func (s *Server) CloseClientConnections() {
	nconns := len(s.conns)
	ch := make(chan struct{}, nconns)

	for conn := range s.conns {
		s.closeChanConn(conn, ch)
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for i := 0; i < nconns; i++ {
		select {
		case <-ch:
		case <-timer.C:
			return
		}
	}
}

func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	s.Listener.Close()
	s.CloseClientConnections()
}
