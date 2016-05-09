package of

import (
	"bytes"
	"io"
	"net"
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
}

func (c *dummyConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func (c *dummyConn) Write(b []byte) (int, error) {
	return c.w.Write(b)
}

func (c *dummyConn) Close() error {
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

type dummyListener struct {
	conn net.Conn
}

func (l *dummyListener) Accept() (c net.Conn, e error) {
	c, l.conn = l.conn, nil
	if c == nil {
		e = io.EOF
	}

	return
}

func (l *dummyListener) Close() error {
	if l.conn != nil {
		return l.conn.Close()
	}

	return nil
}

func (l *dummyListener) Addr() net.Addr {
	return dummyAddr("dummy-address")
}

func TestListener(t *testing.T) {
	ln, err := Listen("tcp6", "[::1]:0")
	if err != nil {
		t.Fatal("Failed to create listener:", err)
	}

	defer ln.ln.Close()

	dconn := &dummyConn{}
	dln := &dummyListener{dconn}

	ln.ln = dln

	conn, err := ln.AcceptOFP()
	if err != nil {
		t.Fatal("Failed to accept a new connection:", err)
	}

	if conn.rwc != dconn {
		t.Fatal("Failed to create OFP connection")
	}
}

func TestDial(t *testing.T) {
	ln, err := Listen("tcp6", "[::1]:0")
	if err != nil {
		t.Fatal("Failed to create listener:", err)
	}

	// Defer the connection closing call, since we are
	// going to replace it with a dummy one.
	defer ln.ln.Close()

	// Define a connection instance that is expecting
	// to be returned on accepting the client connection.
	serverAddr := ln.Addr().String()
	serverConn := &dummyConn{rAddr: serverAddr}

	// Put into the read buffer of the defined connection
	// a single OpenFlow Hello request, so we could ensure
	// that data is read from the connection correctly.
	r, _ := NewRequest(TypeHello, nil)
	r.WriteTo(&serverConn.r)

	// Perform the actual connection replacement.
	ln.ln = &dummyListener{serverConn}

	rwc, err := Dial("tcp6", serverAddr)
	if err != nil {
		t.Fatal("Failed to dial listener:", err)
	}

	defer rwc.Close()

	// Replace the client connection with a dummy one,
	// so we could perform damn simple unit test.
	clientConn := &dummyConn{}
	rwc.(*OFPConn).rwc = clientConn

	cconn, err := ln.AcceptOFP()
	if err != nil {
		t.Fatalf("Failed to accept client connection:", err)
	}

	// Define a new OpenFlow Hello message and send it into
	// the client connection.
	r, _ = NewRequest(TypeHello, nil)
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
