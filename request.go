package openflow

import (
	"bytes"
	"io"
)

const headerlen = 8

type Request struct {
	Header header
	Body   io.Reader

	ContentLength int
}

func (req *Request) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = req.Header.ReadFrom(r)
	if err != nil {
		return
	}

	var nn int
	contentlen := req.Header.Len() - headerlen

	buf := make([]byte, contentlen)
	nn, err = r.Read(buf)
	n += int64(nn)

	if err != nil {
		return
	}

	req.Body = bytes.NewBuffer(buf)
	req.ContentLength = nn
	return
}
