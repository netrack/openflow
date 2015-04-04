package of

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
)

var (
	ErrUnknownVersion = errors.New("openflow: Unknown version passed")
	ErrBodyTooLong    = errors.New("openflow: Request body is too long")
)

const headerlen = 8

type Request struct {
	Header header

	// Body is the request's body. For client requests a nil
	// body means the request has no body, such as a echo requests.
	//
	// For server requests the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	Body io.Reader

	// The protocol version for incoming requests.
	// Client requests always use OFP/1.3.
	Proto      string
	ProtoMajor int
	ProtoMinor int

	Addr net.Addr

	// ContentLength records the length of the associated content.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64
}

// NewRequest returns a new Request given a type, address, and optional body
func NewRequest(t Type, body io.Reader) (*Request, error) {
	req := &Request{Body: body, Proto: "OFP/1.3", ProtoMajor: 1, ProtoMinor: 3}

	req.Header.Version = uint8(req.ProtoMajor + req.ProtoMinor)
	req.Header.Type = t

	return req, nil
}

// ProtoAtLeast reports whether the OpenFlow protocol used
// in the request is at least major.minor.
func (r *Request) ProtoAtLeast(major, minor int) bool {
	return r.ProtoMajor > major ||
		r.ProtoMajor == major && r.ProtoMinor >= minor
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
		return 0, ErrBodyTooLong
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

	r.ProtoMajor = 1
	r.ProtoMinor = int(r.Header.Version - 1)

	//FIXME: wrong for version 2
	r.Proto = fmt.Sprintf("OFP/1.%d", r.ProtoMinor)

	var nn int
	contentlen := r.Header.Len() - headerlen

	buf := make([]byte, contentlen)
	nn, err = rd.Read(buf)
	n += int64(nn)

	if err != nil {
		return
	}

	r.Body = bytes.NewBuffer(buf)
	r.ContentLength = int64(nn)
	return
}
