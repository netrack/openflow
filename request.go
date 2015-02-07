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

func (req *Request) Read(r io.Reader) error {
	err := req.Header.Read(r)
	if err != nil {
		return err
	}

	contentlen := req.Header.Len() - headerlen

	buf := make([]byte, contentlen)

	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	req.Body = bytes.NewBuffer(buf)
	req.ContentLength = n
	return nil
}
