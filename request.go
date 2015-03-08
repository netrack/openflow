package of

import (
	"bytes"
	"errors"
	"io"
	"math"
	"net"
)

const headerlen = 8

type Request struct {
	Header header
	Body   io.Reader

	Proto string
	Addr  net.Addr

	ContentLength int
}

// NewRequest returns a new Request given a type, address, and optional body
func NewRequest(t Type, body io.Reader) (*Request, error) {
	req := &Request{Body: body, Proto: "ofp1.3"}
	req.Header.Type = t

	//TODO: allow proto modification
	req.Header.Version = 0x4

	return req, nil
}

func (r *Request) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	r.Header.Length = uint16(headerlen)

	if r.Body == nil {
		return r.Header.WriteTo(w)
	}

	n, err = io.Copy(&buf, r.Body)
	if err != nil {
		return
	}

	if n > math.MaxUint16 {
		return 0, errors.New("openflow: body too long")
	}

	r.Header.Length += uint16(buf.Len())

	var wbuf bytes.Buffer
	_, err = r.Header.WriteTo(&wbuf)
	if err != nil {
		return
	}

	_, err = io.Copy(&wbuf, &buf)
	if err != nil {
		return
	}

	return wbuf.WriteTo(w)
}

func (r *Request) ReadFrom(rd io.Reader) (n int64, err error) {
	n, err = r.Header.ReadFrom(rd)
	if err != nil {
		return
	}

	var nn int
	contentlen := r.Header.Len() - headerlen

	buf := make([]byte, contentlen)
	nn, err = rd.Read(buf)
	n += int64(nn)

	if err != nil {
		return
	}

	r.Body = bytes.NewBuffer(buf)
	r.ContentLength = nn
	return
}
