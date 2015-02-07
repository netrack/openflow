package openflow

import (
	"bytes"
	"testing"
)

func TestReadRequest(t *testing.T) {
	var buf bytes.Buffer

	h := Header{4, T_HELLO, 8, 0}
	h.Write(&buf)

	req, err := ReadRequest(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if req.ContentLength != 0 {
		t.Fatal("Wrong content length:", req.ContentLength)
	}

	if req.Header.Type != T_HELLO {
		t.Fatal("Wrong header type:", req.Header.Type)
	}

	if req.Header.Version != 4 {
		t.Fatal("Wrong header version:", req.Header.Version)
	}
}
