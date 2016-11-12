package ofp

import (
	"io"

	"github.com/netrack/openflow/encoding"
)

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

type MultipartType uint16

const (
	MultipartRequestMode MultipartRequestFlag = 1 << iota
)

type MultipartRequestFlag uint16

const (
	MultipartReplyMode MultipartReplyFlag = 1 << iota
)

type MultipartReplyFlag uint16

// While the system is running, the controller may request
// state from the datapath using the T_MULTIPART_REQUEST message
type MultipartRequest struct {
	Type  MultipartType
	Flags MultipartRequestFlag
}

func (m *MultipartRequest) Bytes() []byte {
	return Bytes(m)
}

func (m *MultipartRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Type, m.Flags, pad4{})
}

func (m *MultipartRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Type, &m.Flags, &defaultPad4)
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
