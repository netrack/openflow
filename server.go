package openflow

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrHijacked = errors.New("openflow: Connection has been hijacked")
)

type ResponseWriter interface {
	Write([]byte) (int, error)
}

type Handler interface {
	Serve(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (h HandlerFunc) Serve(rw ResponseWriter, r *Request) {
	h(rw, r)
}

func Discard(rw ResponseWriter, r *Request) {}

var DiscardHandler = HandlerFunc(Discard)

type Hijacker interface {
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

type response struct {
	conn *conn
}

func (w *response) Write(b []byte) (int, error) {
	return w.conn.write(b)
}

func (w *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn.hijack()
}

type conn struct {
	rwc net.Conn
	buf *bufio.ReadWriter

	rtimeout time.Duration
	wtimeout time.Duration

	mu sync.Mutex

	hijackedv bool
}

func (c *conn) hijacked() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hijackedv
}

func (c *conn) hijack() (net.Conn, *bufio.ReadWriter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hijackedv {
		return nil, nil, ErrHijacked
	}

	c.hijackedv = true
	rwc := c.rwc
	buf := c.buf

	c.rwc = nil
	c.buf = nil
	return rwc, buf, nil
}

func (c *conn) read() (*Request, error) {
	if c.hijacked() {
		return nil, ErrHijacked
	}

	if d := c.rtimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}

	if d := c.wtimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	var req Request

	err := req.Read(c.buf)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (c *conn) write(b []byte) (int, error) {
	if c.hijacked() {
		return 0, ErrHijacked
	}

	defer c.buf.Flush()
	return c.buf.Write(b)
}

func (c *conn) serve(h Handler) {
	origconn := c.rwc
	defer func() {
		if !c.hijacked() {
			origconn.Close()
		}
	}()

	for {
		req, err := c.read()
		if err != nil {
			return
		}

		h.Serve(&response{c}, req)
	}
}

type Server struct {
	Addr    string
	Handler Handler

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		rwc, err := l.Accept()
		if err != nil {
			return err
		}

		br := bufio.NewReader(rwc)
		bw := bufio.NewWriter(rwc)

		brw := bufio.NewReadWriter(br, bw)
		c := &conn{
			rwc: rwc, buf: brw,
			rtimeout: srv.ReadTimeout,
			wtimeout: srv.WriteTimeout,
		}

		go c.serve(srv.Handler)
	}
}

var DefaultServeMux = NewServeMux()

type ServeMux struct {
	mu sync.RWMutex
	m  map[Type]Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{m: make(map[Type]Handler)}
}

// Handle registers the handler for the given message type.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(t Type, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("openflow: nil handler")
	}

	if _, dup := mux.m[t]; dup {
		panic("openflow: multiple registrations")
	}

	mux.m[t] = handler
}

func (mux *ServeMux) HandleFunc(t Type, f func(ResponseWriter, *Request)) {
	mux.Handle(t, HandlerFunc(f))
}

func (mux *ServeMux) Handler(r *Request) (Handler, Type) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, ok := mux.m[r.Header.Type]
	if !ok {
		h = DiscardHandler
	}

	return h, r.Header.Type
}

func (mux *ServeMux) Serve(rw ResponseWriter, r *Request) {
	h, _ := mux.Handler(r)
	h.Serve(rw, r)
}
