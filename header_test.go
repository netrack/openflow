package of

import (
	"testing"
)

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
