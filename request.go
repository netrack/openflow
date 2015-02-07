package openflow

import (
	"bytes"
	"io"
)

const headerlen = 8

type Request struct {
	Header Header
	Body   io.Reader

	ContentLength int
}

func ReadRequest(r io.Reader) (*Request, error) {
	var header Header

	err := header.Read(r)
	if err != nil {
		return nil, err
	}

	contentlen := header.Len() - headerlen
	buf := make([]byte, contentlen)

	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	return &Request{header, bytes.NewBuffer(buf), n}, nil
}
