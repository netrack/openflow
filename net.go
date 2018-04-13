package openflow

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"
)

type ConnState int

const (
	StateNew ConnState = iota
	StateHandshake
	StateActive
	StateIdle
	StateClosed
)

func (c ConnState) String() string {
	text, ok := connStateText[c]
	if !ok {
		return fmt.Sprintf("ConnState(%d)", c)
	}
	return text
}

var connStateText = map[ConnState]string{
	StateNew:       "StateNew",
	StateHandshake: "StateHandshake",
	StateActive:    "StateActive",
	StateIdle:      "StateIdle",
	StateClosed:    "StateClosed",
}

// Conn is an generic OpenFlow connection.
//
// Multiple goroutines may invoke methods on conn simultaneously.
type Conn interface {
	// Receive receives message from input buffer
	Receive() (*Request, error)

	// Send writes message to output buffer
	Send(*Request) error

	// Close closes the connection. Any blocked Read or Write operations
	// will be unblocked and return errors.
	Close() error

	// Flush writes the messages from output buffer to the connection.
	Flush() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr

	// SetDeadline sets the read and write deadlines associated with the
	// connection.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for the future Receive calls.
	// If the deadline is reached, Receive will fail with a timeout (see
	// type Error) instead of blocking.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadLine sets the deadline for the future Send calls.
	// If the deadline is reached, Send will fail with a timeout (see
	// type Error) instead of blocking.
	SetWriteDeadline(t time.Time) error
}

// conn is an OpenFlow protocol connection.
type conn struct {
	// A read-write connection.
	rwc net.Conn

	// An input and output buffer.
	buf *bufio.ReadWriter
	mu  sync.Mutex

	// Maximum duration before timing out the read of the request.
	ReadTimeout time.Duration
	// Maximum duration before timing out the write of the response.
	WriteTimeout time.Duration
}

// NewConn creates a new OpenFlow protocol connection.
func NewConn(c net.Conn) Conn {
	return newConn(c)
}

// newConn creates a new OpenFlow protocol connection from the
// given one.
func newConn(c net.Conn) *conn {
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	brw := bufio.NewReadWriter(br, bw)
	return &conn{rwc: c, buf: brw}
}

// Read reads data from the connection.
func (c *conn) Read(b []byte) (int, error) {
	return c.buf.Read(b)
}

// Receive reads OpenFlow data from the connection.
func (c *conn) Receive() (*Request, error) {
	if d := c.ReadTimeout; d != 0 {
		c.SetReadDeadline(time.Now().Add(d))
	}

	r := &Request{Addr: c.rwc.RemoteAddr(), conn: c}
	_, err := r.ReadFrom(c)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Write writes data to the connection. Write can be made to time out.
func (c *conn) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Write(b)
}

// Flush writes any buffered data to the connection.
func (c *conn) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Flush()
}

// forceWrite writes given data and any buffered data to the connection.
func (c *conn) forceWrite(b []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.buf.Write(b)
	if err != nil {
		return err
	}

	return c.buf.Flush()
}

// Send writes OpenFlow data to the connection.
func (c *conn) Send(r *Request) error {
	if d := c.WriteTimeout; d != 0 {
		defer func() {
			c.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	_, err := r.WriteTo(c)
	return err
}

// Close closes the connection. Any blocked Read or Write operations will
// be unblocked and return errors.
func (c *conn) Close() error {
	return c.rwc.Close()
}

// LocalAddr returns the local network address.
func (c *conn) LocalAddr() net.Addr {
	return c.rwc.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *conn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated with the
// connection.
func (c *conn) SetDeadline(t time.Time) error {
	return c.rwc.SetDeadline(t)
}

// SetReadDeadline sets the deadline for the future Receive calls. If the
// deadline is reached, Receive will fail with a timeout (see type Error)
// instead of blocking.
func (c *conn) SetReadDeadline(t time.Time) error {
	return c.rwc.SetReadDeadline(t)
}

// SetWriteDeadLine sets the deadline for the future Send calls. If the
// deadline is reached, Send will fail with a timeout (see type Error)
// instead of blocking.
func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.rwc.SetWriteDeadline(t)
}

// Send allows to send multiple requests at once to the connection.
//
// The requests will be written to the per-call buffer. If the
// serialization of all given requests succeeded it will be flushed
// to the OpenFlow connection.
//
// No data will be written when any of the request failed.
func Send(c Conn, requests ...*Request) error {
	for _, request := range requests {
		if err := c.Send(request); err != nil {
			return err
		}
	}

	return c.Flush()
}

// Dial establishes the remote connection to the address on the
// given network.
func Dial(network, addr string) (Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

// DialTLS establishes the remote connection to the address on the
// given network and then initiates TLS handshake, returning the
// resulting TLS connection.
func DialTLS(network, addr string, config *tls.Config) (Conn, error) {
	conn, err := tls.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

// Listener is an OpenFlow network listener. Clients should typically
// use variables of type net.Listener instead of assuming OFP.
type Listener interface {
	// Accept waits for and returns the next connection.
	Accept() (Conn, error)

	// Close closes the listener.
	Close() error

	// Addr returns the listener network address.
	Addr() net.Addr
}

// listener implements the Listener interface.
type listener struct {
	ln net.Listener
}

// NewListener creates a new instance of the OpenFlow listener from
// the given network listener.
//
// This function could be used to establish the communication channel
// over non-TCP sockets, like Unix sockets.
func NewListener(ln net.Listener) Listener {
	return &listener{ln}
}

// Accept waits for and returns the next connection to the listener.
func (l *listener) Accept() (Conn, error) {
	conn, err := l.ln.Accept()
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

// Close closes an OpenFlow server connection.
func (l *listener) Close() error {
	return l.ln.Close()
}

// Addr returns the network address of the listener.
func (l *listener) Addr() net.Addr {
	return l.ln.Addr()
}

// Listen announces on the local network address laddr.
func Listen(network, laddr string) (Listener, error) {
	tcpaddr, err := net.ResolveTCPAddr(network, laddr)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP(network, tcpaddr)
	if err != nil {
		return nil, err
	}

	return &listener{ln}, err
}

// ListenTLS announces on the local network address.
func ListenTLS(network, laddr string, config *tls.Config) (Listener, error) {
	ln, err := tls.Listen(network, laddr, config)
	if err != nil {
		return nil, err
	}

	return &listener{ln}, err
}
