package openflow

import (
	"bytes"
	"io"
	"net"
	"sync"
	"sync/atomic"
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

	return r.conn.forceWrite(r.buf.Bytes())
}

// The reqwrap defines a placeholder for request and error returned
// from Receive call.
type reqwrap struct {
	req *Request
	err error
}

// The receiver defines a wrapper around connection used to create
// a channel of Requests by continuously fetching data from connection.
type receiver struct {
	// Conn is a client connection.
	Conn *conn

	once sync.Once
	ch   chan reqwrap
}

// The receive starts an infinite loop of reading the requests from
// the client connection. When the receive call returns an error, a
// submission channel closes and loop terminates.
func (r *receiver) receive() {
	for {
		req, err := r.Conn.Receive()
		r.ch <- reqwrap{req, err}

		if err != nil {
			close(r.ch)
			return
		}
	}
}

// The starts initializes all required attributes of the instance and
// spawns a new goroutine to receive requests from the channel.
func (r *receiver) start() {
	r.ch = make(chan reqwrap)
	go r.receive()
}

// C returns a read-only channel of request plus error.
func (r *receiver) C() <-chan reqwrap {
	// Wait for a new request from the client.
	r.once.Do(r.start)
	return r.ch
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
//	of.ListenAndServe(":6633", nil)
//
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

// A Server defines parameters for running OpenFlow server.
type Server struct {
	// Addr is an address to listen on.
	Addr string

	// Runner defines concurrency model of handling requests from the
	// switch within a single connection. By default OnDemandRoutineRunner
	// is used.
	Runner Runner

	// Handler to invoke on the incoming requests.
	Handler Handler

	// Maximum duration before timing out the read of the request.
	ReadTimeout time.Duration

	// Maximum duration before timing out the write of the response.
	WriteTimeout time.Duration

	// ConnState specifies an optional callback function that is called
	// when a client connection changes state.
	ConnState func(Conn, ConnState)

	// MaxConns defines the maximum number of client connections server
	// handles, the rest will be explicitly closed. Zero means no limit.
	MaxConns int

	// The conns store the count of the client connections. This value
	// is incremented on each new connection and decremented on each
	// closed connection.
	conns int32

	// The stop is a channel used to terminate the server in a hard
	// way. It will be used only for testing purposes.
	stop chan struct{}
	once sync.Once
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
		handler = DefaultMux
	}

	runner := srv.Runner
	if runner == nil {
		runner = OnDemandRoutineRunner{}
	}

	// Initialize a channel used to terminate the main handling loop.
	srv.once.Do(func() { srv.stop = make(chan struct{}) })

	for {
		err := srv.accept(l, runner, handler)
		if err != nil {
			return err
		}
	}
}

// The accept block until a new connection will be extracted from the
// queue. It keeps track of count of incoming connections and closes all
// that exceed the MaxConns threshold.
func (srv *Server) accept(l net.Listener, r Runner, h Handler) error {
	rwc, err := l.Accept()
	if err != nil {
		return err
	}

	c := newConn(rwc)
	c.ReadTimeout = srv.ReadTimeout
	c.WriteTimeout = srv.WriteTimeout

	srv.setState(c, StateNew)

	// Terminate the
	numConns := atomic.LoadInt32(&srv.conns)
	if srv.MaxConns != 0 && numConns >= int32(srv.MaxConns) {
		// Make the close a deferred call in case of overridden ConnState
		// function produce a panic error (to prevent file descriptor leak).
		defer c.Close()
		srv.setState(c, StateClosed)
		return nil
	}

	atomic.AddInt32(&srv.conns, 1)
	go srv.serve(c, r, h)

	return nil
}

func (srv *Server) serve(c *conn, r Runner, h Handler) {
	// Define a deferred call to close the connection.
	defer c.Close()
	defer srv.setState(c, StateClosed)
	rcvr := &receiver{Conn: c}

	for {
		select {
		// Receive new requests in the infinite loop and handle them
		// according to the Runner algorithm. When error is returned
		// from the receive loop, this routine will be stopped and the
		// client connection closed.
		case entry := <-rcvr.C():
			if entry.err != nil {
				atomic.AddInt32(&srv.conns, -1)
				return
			}

			r.Run(func() { srv.serveReq(c, entry.req, h) })

		// Stop channel used here only for testing purposes to terminate the
		// infinite receive loop. As a result the all client connections will
		// be closed, as well as receive routine.
		case <-srv.stop:
			return
		}
	}
}

// The serveReq serves a single request from the given connection using
// specified handler.
func (srv *Server) serveReq(c *conn, req *Request, h Handler) {
	state := StateActive
	if req.Header.Type == TypeHello {
		state = StateHandshake
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
	// pending messages will be written to the wire.
	c.Flush()

	// Update the state as the handler just processed the
	// received request and now the server will wait for a
	// new one.
	srv.setState(c, StateIdle)
}

// The close closes all running handlers of the client connections.
func (srv *Server) close() {
	close(srv.stop)
}
