package ofp

import (
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// EchoRequest is a message with arbitrary-length data field.
//
// For example, to create a request with a timestamp to check the
// latency between the switch and controller:
//
//	now, _ := time.Now().MarshalBinary()
//	req := of.NewRequest(of.TypeEchoRequest, &EchoRequest{now})
//	...
type EchoRequest struct {
	// The arbitrary-length chunk of bytes.
	Data []byte
}

// WriteTo implements io.WriterTo interface. It serializes the echo
// request message into the wire format.
func (er *EchoRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, er.Data)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the echo
// request message from the wire format.
//
// The passed reader instance have to be limited up to the length of message
// body.
func (er *EchoRequest) ReadFrom(r io.Reader) (n int64, err error) {
	// The header of the echo request will be unmarshalled by the
	// previous read, so we need to read up to the end of the
	// buffer. So it is responsibility of the called to provide
	// fixed-length reader.
	er.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return int64(len(er.Data)), nil
}

// EchoReply is a message with unmodified data field of an echo request
// message.
//
// The message is used to respond on echo-requests submitted by the other
// side of communication channel. For example, to send an empty echo-reply
// message, you could create the following request:
//
//	req := of.NewRequest(of.TypeEchoReply, &EchoRequest{data})
//	...
type EchoReply struct {
	// The data copied from the received echo-request.
	Data []byte
}

// WriteTo implements io.WriterTo interface. It serializes the echo request
// to the wire format.
func (er *EchoReply) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, er.Data)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the echo
// request from the wire format.
func (er *EchoReply) ReadFrom(r io.Reader) (n int64, err error) {
	er.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return int64(len(er.Data)), nil
}
