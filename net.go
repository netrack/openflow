package of

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

// Conn is an generic OpenFlow connection.
//
// Multiple goroutines may invoke methods on OFPConn simultaneously.
type Conn interface {
	net.Conn

	// Hijack lets the caller take over the connection.
	Hijacker

	// Receive receives message from input buffer
	Receive() (*Request, error)

	// Send writes message to output buffer
	Send(*Request) error

	// Flush writes the messages from output buffer to the connection.
	Flush() error
}

// OFPConn is an OpenFlow protocol connection.
type OFPConn struct {
	// A read-write connection.
	rwc net.Conn

	// An input and output buffer.
	buf *bufio.ReadWriter
	// Maximum duration before timing out the read of the request.
	ReadTimeout time.Duration

	// Maximum duration before timing out the write of the response.
	WriteTimeout time.Duration

	// A mutex to access the hijack-ed flag
	mu sync.RWMutex
	// A flag set when the connection hijacked.
	hijackedv bool
}

// NewConn creates a new OpenFlow protocol connection.
func NewConn(conn net.Conn) *OFPConn {
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	brw := bufio.NewReadWriter(br, bw)
	return &OFPConn{rwc: conn, buf: brw}
}

// Returns true when the connection hijacked.
func (c *OFPConn) hijacked() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hijackedv
}

// Hijack takes over the connection.
func (c *OFPConn) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Return an error, when connection already hijacked.
	if c.hijackedv {
		return nil, nil, ErrHijacked
	}

	// Mark the connection hijacked.
	c.hijackedv = true
	rwc := c.rwc
	buf := c.buf

	c.rwc = nil
	c.buf = nil
	return rwc, buf, nil
}

// Read reads data from the connection.
func (c *OFPConn) Read(b []byte) (int, error) {
	if c.hijacked() {
		return 0, ErrHijacked
	}

	return c.buf.Read(b)
}

// Receive reads OpenFlow data from the connection.
func (c *OFPConn) Receive() (*Request, error) {
	if c.hijacked() {
		return nil, ErrHijacked
	}

	if d := c.ReadTimeout; d != 0 {
		c.SetReadDeadline(time.Now().Add(d))
	}

	r := &Request{Addr: c.rwc.RemoteAddr()}
	_, err := r.ReadFrom(c.buf)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Write writes data to the connection. Write can be made to time out.
func (c *OFPConn) Write(b []byte) (int, error) {
	if c.hijacked() {
		return 0, ErrHijacked
	}

	return c.buf.Write(b)
}

// Flush writes any buffered data to the connection.
func (c *OFPConn) Flush() error {
	return c.buf.Flush()
}

// Send writes OpenFlow data to the connection.
func (c *OFPConn) Send(r *Request) error {
	if c.hijacked() {
		return ErrHijacked
	}

	if d := c.WriteTimeout; d != 0 {
		defer func() {
			c.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	_, err := r.WriteTo(c.buf)
	return err
}

// Close closes the connection. Any blocked Read or Write operations will
// be unblocked and return errors.
func (c *OFPConn) Close() error {
	return c.rwc.Close()
}

// LocalAddr returns the local network address.
func (c *OFPConn) LocalAddr() net.Addr {
	return c.rwc.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *OFPConn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated with the
// connection.
func (c *OFPConn) SetDeadline(t time.Time) error {
	return c.rwc.SetDeadline(t)
}

// SetReadDeadline sets the deadline for the future Receive calls. If the
// deadline is reached, Receive will fail with a timeout (see type Error)
// instead of blocking.
func (c *OFPConn) SetReadDeadline(t time.Time) error {
	return c.rwc.SetReadDeadline(t)
}

// SetWriteDeadLine sets the deadline for the future Send calls. If the
// deadline is reached, Send will fail with a timeout (see type Error)
// instead of blocking.
func (c *OFPConn) SetWriteDeadline(t time.Time) error {
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
	var buf bytes.Buffer

	for _, request := range requests {
		if _, err := request.WriteTo(&buf); err != nil {
			return err
		}
	}

	if _, err := buf.WriteTo(c); err != nil {
		return err
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

// OFPListener is an OpenFlow network listener. Clients should typically
// use variables of type net.Listener instead of assuming OFP.
type OFPListener struct {
	ln net.Listener
}

// Accept waits for and returns the next connection to the listener.
func (l *OFPListener) Accept() (net.Conn, error) {
	return l.AcceptOFP()
}

// Accepts accepts the next incoming call and returns new connection.
func (l *OFPListener) AcceptOFP() (*OFPConn, error) {
	conn, err := l.ln.Accept()
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

// Close closes an OpenFlow server connection.
func (l *OFPListener) Close() error {
	return l.ln.Close()
}

// Addr returns the network address of the listener.
func (l *OFPListener) Addr() net.Addr {
	return l.ln.Addr()
}

// Listen announces on the local network address laddr.
func Listen(network, laddr string) (*OFPListener, error) {
	tcpaddr, err := net.ResolveTCPAddr(network, laddr)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP(network, tcpaddr)
	if err != nil {
		return nil, err
	}

	return &OFPListener{ln}, err
}

// ListenTLS announces on the local network address.
func ListenTLS(network, laddr string, config *tls.Config) (*OFPListener, error) {
	ln, err := tls.Listen(network, laddr, config)
	if err != nil {
		return nil, err
	}

	return &OFPListener{ln}, err
}
