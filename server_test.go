package openflow

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	//"github.com/netrack/net/pkg"
	//ofp "github.com/netrack/openflow/ofp13"
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

func TestServerMux(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := NewServeMux()
	mux.HandleFunc(T_HELLO, func(rw ResponseWriter, r *Request) {
		rw.Write([]byte("response"))
		wg.Done()
	})

	reader := bytes.NewBuffer([]byte{4, 0, 0, 8, 0, 0, 0, 0})
	conn := &dummyConn{r: *reader}

	s := Server{Addr: "0.0.0.0:6633", Handler: mux}
	err := s.Serve(&dummyListener{conn})

	if err != io.EOF {
		t.Fatal("Serve failed:", err)
	}

	wg.Wait()

	if conn.w.String() != "response" {
		t.Fatal("Invalid data returned:", conn.w.String())
	}
}

/*
 *func TestServer(t *testing.T) {
 *    handler := HandlerFunc(func(rw ResponseWriter, r *Request) {
 *        switch r.Header.Type {
 *        case T_HELLO:
 *            fmt.Println("GOT HELLO:", r.Header)
 *            hello := Header{r.Header.Version, T_HELLO, 8, r.Header.Xid}
 *            hello.Write(rw)
 *        case T_PACKET_IN:
 *            var pin ofp.PacketIn
 *            var eth pkg.EthernetII
 *
 *            err1 := pin.Read(r.Body)
 *            err2 := eth.Read(r.Body)
 *            fmt.Println(pin)
 *            fmt.Println("GOT PACKET_IN:", err1, err2, eth, eth.HWDst)
 *
 *            pout := &ofp.PacketOut{
 *                BufferID: pin.BufferID,
 *                Actions: []interface {
 *                    Write(io.Writer) error
 *                }{
 *                    ofp.ActionOutput{ofp.P_FLOOD, 0},
 *                },
 *            }
 *
 *            var buf bytes.Buffer
 *            err := pout.Write(&buf)
 *
 *            header := Header{r.Header.Version, T_PACKET_OUT, 8 + uint16(buf.Len()), r.Header.Xid}
 *            header.Write(rw)
 *            rw.Write(buf.Bytes())
 *            fmt.Println("PACKET_OUT:", err, buf.Bytes())
 *        case T_ERROR:
 *            fmt.Println("GOT ERROR:", r.Header)
 *        }
 *    })
 *
 *    s := Server{Addr: "0.0.0.0:6633", Handler: handler}
 *    s.ListenAndServe()
 *}
 */
