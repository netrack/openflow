package of

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"time"
)

var (
	ErrHijacked = errors.New("openflow: Connection has been hijacked")
)

type Hijacker interface {
	// Hijack lets the caller take over the connection.
	//
	// After a call to Hijack(), the OFP server library will not do
	// anything else with a connection.
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

// A ResponseWriter interface is used by an OpenFlow handler to construct
// an OpenFlow Response.
type ResponseWriter interface {
	// Hijack lets the caller take over the connection.
	Hijacker

	// Write writes the data to the connection as part of an OpenFlow reply.
	Write([]byte) (int, error)

	// WriteHeader sends an Response header as part of an OpenFlow reply.
	WriteHeader(*Header) error

	// Close closes connection.
	Close() error
}

// A Handler responds to an OpenFlow request.
//
// Serve should write the reply headers and the payload to the
// ResponseWriter and then return. Returning signals that the request is
// finished.
type Handler interface {
	Serve(ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow use of ordinary functions
// as OpenFlow handlers.
type HandlerFunc func(ResponseWriter, *Request)

// Serve calls f(rw, r).
func (h HandlerFunc) Serve(rw ResponseWriter, r *Request) {
	h(rw, r)
}

// DiscardHandler is a Handler instance to discard the remote OpenFlow
// requests.
var DiscardHandler = HandlerFunc(func(rw ResponseWriter, r *Request) {})

// Response implements ResponseWriter interface and represents the
// response from an OpenFlow request.
type Response struct {
	// Conn is an OpenFlow connection instance.
	Conn Conn

	// Request header, that could be used to configure a few of
	// response attributes.
	header Header

	// The buf is a response message buffer. We use this buffer
	// to calculate the length of the payload.
	buf bytes.Buffer
}

// Write writes the given byte slice to the output buffer.
func (w *Response) Write(b []byte) (n int, err error) {
	return w.buf.Write(b)
}

// WriteHeader sends an OpenFlow response.
func (w *Response) WriteHeader(header *Header) (err error) {
	var buf bytes.Buffer

	// When the version is not configured properly, we will
	// use the version from the response header, to minimize
	// manual configuration.
	if header.Version == 0 {
		header.Version = w.header.Version
	}

	header.Length = headerlen + uint16(w.buf.Len())
	defer w.buf.Reset()

	_, err = header.WriteTo(&buf)
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

// Close closes connection.
func (w *Response) Close() error {
	return w.Close()
}

// Hijack lets the caller take over the connection.
func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Conn.Hijack()
}

// ListenAndServe listens on the given TCP address the handler. When
// handler set to nil, the default handler will be used.
//
// A trivial example is:
//
//	of.HandleFunc(of.TypeEchoRequest, func(rw of.ResponseWriter, r *of.Request) {
//		rw.WriteHeader(&of.Header{Type: of.TypeEchoReply})
//	})
//
//	ListenAndServe(":6633", nil)
//
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

// A Server defines parameters for running OpenFlow server.
type Server struct {
	// Addr is an address to listen on.
	Addr string

	// Handler to invoke on the incoming requests.
	Handler Handler

	// Maximum duration before timing out the read of the request.
	ReadTimeout time.Duration

	// Maximum duration before timing out the write of the response.
	WriteTimeout time.Duration
}

// ListenAndServe listens on the network address and then calls Server
// to handle requests on the incoming connections.
func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	// When the handler is not specified, the default dispatcher
	// will be used instead.
	handler := srv.Handler
	if handler == nil {
		handler = DefaultDispatcher
	}

	for {
		rwc, err := l.Accept()
		if err != nil {
			return err
		}

		c := NewConn(rwc)
		c.ReadTimeout = srv.ReadTimeout
		c.WriteTimeout = srv.WriteTimeout

		go srv.serve(c, handler)
	}
}

func (srv *Server) serve(c *OFPConn, h Handler) {
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
		header := Header{Version: req.Header.Version}

		// Construct a new response instance with some default
		// attributes and execute respective handler for it.
		resp := &Response{Conn: c, header: header}
		h.Serve(resp, req)

		// Write the buffer content to the connection, so the
		// pending messages will written.
		c.buf.Flush()
	}
}
