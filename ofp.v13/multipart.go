package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// Description of this OpenFlow switch.
	// The request body is empty.
	// The reply body is struct desc.
	MP_DESC MultipartType = iota

	// Individual flow statistics.
	// The request body is struct FlowStatsRequest.
	// The reply body is an array of struct FlowStats.
	MP_FLOW

	// Aggregate flow statistics.
	// The request body is struct AggregateStatsRequest.
	// The reply body is struct AggregateStatsReply.
	MP_AGGREGATE

	// Flow table statistics. The request body is empty.
	// The reply body is an array of struct TableStats.
	MP_TABLE

	// Port statistics. The request body is struct PortStatsRequest.
	// The reply body is an array of struct PortStats.
	MP_PORT_STATS

	// Queue statistics for a port
	// The request body is struct QueueStatsRequest.
	// The reply body is an array of struct QueueStats
	MP_QUEUE

	// Group counter statistics. The request body is struct GroupStatsRequest.
	// The reply is an array of struct GroupStats.
	MP_GROUP

	// Group description. The request body is empty.
	// The reply body is an array of struct GroupDescStats.
	MP_GROUP_DESC

	// Group features. The request body is empty.
	// The reply body is struct GroupFeatures.
	MP_GROUP_FEATURES

	// Meter statistics.
	// The request body is struct MeterMultipartRequests.
	// The reply body is an array of struct MeterStats.
	MP_METER

	// Meter configuration.
	// The request body is struct MeterMultipartRequests.
	// The reply body is an array of struct MeterConfig.
	MP_METER_CONFIG

	// Meter features. The request body is empty.
	// The reply body is struct MeterFeatures.
	MP_METER_FEATURES

	// Table features.
	// The request body is either empty or contains an array of
	// struct ofp_table_features containing the controller’s
	// desired view of the switch. If the switch is unable to
	// set the specified view an error is returned.
	// The reply body is an array of struct TableFeatures.
	MP_TABLE_FEATURES

	// Port description. The request body is empty.
	// The reply body is an array of struct Port.
	MP_PORT_DESC

	// Experimenter extension. The request and reply bodies begin with
	// struct ExperimenterMultipartHeader.
	// The request and reply bodies are otherwise experimenter-defined.
	MP_EXPERIMENTER MultipartType = 0xffff
)

type MultipartType uint16

const (
	MPF_REQ_MODE MultipartRequestFlags = 1 << iota
)

type MultipartRequestFlags uint16

const (
	MPF_REPLY_MODE MultipartReplyFlags = 1 << iota
)

type MultipartReplyFlags uint16

// While the system is running, the controller may request
// state from the datapath using the T_MULTIPART_REQUEST message
type MultipartRequest struct {
	Type  MultipartType
	Flags MultipartRequestFlags
	Body  io.WriterTo
}

func (m *MultipartRequest) Bytes() []byte {
	return Bytes(m)
}

func (m *MultipartRequest) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	if m.Body != nil {
		_, err = m.Body.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		m.Type, m.Flags, pad4{}, buf.Bytes(),
	})
}

// The switch responds on T_MULTIPART_REQUEST with one
// or more T_MULTIPART_REPLY messages
type MultipartReply struct {
	Type  MultipartType
	Flags MultipartReplyFlags
}

func (m *MultipartReply) ReadFrom(r io.Reader) (int64, error) {
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&m.Type, &m.Flags, &pad4{},
	})
}

type ExperimenterMultipartHeader struct {
	Experimenter uint32
	ExpType      uint32
}
