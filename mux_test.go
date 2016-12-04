package openflow

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

func TestMultiMatcher(t *testing.T) {
	txn := uint32(42)

	// A function, that matches the type of the request.
	mf1 := func(r *Request) bool {
		return r.Header.Type == TypeHello
	}

	// A function, that matches transaction ID.
	mf2 := func(r *Request) bool {
		return r.Header.Transaction == txn
	}

	matcher := MultiMatcher(&MatcherFunc{mf1}, &MatcherFunc{mf2})

	r := NewRequest(TypePacketIn, nil)
	if matcher.Match(r) {
		t.Errorf("Matched request with different type")
	}

	r = NewRequest(TypeHello, nil)
	r.Header.Transaction = txn + 1

	if matcher.Match(r) {
		t.Errorf("Matched request with different transaction ID")
	}

	r.Header.Transaction = txn
	if !matcher.Match(r) {
		t.Errorf("Request supposed to match")
	}
}

func TestTypeMux(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := NewTypeMux()
	mux.HandleFunc(TypeHello, func(rw ResponseWriter, r *Request) {
		defer wg.Done()

		wbuf := bytes.NewBuffer([]byte{0, 0, 0, 0})
		rw.Write(r.Header.Copy(), wbuf)
	})

	mux.HandleFunc(TypeEchoRequest, func(rw ResponseWriter, r *Request) {
		t.Errorf("This handler should never be called")
	})

	reader := bytes.NewBuffer([]byte{4, 0, 0, 8, 0, 0, 0, 0})
	conn := &dummyConn{r: *reader}

	s := Server{Addr: "0.0.0.0:6633", Handler: mux}
	err := s.Serve(&dummyListener{conn})

	// Serve function will treat the connection as a regular
	// connection, thus will try to read the next message after
	// the reading the first one. And as the buffer is empty
	// it will return EOF, which will be used to identify the
	// successful read of the message.
	if err != io.EOF {
		t.Errorf("Serve failed:", err)
	}

	wg.Wait()

	returned := fmt.Sprintf("%x", conn.w.Bytes())
	if returned != "0400000c0000000000000000" {
		t.Errorf("Invalid data returned: ", returned)
	}
}
