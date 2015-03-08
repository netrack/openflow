package openflow

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

type dummyAddr string

func (a dummyAddr) Network() string {
	return "dummy"
}

func (a dummyAddr) String() string {
	return string(a)
}

type dummyConn struct {
	r bytes.Buffer
	w bytes.Buffer
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
	return dummyAddr("local-addr")
}

func (c *dummyConn) RemoteAddr() net.Addr {
	return dummyAddr("remote-addr")
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
	ln, err := Listen("tcp6", "[::1]:6633")
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
		t.Fatal("Faile to create OFP connection")
	}
}

func TestDial(t *testing.T) {
	ln, err := Listen("tcp6", "[::1]:6633")
	if err != nil {
		t.Fatal("Failed to create listener:", err)
	}

	defer ln.Close()

	rwc, err := Dial("tcp6", "[::1]:6633")
	if err != nil {
		t.Fatal("Failed to dial listener:", err)
	}

	defer rwc.Close()

	cconn, err := ln.AcceptOFP()
	if err != nil {
		t.Fatalf("Failed to accept client connection:", err)
	}

	defer cconn.Close()

	r, err := NewRequest(T_HELLO, nil)
	if err != nil {
		t.Fatal("Failed to create a new request:", err)
	}

	err = rwc.Send(r)
	if err != nil {
		t.Fatal("Failed to send request:", err)
	}

	r, err = cconn.Recv()
	if err != nil {
		t.Fatal("Failed to receive request:", err)
	}

	if r.Addr.String() != rwc.LocalAddr().String() {
		t.Fatal("Wrong address returned:", r.Addr.String())
	}

	if r.Header.Type != T_HELLO {
		t.Fatal("Wrong message type returned:", r.Header.Type)
	}

	if r.Header.Length != 8 {
		t.Fatal("Wrong length returned:", r.Header.Length)
	}

	if r.ContentLength != 0 {
		t.Fatal("Wrong content length returned:", r.ContentLength)
	}
}
