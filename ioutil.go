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

func NewReader(w io.WriterTo) io.Reader {
	var buf bytes.Buffer
	_, err := w.WriteTo(&buf)

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
