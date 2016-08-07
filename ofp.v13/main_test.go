package ofp

import (
	"bytes"
	"io"
	"testing"
)

type testValue struct {
	w io.WriterTo
	b []byte
}

// testMarshal validate, that each passed writer produces
// exact sequence of bytes, as specified.
func testMarshal(t *testing.T, tests []testValue) {
	for _, test := range tests {
		var buf bytes.Buffer
		nn, err := test.w.WriteTo(&buf)

		if err != nil {
			t.Fatalf("Failed to marshal the given packet:"+
				"%x, got error: %s", test.b, err)
		}

		if nn != int64(len(test.b)) {
			t.Fatalf("Invalid length returned on attempt to "+
				"marshal: `%x`: %d, expected %d",
				test.b, nn, len(test.b))
		}

		if bytes.Compare(test.b, buf.Bytes()) != 0 {
			t.Fatalf("The marshaled result is not equal to "+
				"the expected one: `%x`, got instead: `%x`",
				test.b, buf.Bytes())
		}
	}
}
