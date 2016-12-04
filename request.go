package openflow

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

// A Request represents an OpenFlow request received by the server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server usage.
type Request struct {
	// Header contains the request header fields either received by
	// the server or sent by the client.
	Header Header

	// Body is the request's body. For client requests a nil
	// body means the request has no body, such as a echo requests.
	//
	// For server requests the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	Body io.Reader

	// The protocol version for incoming requests.
	// Client requests always use OFP/1.3.
	Proto      string
	ProtoMajor int // 1
	ProtoMinor int // 3

	Addr net.Addr

	// ContentLength records the length of the associated content.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64

	// Connection instance.
	conn Conn
}

// NewRequest returns a new Request given a type, address, and optional
// body.
func NewRequest(t Type, body io.WriterTo) *Request {
	req := &Request{
		Body:       newReader(body),
		Proto:      "OFP/1.3",
		ProtoMajor: 1, ProtoMinor: 3,
	}

	req.Header.Version = uint8(req.ProtoMajor + req.ProtoMinor)
	req.Header.Type = t

	return req
}

// ProtoAtLeast reports whether the OpenFlow protocol used in the request
// is at least major.minor.
func (r *Request) ProtoAtLeast(major, minor int) bool {
	return r.ProtoMajor > major ||
		r.ProtoMajor == major && r.ProtoMinor >= minor
}

// Conn returns the instance of the OpenFlow protocol connection.
func (r *Request) Conn() Conn {
	return r.conn
}

// WriteTo implements WriterTo interface. Writes the request in wire format
// to w until there's no more data to write or when an error occurs.
func (r *Request) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	// Save the header length into the request.
	r.Header.Length = uint16(headerlen)

	// If the body of the request is not specified (like in case when
	// the echo requests or reply are used), keep the fast path.
	if r.Body == nil {
		return r.Header.WriteTo(w)
	}

	// Previous to the request serialization we have to specify the
	// length of the body, so there is no choice unless copy data
	// from the reader to buffer.
	n, err = io.Copy(&buf, r.Body)
	if err != nil {
		return
	}

	// For sure we need to double check that body length fits into
	// the header length.
	if n+headerlen > math.MaxUint16 {
		return 0, ErrBodyTooLong
	}

	r.Header.Length += uint16(buf.Len())

	var wbuf bytes.Buffer
	// Write the header of the OpenFlow packet first, and then the
	// body should be written accordingly to the buffer.
	_, err = r.Header.WriteTo(&wbuf)
	if err != nil {
		return
	}

	_, err = io.Copy(&wbuf, &buf)
	if err != nil {
		return
	}

	// At the end, write the buffer to the specified writer instance.
	return wbuf.WriteTo(w)
}

// ReadFrom implements ReaderFrom interface. Reads the request in wire
// format from the r to the Request structure.
func (r *Request) ReadFrom(rd io.Reader) (n int64, err error) {
	// On the first step we will read the OpenFlow header, so
	// we could get the total length of the OpenFlow message.
	n, err = r.Header.ReadFrom(rd)
	if err != nil {
		return
	}

	// Decode the protocol version major and minor number to
	// make the request interface more or less friendly.
	r.ProtoMajor = 1
	r.ProtoMinor = int(r.Header.Version - 1)

	// FIXME: wrong for version 2
	r.Proto = fmt.Sprintf("OFP/1.%d", r.ProtoMinor)

	var nn int
	contentlen := r.Header.Len() - headerlen

	// Define a buffer to fit the content of the OpenFlow package.
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
