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

// A ResponseWriter interface is used by an OpenFlow handler to construct
// an OpenFlow Response.
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

// A Handler responds to an OpenFlow request.
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
	// Define a deferred function to close the connection.
	defer func() {
		if !c.hijacked() {
			origconn.Close()
		}
	}()

	for {
		// Wait for the new request from either the Switch or
		// the Controller.
		req, err := c.Receive()
		if err != nil {
			return
		}

		// Define a response version from the request version, so
		// it will potentially reduce the amount of additional
		// header configurations.
		header := header{Version: req.Header.Version}

		// Construct a new response instance with some default
		// attributes and execute respective handler for it.
		resp := &Response{Conn: c, header: header}
		h.Serve(resp, req)

		// Write the buffer content to the connection, so the
		// pending messages will written.
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
	m  map[Type][]*muxEntry
}

type muxEntry struct {
	h Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{m: make(map[Type][]*muxEntry)}
}

// Handle registers the handler for the given message type.
func (mux *ServeMux) Handle(t Type, handler Handler) *muxEntry {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("mux: nil handler")
	}

	entry := &muxEntry{h: handler}
	mux.m[t] = append(mux.m[t], entry)
	return entry
}

func (mux *ServeMux) HandleFunc(t Type, f HandlerFunc) *muxEntry {
	return mux.Handle(t, f)
}

func (mux *ServeMux) Unhandle(t Type, entry *muxEntry) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	typeEntries, ok := mux.m[t]
	if !ok {
		return
	}

	var entries []*muxEntry
	for _, e := range typeEntries {
		// Remove all of them
		if e == entry {
			continue
		}

		entries = append(entries, e)
	}

	mux.m[t] = entries
}

func (mux *ServeMux) Handler(r *Request) (Handler, Type) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	entries, ok := mux.m[r.Header.Type]
	if !ok {
		entries = append(entries, &muxEntry{h: DiscardHandler})
	}

	h := HandlerFunc(func(rw ResponseWriter, r *Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		for _, entry := range entries {
			r.Body = bytes.NewBuffer(body)
			entry.h.Serve(rw, r)
		}
	})

	return h, r.Header.Type
}

func (mux *ServeMux) Serve(rw ResponseWriter, r *Request) {
	h, _ := mux.Handler(r)
	h.Serve(rw, r)
}
