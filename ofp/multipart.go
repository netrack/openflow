package ofp

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/netrack/openflow/internal/encoding"
)

// MultipartType defines the type of multipart request. It specifies the
// kind of information being passed and determined how the body of the
// request is interpreted.
type MultipartType uint16

const (
	// MultipartTypeDescription is used to retrieve description of the
	// OpenFlow switch.
	//
	// The request body is empty. The reply body is struct Description.
	MultipartTypeDescription MultipartType = iota

	// MultipartTypeFlow is used to retrieve individual flow statistics.
	//
	// The request body is struct FlowStatsRequest. The reply body is an
	// array of struct FlowStats.
	MultipartTypeFlow

	// MultipartTypeAggregate is used to aggregate flow statistics.
	//
	// The request body is struct AggregateStatsRequest. The reply body is
	// struct AggregateStatsReply.
	MultipartTypeAggregate

	// MultipartTypeTable is used to retrieve flow table statistics.
	//
	// The request body is empty. The reply body is an array of struct
	// TableStats.
	MultipartTypeTable

	// MultipartTypePortStats is used to retrieve port statistics.
	//
	// The request body is struct PortStatsRequest. The reply body is an
	// array of struct PortStats.
	MultipartTypePortStats

	// MultipartTypeQueue is used to retrieve queue statistics for a port.
	//
	// The request body is struct QueueStatsRequest. The reply body is an
	// array of struct QueueStats.
	MultipartTypeQueue

	// MultipartTypeGroup is used to retrieve group counter statistics.
	//
	// The request body is struct GroupStatsRequest. The reply is an array
	// of struct GroupStats.
	MultipartTypeGroup

	// MultipartTypeGroupDescription is used to retrieve group description.
	//
	// The request body is empty. The reply body is an array of struct
	// GroupDescStats.
	MultipartTypeGroupDescription

	// MultipartTypeGroupFeatures is used to retrieve group features.
	//
	// The request body is empty. The reply body is struct GroupFeatures.
	MultipartTypeGroupFeatures

	// MultipartTypeMeter is used to retrieve meter statistics.
	//
	// The request body is struct MeterMultipartRequests. The reply body is
	// an array of struct MeterStats.
	MultipartTypeMeter

	// MultipartTypeMeterConfig is used to retrieve meter configuration.
	//
	// The request body is struct MeterMultipartRequests. The reply body is
	// an array of struct MeterConfig.
	MultipartTypeMeterConfig

	// MultipartTypeMeterFeatures is used to retrieve meter features.
	//
	// The request body is empty. The reply body is struct MeterFeatures.
	MultipartTypeMeterFeatures

	// MultipartTypeTableFeatures is used to retrieve table features.
	//
	// The request body is either empty or contains an array of struct
	// TableFeatures containing the controller’s desired view of the
	// switch. If the switch is unable to set the specified view an error
	// is returned.
	//
	// The reply body is an array of struct TableFeatures.
	MultipartTypeTableFeatures

	// MultipartTypePortDescription is used to retrieve port description.
	//
	// The request body is empty. The reply body is an array of struct Port.
	MultipartTypePortDescription

	// MultipartTypeExperimenter is an experimenter extension.
	//
	// The request and reply bodies begin with struct
	// ExperimenterMultipartHeader. The request and reply bodies are
	// otherwise experimenter-defined.
	MultipartTypeExperimenter MultipartType = 0xffff
)

func (t MultipartType) String() string {
	text, ok := multipartTypeText[t]
	if !ok {
		return fmt.Sprintf("MultipartType(%d)", t)
	}
	return text
}

var multipartTypeText = map[MultipartType]string{
	MultipartTypeDescription:      "MultipartTypeDescription",
	MultipartTypeFlow:             "MultipartTypeFlow",
	MultipartTypeAggregate:        "MultipartTypeAggregate",
	MultipartTypeTable:            "MultipartTypeTable",
	MultipartTypePortStats:        "MultipartTypePortStats",
	MultipartTypeQueue:            "MultipartTypeQueue",
	MultipartTypeGroup:            "MultipartTypeGroup",
	MultipartTypeGroupDescription: "MultipartTypeGroupDescription",
	MultipartTypeGroupFeatures:    "MultipartTypeGroupFeatures",
	MultipartTypeMeter:            "MultipartTypeMeter",
	MultipartTypeMeterConfig:      "MultipartTypeMeterConfig",
	MultipartTypeMeterFeatures:    "MultipartTypeMeterFeatures",
	MultipartTypeTableFeatures:    "MultipartTypeTableFeatures",
	MultipartTypePortDescription:  "MultipartTypePortDescription",
	MultipartTypeExperimenter:     "MultipartTypeExperimenter",
}

