package ofptest

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/netrack/openflow"
)

// Server is an OpenFlow server listening on a system-chosen port
// on the loopback interface, for use in end-to-end OpenFlow tests.
type Server struct {
	// Listener specifies the server listener. This attribute is
	// optional, if specified, will be used instead of the default.
	Listener net.Listener

	// Config may be changed after calling NewUnstartedServer and
	// before Start.
	Config *of.Server

	closed bool
	conns  map[of.Conn]struct{}

	mu sync.Mutex
}

// NewServer starts and returns a new Server. The caller should call
// Close when finished to shut it down.
func NewServer(handler of.Handler) *Server {
	srv := NewUnstartedServer(handler, nil)
	srv.Start()
	return srv
}

// newLocalListener creates a new listener on loopback interface with
// a system-chosen port.
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

// NewUnstartedServer returns a new Server but doesn't start it.
//
// After changing its configuration, the caller should call Start.
//
// The caller should call Close when finished to shut it down.
func NewUnstartedServer(handler of.Handler, listener net.Listener) *Server {
	if listener == nil {
		listener = newLocalListener()
	}

	return &Server{
		Listener: listener,
		Config:   &of.Server{Handler: handler},
	}
}

// wrap wraps the server ConnState function in order to save all
// incoming connections to clean them up at the server shut down.
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

// Start starts a server on system-chosen port.
func (s *Server) Start() {
	s.wrap()

	go func() {
		s.Config.Serve(s.Listener)
	}()
}

// closeChanConn closes a specified connection an sends a message
// to the cannel to notify the caller about the termination.
func (s *Server) closeChanConn(conn of.Conn, ch chan struct{}) {
	conn.Close()
	ch <- struct{}{}
}

// CloseClientConnections closes all open connections to the Server.
func (s *Server) CloseClientConnections() {
	nconns := len(s.conns)
	ch := make(chan struct{}, nconns)

	for conn := range s.conns {
		go s.closeChanConn(conn, ch)
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

// Close shuts down the server and closes all open connections to the
// Server.
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
