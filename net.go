package of

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type OFPConn interface {
	net.Conn

	Hijack() (net.Conn, *bufio.ReadWriter, error)

	// Receive receives message from input buffer
	Receive() (*Request, error)

	// Send writes message to output buffer
	Send(*Request) error

	Flush() error
}

type Conn struct {
	rwc net.Conn
	buf *bufio.ReadWriter

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	mu sync.Mutex

	hijackedv bool
}

func NewConn(conn net.Conn) *Conn {
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	brw := bufio.NewReadWriter(br, bw)
	return &Conn{rwc: conn, buf: brw}
}

func (c *Conn) hijacked() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hijackedv
}

func (c *Conn) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hijackedv {
		return nil, nil, ErrHijacked
	}

	c.hijackedv = true
	rwc := c.rwc
	buf := c.buf

	c.rwc = nil
	c.buf = nil
	return rwc, buf, nil
}

func (c *Conn) Read(b []byte) (int, error) {
	if c.hijacked() {
		return 0, ErrHijacked
	}

	return c.buf.Read(b)
}

func (c *Conn) Receive() (*Request, error) {
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

func (c *Conn) Write(b []byte) (int, error) {
	if c.hijacked() {
		return 0, ErrHijacked
	}

	return c.buf.Write(b)
}

func (c *Conn) Flush() error {
	return c.buf.Flush()
}

func (c *Conn) Send(r *Request) error {
	if c.hijacked() {
		return ErrHijacked
	}

	if d := c.WriteTimeout; d != 0 {
		defer func() {
			c.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	_, err := r.WriteTo(c.buf)
	if err != nil {
		return err
	}

	return c.buf.Flush()
}

func (c *Conn) Close() error {
	return c.rwc.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.rwc.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.rwc.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.rwc.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.rwc.SetWriteDeadline(t)
}

func Dial(network, addr string) (OFPConn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

func DialTLS(network, addr string) (OFPConn, error) {
	panic("PANIC PANIC PANIC")
}

type Listener struct {
	ln net.Listener
}

func (l *Listener) Accept() (net.Conn, error) {
	return l.AcceptOFP()
}

func (l *Listener) AcceptOFP() (OFPConn, error) {
	conn, err := l.ln.Accept()
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

func (l *Listener) Close() error {
	return l.ln.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.ln.Addr()
}

func Listen(network, laddr string) (*Listener, error) {
	tcpaddr, err := net.ResolveTCPAddr(network, laddr)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP(network, tcpaddr)
	if err != nil {
		return nil, err
	}

	return &Listener{ln}, err
}
