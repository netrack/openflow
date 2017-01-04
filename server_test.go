package openflow

import (
	"io"
	"net"
	"testing"
)

func TestReceiver(t *testing.T) {
	dconn := new(dummyConn)
	dconn.r.Write(newHeader(TypeHello))

	rcvr := &receiver{Conn: newConn(dconn)}
	wrap := <-rcvr.C()

	if wrap.err != nil {
		t.Fatalf("Error returned from receiver: %s", wrap.err)
	}

	htype := wrap.req.Header.Type
	if htype != TypeHello {
		t.Fatalf("Incorrect request type returned: %d", htype)
	}

	// Receive of another request should fail with EOF.
	wrap = <-rcvr.C()
	if wrap.err != io.EOF {
		t.Fatalf("Expected EOF on second receive: %s", wrap.err)
	}

	_, ok := <-rcvr.C()
	if ok {
		t.Fatalf("Channel should be closed on error")
	}
}

func TestServeMaxConns(t *testing.T) {
	s := Server{MaxConns: 2}
	defer s.close()

	// Create three connection, with a maximum connections set to two.
	// This means, a third client should be closed by the server.
	dconn1 := new(dummyBlockConn)
	dconn2 := new(dummyBlockConn)
	dconn3 := new(dummyBlockConn)

	defer dconn1.Close()
	defer dconn2.Close()

	dln := &dummyListener{[]net.Conn{dconn1, dconn2, dconn3}}
	err := s.Serve(dln)

	// The mock of the listener returns an error, when connections
	// have been extracted from the queue completely.
	if err != io.EOF {
		t.Errorf("Serving of the listener failed: %s", err)
	}

	if dconn1.closed || dconn2.closed {
		t.Errorf("Two first connection expected to be alive")
	}

	if !dconn3.closed {
		t.Errorf("Third connection expected to be closed")
	}
}

func TestServerServe(t *testing.T) {
}
