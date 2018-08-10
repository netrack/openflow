package openflow

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

type dummyConn struct {
	r bytes.Buffer
	w bytes.Buffer

	lAddr string
	rAddr string

	closed bool
}

func (c *dummyConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func (c *dummyConn) Write(b []byte) (int, error) {
	return c.w.Write(b)
}

// Close implement net.Conn interface.
func (c *dummyConn) Close() error {
	c.closed = true
	return nil
}

func (c *dummyConn) LocalAddr() net.Addr {
	return dummyAddr(c.lAddr)
}

func (c *dummyConn) RemoteAddr() net.Addr {
	return dummyAddr(c.rAddr)
}

func (c *dummyConn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *dummyConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *dummyConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

// Type dummyBlockConn defines the implementation of net.Conn interface
// used for testing. The Read and Write methods blocks until the conn
// is closed.
type dummyBlockConn struct {
	dummyConn

	// Cancel is a channel used to simulate the blocking read and
	// write operations from the connection. It will be closed on
	// close of the client connection, thus release all blocked ops.
	cancel     chan struct{}
	cancelOnce sync.Once

	once sync.Once
}

// init initializes the attributes of the blocking connection struct.
func (c *dummyBlockConn) init() {
	c.cancel = make(chan struct{})
}

// Read implements net.Conn interface. It blocks until connection will
// be closed and returns an error then.
func (c *dummyBlockConn) Read(b []byte) (int, error) {
	c.once.Do(c.init)
	_, _ = <-c.cancel
	return 0, errors.New("conn: connection closed")
}

// Write implement net.Conn interface. It block until connection will
// be closed and returns an error then.
func (c *dummyBlockConn) Write(b []byte) (int, error) {
	c.once.Do(c.init)
	_, _ = <-c.cancel
	return 0, errors.New("conn: connection closed")
}

// Close implements net.Conn interface. It releases the blocked Read
// and Write methods by closing the internal channel.
func (c *dummyBlockConn) Close() error {
	c.once.Do(c.init)

	// Cancel the close channel only once to prevent panics.
	c.cancelOnce.Do(func() { close(c.cancel) })
	return c.dummyConn.Close()
}

// The dummyListerner defines the mock implementation of the
// net.Listener. It sequentially returns a single connection from the
// list of the connection.
type dummyListener struct {
	conns []net.Conn
}

func (l *dummyListener) Accept() (net.Conn, error) {
	var conn net.Conn
	if len(l.conns) == 0 {
		return conn, io.EOF
	}

	conn, l.conns = l.conns[0], l.conns[1:]
	return conn, nil
}

func (l *dummyListener) Close() error {
	return nil
}

func (l *dummyListener) Addr() net.Addr {
	return dummyAddr("dummy-address")
}

func TestListener(t *testing.T) {
	ln, err := Listen("tcp", ":6666")
	if err != nil {
		t.Fatal("Failed to create listener:", err)
	}

	ofpLn := ln.(*listener)
	defer ofpLn.ln.Close()

	dconn := &dummyConn{}
	dln := &dummyListener{[]net.Conn{dconn}}

	ofpLn.ln = dln

	c, err := ln.Accept()
	ofpConn := c.(*conn)

	if err != nil {
		t.Fatal("Failed to accept a new connection:", err)
	}

	if ofpConn.rwc != dconn {
		t.Fatal("Failed to create OFP connection")
	}
}

func TestDial(t *testing.T) {
	ln, err := Listen("tcp", "localhost:6667")
	if err != nil {
		t.Fatal("Failed to create listener:", err)
	}

	// Defer the connection closing call, since we are
	// going to replace it with a dummy one.
	ofpLn := ln.(*listener)
	defer ofpLn.ln.Close()

	// Define a connection instance that is expecting
	// to be returned on accepting the client connection.
	serverAddr := ln.Addr().String()
	serverConn := &dummyConn{rAddr: serverAddr}

	// Put into the read buffer of the defined connection
	// a single OpenFlow Hello request, so we could ensure
	// that data is read from the connection correctly.
	r := NewRequest(TypeHello, nil)
	r.WriteTo(&serverConn.r)

	// Perform the actual connection replacement.
	ofpLn.ln = &dummyListener{[]net.Conn{serverConn}}

	rwc, err := Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("Failed to dial listener:", err)
	}

	defer rwc.Close()

	// Replace the client connection with a dummy one,
	// so we could perform damn simple unit test.
	clientConn := &dummyConn{}
	rwc.(*conn).rwc = clientConn

	cconn, err := ln.Accept()
	if err != nil {
		t.Fatal("Failed to accept client connection:", err)
	}

	// Define a new OpenFlow Hello message and send it into
	// the client connection.
	r = NewRequest(TypeHello, nil)
	err = rwc.Send(r)
	if err != nil {
		t.Fatal("Failed to send request:", err)
	}

	r, err = cconn.Receive()
	if err != nil {
		t.Fatal("Failed to receive request:", err)
	}

	if r.Addr.String() != serverAddr {
		t.Fatal("Wrong address returned:", r.Addr.String())
	}

	// Validate attributes of the retrieved OpenFlow request.
	// At first, it certainly should be a Hello message.
	if r.Header.Type != TypeHello {
		t.Fatal("Wrong message type returned:", r.Header.Type)
	}

	// On the other hand nothing additional should be presented
	// in the request apart of required fields.
	if r.Header.Length != 8 {
		t.Fatal("Wrong length returned:", r.Header.Length)
	}

	// No content expected inside the request packet.
	if r.ContentLength != 0 {
		t.Fatal("Wrong content length returned:", r.ContentLength)
	}
}
