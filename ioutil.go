package of

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

func NewReader(w io.WriterTo) (io.Reader, error) {
	var buf bytes.Buffer

	_, err := w.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func Bytes(v interface{}) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, v)
	return buf.Bytes()
}
