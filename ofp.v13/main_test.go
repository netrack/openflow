package ofp

import (
	"bytes"
	"encoding/gob"
	"io"
	"testing"
)

// ReadFromWriterTo describes the types that implement
// the interfaces of io.ReaderFrom and io.WriterTo.
type ReadFromWriterTo interface {
	io.ReaderFrom
	io.WriterTo
}

// testM defines the marshaling testing type.
type testM struct {
	w io.WriterTo
	b []byte
}

// testU defines the unmarshaling testing type.
type testU struct {
	r io.ReaderFrom
	b []byte
}

// testMU defines the marshaling/unmarshaling
// testing type.
type testMU struct {
	rw ReadFromWriterTo
	b  []byte
}

// testMarshal validates, that each passed writer produces
// exact sequence of bytes, as specified.
func testMarshal(t *testing.T, tests []testM) {
	for _, test := range tests {
		var buf bytes.Buffer
		nn, err := test.w.WriteTo(&buf)

		if err != nil {
			t.Fatalf("Failed to marshal the given packet: "+
				"`%x`, got error: %s", test.b, err)
		}

		if nn != int64(len(test.b)) {
			t.Fatalf("Invalid length returned on attempt to "+
				"marshal: `%x`: %d, expected %d",
				test.b, nn, len(test.b))
		}

		if bytes.Compare(test.b, buf.Bytes()) != 0 {
			t.Fatalf("The marshaled result is not equal to "+
				"the\nexpected:\n`%x`,\ngot instead:\n`%x`",
				test.b, buf.Bytes())
		}
	}
}

// testUnmarshal validates each passed reader produces
// exact objects, as expected.
func testUnmarshal(t *testing.T, tests []testU) {
	for _, test := range tests {
		var before bytes.Buffer
		enc := gob.NewEncoder(&before)

		err := enc.Encode(test.r)
		if err != nil {
			t.Fatalf("Failed to encode Go object: `%v`: %s",
				test.r, err)
		}

		buf := bytes.NewBuffer(test.b)
		nn, err := test.r.ReadFrom(buf)

		if err != nil {
			t.Fatalf("Failed to unmarshal the given packet: "+
				"`%x`, got error: %s", test.b, err)
		}

		if nn != int64(len(test.b)) {
			t.Fatalf("Invalid length returned on attempt to "+
				"marshal: `%x`: %d, expected %d",
				test.b, nn, len(test.b))
		}

		var after bytes.Buffer
		enc = gob.NewEncoder(&after)

		err = enc.Encode(test.r)
		if err != nil {
			t.Fatalf("Failed to encode Go object: `%v`: %s",
				test.r, err)
		}

		if bytes.Compare(before.Bytes(), after.Bytes()) != 0 {
			t.Fatalf("The unmarshaled result is not equal to "+
				"the expected one:\n`%v`,\ngot instead:\n`%v`\n%v",
				before.Bytes(), after.Bytes(), test.r)
		}
	}
}

// testMarshalUnmarshal executes the marshaling and unmarshaling
// test for the given sequence of tests.
func testMarshalUnmarshal(t *testing.T, tests []testMU) {
	for _, test := range tests {
		testMarshal(t, []testM{{test.rw, test.b}})
		testUnmarshal(t, []testU{{test.rw, test.b}})
	}
}
