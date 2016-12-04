package ofputil

import (
	"testing"

	"github.com/netrack/openflow/ofp"
	"github.com/netrack/openflow/ofptest"
	of "github.com/netrack/openflow"
)

func TestEchoHandler(t *testing.T) {
	ver := uint8(4)

	rw := ofptest.NewRecorder()
	h := HelloHandler(ver, nil)

	req := of.NewRequest(of.TypeHello, nil)
	req.Header.Version = 3
	req.Header.Transaction = 42

	h.Serve(rw, req)

	resp := rw.First()
	if resp.Header.Type != of.TypeHello {
		text := "hello message expected: %d"
		t.Errorf(text, resp.Header.Type)
	}

	if resp.Header.Version != ver {
		text := "unexpected version returned: %d"
		t.Errorf(text, resp.Header.Version)
	}

	if resp.Header.Transaction != req.Header.Transaction {
		text := "transaction identifier changed: %d"
		t.Errorf(text, resp.Header.Transaction)
	}
}

func TestHelloHandler(t *testing.T) {
	rw := ofptest.NewRecorder()
	h := EchoHandler(nil)

	echo := &ofp.EchoRequest{Data: []byte{1, 2, 3, 4}}
	req := of.NewRequest(of.TypeEchoReply, echo)
	req.Header.Transaction = 43

	h.Serve(rw, req)

	resp := rw.First()
	if resp.Header.Type != of.TypeEchoReply {
		text := "echo reply message expected: %d"
		t.Errorf(text, resp.Header.Type)
	}

	if resp.Header.Transaction != req.Header.Transaction {
		text := "transaction identifier changed: %d"
		t.Errorf(text, resp.Header.Transaction)
	}
}
