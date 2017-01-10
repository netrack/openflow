package openflow

import (
	"bytes"
	"testing"
)

func TestCopyReaderWriteTo(t *testing.T) {
	var buf bytes.Buffer
	rd := copyReader{WriterTo: nil}

	// Ensure the nil writer is acceptable.
	n, err := rd.WriteTo(&buf)
	if n != 0 || err != nil {
		t.Errorf("Expected successful write to buffer: %s", err)
	}
}

func TestCopyReaderRead(t *testing.T) {
	// Put some data into the writer and ensure all bits
	// have been returned during the Read call.
	bits := []byte{1, 2, 3, 4}
	buf := bytes.NewBuffer(bits)
	rd := copyReader{WriterTo: buf}

	p := make([]byte, len(bits))
	nn, err := rd.Read(p)

	if nn != len(bits) || err != nil {
		t.Errorf("Expected successful read from reader: %s", err)
	}
}

func TestReadRequest(t *testing.T) {
	var req Request
	var buf bytes.Buffer

	h := Header{4, TypeHello, 8, 0}
	h.WriteTo(&buf)

	_, err := req.ReadFrom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if req.ContentLength != 0 {
		t.Fatal("Wrong content length:", req.ContentLength)
	}

	if req.Header.Type != TypeHello {
		t.Fatal("Wrong header type:", req.Header.Type)
	}

	if req.Header.Version != 4 {
		t.Fatal("Wrong header version:", req.Header.Version)
	}
}
