package encodingtest

import (
	"bytes"
	"encoding/gob"
	"io"
	"testing"
)

// M defines the marshaling testing type.
type M struct {
	Writer io.WriterTo
	Bytes  []byte
}

// RunM validates, that each passed marshaler produces
// exact sequence of bytes, as specified.
func RunM(t *testing.T, tests []M) {
	for _, test := range tests {
		var buf bytes.Buffer
		nn, err := test.Writer.WriteTo(&buf)

		if err != nil {
			t.Fatalf("Failed to marshal the given packet: "+
				"`%x`, got error: %s", test.Bytes, err)
		}

		if nn != int64(len(test.Bytes)) {
			t.Fatalf("Invalid length returned on attempt to "+
				"marshal:\n`%x`: %d,\nexpected:\n`%x`: %d\n",
				buf.Bytes(), nn, test.Bytes, len(test.Bytes))
		}

		if bytes.Compare(test.Bytes, buf.Bytes()) != 0 {
			t.Fatalf("The marshaled result is not equal to "+
				"the\nexpected:\n`%x`,\ngot instead:\n`%x`",
				test.Bytes, buf.Bytes())
		}
	}
}

// U defines the unmarshaling testing type.
type U struct {
	Reader io.ReaderFrom
	Bytes  []byte
}

// RunU validates each passed reader produces
// exact objects, as expected.
func RunU(t *testing.T, tests []U) {
	for _, test := range tests {
		var before bytes.Buffer
		enc := gob.NewEncoder(&before)

		err := enc.Encode(test.Reader)
		if err != nil {
			t.Fatalf("Failed to encode Go object: `%v`: %s",
				test.Reader, err)
		}

		buf := bytes.NewBuffer(test.Bytes)
		nn, err := test.Reader.ReadFrom(buf)

		if err != nil {
			t.Fatalf("Failed to unmarshal the given packet: "+
				"`%x`, got error: %s", test.Bytes, err)
		}

		if nn != int64(len(test.Bytes)) {
			t.Fatalf("Invalid length returned on attempt to "+
				"marshal: `%x`: %d, expected %d",
				test.Bytes, nn, len(test.Bytes))
		}

		var after bytes.Buffer
		enc = gob.NewEncoder(&after)

		err = enc.Encode(test.Reader)
		if err != nil {
			t.Fatalf("Failed to encode Go object: `%v`: %s",
				test.Reader, err)
		}

		if bytes.Compare(before.Bytes(), after.Bytes()) != 0 {
			t.Fatalf("The unmarshaled result is not equal to "+
				"the expected one:\n`%x`,\ngot instead:\n`%x`\n%v",
				before.Bytes(), after.Bytes(), test.Reader)
		}
	}
}

// MU defines the marshaling/unmarshaling
// testing type.
type MU struct {
	ReadWriter interface {
		io.ReaderFrom
		io.WriterTo
	}

	Bytes []byte
}

// RunMU executes the marshaling and unmarshaling
// test for the given sequence of tests.
func RunMU(t *testing.T, tests []MU) {
	for _, test := range tests {
		RunM(t, []M{{test.ReadWriter, test.Bytes}})
		RunU(t, []U{{test.ReadWriter, test.Bytes}})
	}
}
