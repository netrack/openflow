package ofp

import (
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding/binary"
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
	return binary.Write(w, binary.BigEndian, er.Data)
}

func (er *EchoRequest) ReadFrom(r io.Reader) (n int64, err error) {
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
	return binary.Write(w, binary.BigEndian, er.Data)
}

func (er *EchoReply) ReadFrom(r io.Reader) (n int64, err error) {
	er.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return int64(len(er.Data)), nil
}
