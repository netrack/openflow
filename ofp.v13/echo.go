package ofp

import (
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

// EchoRequest message consists of an OpenFlow header
// plus an arbitrary-length data field. The data field
// might be a message timestamp to check latency, various
// lengths to measure bandwidth, or zero-size to verify
// liveness between the switch and controller.
type EchoRequest struct {
	Data []byte
}

func (er *EchoRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, er.Data)
}

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

// EchoReply message consists of an OpenFlow header
// plus the unmodified data field of an echo request message.
type EchoReply struct {
	Data []byte
}

func (er *EchoReply) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, er.Data)
}

func (er *EchoReply) ReadFrom(r io.Reader) (n int64, err error) {
	er.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return int64(len(er.Data)), nil
}
