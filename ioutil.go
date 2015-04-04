package of

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

type ReaderFunc func([]byte) (int, error)

func (fn ReaderFunc) Read(b []byte) (int, error) {
	return fn(b)
}

func NewReader(w ...io.WriterTo) io.Reader {
	var buf bytes.Buffer
	var err error

	for _, wt := range w {
		if _, err = wt.WriteTo(&buf); err != nil {
			break
		}
	}

	return ReaderFunc(func(b []byte) (int, error) {
		if err != nil {
			return 0, err
		}

		return buf.Read(b)
	})
}

func Bytes(v interface{}) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, v)
	return buf.Bytes()
}

func WriteAllTo(w io.Writer, wt ...io.WriterTo) (int64, error) {
	var n int64

	for _, wrt := range wt {
		nn, err := wrt.WriteTo(w)
		if n += nn; err != nil {
			return n, err
		}
	}

	return n, nil
}

func ReadAllFrom(r io.Reader, rf ...io.ReaderFrom) (int64, error) {
	var n int64

	for _, rfr := range rf {
		nn, err := rfr.ReadFrom(r)
		if n += nn; err != nil {
			return n, err
		}
	}

	return n, nil
}
