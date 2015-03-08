package openflow

import (
	"bytes"
	"io"
)

func NewReader(w io.WriterTo) (io.Reader, error) {
	var buf bytes.Buffer

	_, err := w.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
