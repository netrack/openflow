package of

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"
)

// A ResponseWriter interface is used by an OpenFlow handler to construct
// an OpenFlow Response.
type ResponseWriter interface {
	// Write writes the data to the connection as part of an OpenFlow reply.
	Write(*Header, io.WriterTo) error
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

// response implements ResponseWriter interface and represents the
// response from an OpenFlow request.
type response struct {
	// conn is an OpenFlow connection instance.
	conn *conn

	// Request header, that could be used to configure a few of
	// response attributes.
	header Header

	// The buf is a response message buffer. We use this buffer
	// to calculate the length of the payload.
	buf   bytes.Buffer
	bufMu sync.Mutex
}

// Write sends an OpenFlow response.
func (r *response) Write(header *Header, w io.WriterTo) (err error) {
	var buf bytes.Buffer

	// When the version is not configured properly, we will
	// use the version from the response header, to minimize
	// manual configuration.
	if header.Version == 0 {
		header.Version = r.header.Version
	}

	r.bufMu.Lock()
	defer r.bufMu.Unlock()
	defer r.buf.Reset()

	if w != nil {
		if _, err = w.WriteTo(&buf); err != nil {
			return
		}
	}

	header.Length = headerlen + uint16(buf.Len())
	_, err = header.WriteTo(&r.buf)
	if err != nil {
		return
	}

	_, err = buf.WriteTo(&r.buf)
	if err != nil {
		return
	}

	_, err = r.conn.Write(r.buf.Bytes())
	if err != nil {
		return
	}

	return r.conn.Flush()
}

// ListenAndServe listens on the given TCP address the handler. When
// handler set to nil, the default handler will be used.
//
// A trivial example is:
//
//	of.HandleFunc(of.TypeEchoRequest, func(rw of.ResponseWriter, r *of.Request) {
//		rw.Write(&of.Header{Type: of.TypeEchoReply}, nil)
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

	ConnState func(Conn, ConnState)
}

func (srv *Server) setState(conn Conn, state ConnState) {
	if cb := srv.ConnState; cb != nil {
		cb(conn, state)
	}
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

		c := newConn(rwc)
		c.ReadTimeout = srv.ReadTimeout
		c.WriteTimeout = srv.WriteTimeout

		srv.setState(c, StateNew)
		go srv.serve(c, handler)
	}
}

func (srv *Server) serve(c *conn, h Handler) {
	// Define a deferred call to close the connection.
	defer c.Close()

	for {
		// Wait for the new request from either the Switch or
		// the Controller.
		req, err := c.Receive()
		if err != nil {
			return
		}

		state := StateActive
		if req.Header.Type == TypeHello {
			state = StateHelloReceived
		}

		srv.setState(c, state)
		// Define a response version from the request version, so
		// it will potentially reduce the amount of additional
		// header configurations.
		header := Header{
			Version:     req.Header.Version,
			Transaction: req.Header.Transaction,
		}

		// Construct a new response instance with some default
		// attributes and execute respective handler for it.
		resp := &response{conn: c, header: header}
		h.Serve(resp, req)

		// Write the buffer content to the connection, so the
		// pending messages will written.
		c.Flush()

		// Update the state as the handler just processed the
		// received request and now the server will wait for a
		// new one.
		srv.setState(c, StateIdle)
	}
}
