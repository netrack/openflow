package of

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

var (
	ErrHijacked = errors.New("openflow: Connection has been hijacked")
)

type Hijacker interface {
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

// A ResponseWriter interface is used by an OpenFlow handler to
// construct an OpenFlow Response.
type ResponseWriter interface {
	Hijacker
	// Header returns the Header interface that will be sent by
	// WriteHeader. Changing the header after a call to WriteHeader
	// (or Write) has no effect
	Header() Header
	// Write writes the data to the connection as part of an OpenFlow reply.
	Write([]byte) (int, error)
	// WriteHeader sends an Response header as part of an OpenFlow reply.
	WriteHeader() error
	// Close closes connection
	Close() error
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

type Response struct {
	Conn   OFPConn
	header header
	buf    bytes.Buffer
}

func (w *Response) Header() Header {
	return &w.header
}

func (w *Response) Write(b []byte) (n int, err error) {
	return w.buf.Write(b)
}

func (w *Response) WriteHeader() (err error) {
	var buf bytes.Buffer

	w.header.Length = headerlen + uint16(w.buf.Len())
	defer w.buf.Reset()

	_, err = w.header.WriteTo(&buf)
	if err != nil {
		return
	}

	_, err = w.buf.WriteTo(&buf)
	if err != nil {
		return
	}

	_, err = w.Conn.Write(buf.Bytes())
	if err != nil {
		return
	}

	return w.Conn.Flush()
}

func (w *Response) Close() error {
	return w.Close()
}

func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Conn.Hijack()
}

var DefaultServer = Server{
	Addr:    "0.0.0.0:6633",
	Handler: DefaultServeMux,
}

func ListenAndServe() error {
	return DefaultServer.ListenAndServe()
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

		c := NewConn(rwc)
		c.ReadTimeout = srv.ReadTimeout
		c.WriteTimeout = srv.WriteTimeout

		go srv.serve(c, srv.Handler)
	}
}

func (srv *Server) serve(c *Conn, h Handler) {
	origconn := c.rwc
	defer func() {
		if !c.hijacked() {
			origconn.Close()
		}
	}()

	for {
		req, err := c.Receive()
		if err != nil {
			return
		}

		resp := &Response{Conn: c}
		h.Serve(resp, req)

		c.buf.Flush()
	}
}

var DefaultServeMux = NewServeMux()

func Handle(t Type, handler Handler) {
	DefaultServeMux.Handle(t, handler)
}

func HandleFunc(t Type, f func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(t, f)
}

type ServeMux struct {
	mu sync.RWMutex
	m  map[Type][]Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{m: make(map[Type][]Handler)}
}

// Handle registers the handler for the given message type.
func (mux *ServeMux) Handle(t Type, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("mux: nil handler")
	}

	mux.m[t] = append(mux.m[t], handler)
}

func (mux *ServeMux) HandleFunc(t Type, f func(ResponseWriter, *Request)) {
	mux.Handle(t, HandlerFunc(f))
}

func (mux *ServeMux) Handler(r *Request) (Handler, Type) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	handlers, ok := mux.m[r.Header.Type]
	if !ok {
		handlers = append(handlers, DiscardHandler)
	}

	h := HandlerFunc(func(rw ResponseWriter, r *Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		for _, handler := range handlers {
			r.Body = bytes.NewBuffer(body)
			handler.Serve(rw, r)
		}
	})

	return h, r.Header.Type
}

func (mux *ServeMux) Serve(rw ResponseWriter, r *Request) {
	h, _ := mux.Handler(r)
	h.Serve(rw, r)
}
