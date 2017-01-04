package openflow

import (
	"bytes"
	"testing"
)

// newHeader creates a new header with a given type and returns a wire
// representation of the created instance. Panics if error returns.
func newHeader(t Type) []byte {
	header := &Header{Version: 4, Type: t, Length: 8}

	var buf bytes.Buffer
	_, err := header.WriteTo(&buf)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func TestTransactionMatcher(t *testing.T) {
	header := new(Header)
	matcher := TransactionMatcher(header)

	r := NewRequest(TypePacketIn, nil)
	r.Header.Transaction = header.Transaction

	if !matcher.Match(r) {
		t.Errorf("Failed to match request of the same transaction")
	}

	r.Header.Transaction = header.Transaction + 1
	if matcher.Match(r) {
		t.Errorf("Matched request of different transaction")
	}
}
