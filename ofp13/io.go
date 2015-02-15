package ofp13

import (
	"bytes"
	"io"
)

func Bytes(w io.WriterTo) []byte {
	var buf bytes.Buffer
	w.WriteTo(&buf)
	return buf.Bytes()
}
