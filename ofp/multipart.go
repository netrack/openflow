package ofp

import (
	"bytes"
	"io"
	"sync"

	"github.com/netrack/openflow/internal/encoding"
)

type MultipartType uint16

const (
	// Description of this OpenFlow switch.
	// The request body is empty.
	// The reply body is struct Desc.
	MultipartTypeDescription MultipartType = iota

	// Individual flow statistics.
	// The request body is struct FlowStatsRequest.
	// The reply body is an array of struct FlowStats.
	MultipartTypeFlow

	// Aggregate flow statistics.
	// The request body is struct AggregateStatsRequest.
	// The reply body is struct AggregateStatsReply.
	MultipartTypeAggregate

	// Flow table statistics. The request body is empty.
	// The reply body is an array of struct TableStats.
	MultipartTypeTable

	// Port statistics. The request body is struct PortStatsRequest.
	// The reply body is an array of struct PortStats.
	MultipartTypePortStats

	// Queue statistics for a port
	// The request body is struct QueueStatsRequest.
	// The reply body is an array of struct QueueStats
	MultipartTypeQueue

	// Group counter statistics. The request body is struct GroupStatsRequest.
	// The reply is an array of struct GroupStats.
	MultipartTypeGroup

	// Group description. The request body is empty.
	// The reply body is an array of struct GroupDescStats.
	MultipartTypeGroupDescription

	// Group features. The request body is empty.
	// The reply body is struct GroupFeatures.
	MultipartTypeGroupFeatures

	// Meter statistics.
	// The request body is struct MeterMultipartRequests.
	// The reply body is an array of struct MeterStats.
	MultipartTypeMeter

	// Meter configuration.
	// The request body is struct MeterMultipartRequests.
	// The reply body is an array of struct MeterConfig.
	MultipartTypeMeterConfig

	// Meter features. The request body is empty.
	// The reply body is struct MeterFeatures.
	MultipartTypeMeterFeatures

	// Table features.
	// The request body is either empty or contains an array of
	// struct ofp_table_features containing the controller’s
	// desired view of the switch. If the switch is unable to
	// set the specified view an error is returned.
	// The reply body is an array of struct TableFeatures.
	MultipartTypeTableFeatures

	// Port description. The request body is empty.
	// The reply body is an array of struct Port.
	MultipartTypePortDescription

	// Experimenter extension. The request and reply bodies begin with
	// struct ExperimenterMultipartHeader.
	// The request and reply bodies are otherwise experimenter-defined.
	MultipartTypeExperimenter MultipartType = 0xffff
)

type MultipartRequestFlag uint16

const (
	MultipartRequestMode MultipartRequestFlag = 1 << iota
)

type MultipartReplyFlag uint16

const (
	MultipartReplyMode MultipartReplyFlag = 1 << iota
)

// reader is a wrapper around the io.WriterTo type.
type reader struct {
	io.WriterTo

	buf  *bytes.Buffer
	bufE error

	once sync.Once
}

func (r *reader) read(b []byte) (int, error) {
	// Initialize a buffer only once, to dump
	// the content of the writer to the buffer.
	r.once.Do(func() {
		r.buf = new(bytes.Buffer)
		_, r.bufE = r.WriterTo.WriteTo(r.buf)
	})

	if r.bufE != nil {
		return 0, r.bufE
	}

	// And then read the data from the buffer back.
	return r.buf.Read(b)
}

// Read implements io.Reader interface.
func (r *reader) Read(b []byte) (int, error) {
	switch r := r.WriterTo.(type) {
	case io.Reader:
		return r.Read(b)
	}

	return r.read(b)
}

// While the system is running, the controller may request
// state from the datapath using the multipart request message.
type MultipartRequest struct {
	Type  MultipartType
	Flags MultipartRequestFlag

	// Body is the request's body.
	//
	// For receiver requests a nil body means the request has no
	// body.
	Body io.Reader
}

// NewMultipartRequest creates a new multipart request.
func NewMultipartRequest(t MultipartType, body io.WriterTo) *MultipartRequest {
	var rd io.Reader

	// When the body is not defined then bypass wrapping
	// it with a reader type.
	if body != nil {
		rd = &reader{WriterTo: body}
	}

	return &MultipartRequest{t, MultipartRequestMode, rd}
}

func (m *MultipartRequest) WriteTo(w io.Writer) (int64, error) {
	// By default, body will be embedded into the reader type.
	return encoding.WriteTo(w, m.Type, m.Flags, pad4{}, m.Body)
}

func (m *MultipartRequest) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, &m.Type, &m.Flags, &defaultPad4)
	if err != nil {
		return n, err
	}

	// Just copy the rest of the body to the requet.
	buf := new(bytes.Buffer)
	m.Body = buf

	nn, err := io.Copy(buf, r)
	return n + nn, err
}

// The switch responds on T_MULTIPART_REQUEST with one
// or more T_MULTIPART_REPLY messages
type MultipartReply struct {
	Type  MultipartType
	Flags MultipartReplyFlag
}

func (m *MultipartReply) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *m, pad4{})
}

func (m *MultipartReply) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Type, &m.Flags, &defaultPad4)
}

type ExperimenterMultipartHeader struct {
	Experimenter uint32
	ExpType      uint32
}

func (m *ExperimenterMultipartHeader) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *m)
}

func (m *ExperimenterMultipartHeader) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Experimenter, &m.ExpType)
}