// MultipartRequestFlag defines the multipart request flags.
type MultipartRequestFlag uint16

const (
	// MultipartRequestMode is set when more requests to follow.
	MultipartRequestMode MultipartRequestFlag = 1 << iota
)

// MultipartReplyFlag defines the multipart reply flags.
type MultipartReplyFlag uint16

const (
	// MultipartReplyMode is set when more responses replies to follow.
	MultipartReplyMode MultipartReplyFlag = 1 << iota
)

// reader is a wrapper around the io.WriterTo type.
type reader struct {
	io.WriterTo

	// Buffer use to store the content of the WriterTo instance.
	//
	// If the write to the buffer will fail, an error will be saved
	// into the bufE error. It will be returned on each attempt to
	// read from this buffer.
	buf  *bytes.Buffer
	bufE error

	// Instance of the once type used to dump content of WriterTo
	// to the buffer only once.
	once sync.Once
}

// read is used to enclose the content of the WriterTo instance into
// the underlying buffer instance.
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

// Read implements io.Reader interface. When the underlying instance of
// the io.WriterTo interface implements io.Reader, it will be used to
// read the given bytes. Otherwise, a buffer of bytes will be used to
// store the bytes.
func (r *reader) Read(b []byte) (int, error) {
	switch r := r.WriterTo.(type) {
	case io.Reader:
		return r.Read(b)
	}

	return r.read(b)
}

// MultipartRequest is a message that controller may send in order to
// request the state from the datapath while the system is running.
//
// For example, to retrieve the flow table statistics, the following
// request can be sent:
//
//	body := ofp.NewMultipartRequest(ofp.MultipartTypeTable, nil)
//
//	req := of.NewRequest(of.TypeMultipartRequest, body)
//	// ...
type MultipartRequest struct {
	Type  MultipartType
	Flags MultipartRequestFlag

	// Body is the request's body.
	//
	// For receiver requests a nil body means the request has no
	// body.
	Body io.Reader
}

// NewMultipartRequest creates a new multipart request of the given
// type. If the body is not equal to nil, it will be sent as part of
// the multipart request message.
func NewMultipartRequest(t MultipartType, body io.WriterTo) *MultipartRequest {
	var rd io.Reader

	// When the body is not defined then bypass wrapping
	// it with a reader type.
	if body != nil {
		rd = &reader{WriterTo: body}
	}

	return &MultipartRequest{t, 0, rd}
}

// WriteTo implements io.WriterTo interface. It serializes the multipart
// request into the wire format.
func (m *MultipartRequest) WriteTo(w io.Writer) (int64, error) {
	// By default, body will be embedded into the reader type.
	return encoding.WriteTo(w, m.Type, m.Flags, pad4{}, m.Body)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// multipart request from the wire format.
func (m *MultipartRequest) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, &m.Type, &m.Flags, &defaultPad4)
	if err != nil {
		return n, err
	}

	// Just copy the rest of the body to the request.
	buf := new(bytes.Buffer)
	m.Body = buf

	nn, err := io.Copy(buf, r)
	return n + nn, err
}

// MultipartReply is a message used by the datapath to reply with the
// data requested by controller while system is running.
type MultipartReply struct {
	// Type of the multipart reply.
	Type MultipartType

	// Flags of the multipart reply message.
	Flags MultipartReplyFlag
}

// WriteTo implements io.WriterTo interface. It serializes the multipart
// reply into the wire format.
func (m *MultipartReply) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *m, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// multipart request from the wire format.
func (m *MultipartReply) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Type, &m.Flags, &defaultPad4)
}

// ExperimenterMultipartHeader is a header of the experimenter multipart
// messages (both, requests and replies).
type ExperimenterMultipartHeader struct {
	// Experimenter identifier.
	Experimenter uint32

	// Experimenter type.
	ExpType uint32
}

// WriteTo implements io.WriterTo interface. It serializes the multipart
// experimenter header into the wire format.
func (m *ExperimenterMultipartHeader) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *m)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the multipart
// experimenter header from the wire format.
func (m *ExperimenterMultipartHeader) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Experimenter, &m.ExpType)
}
