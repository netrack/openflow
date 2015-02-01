package ofp13

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

type MultipartRequest struct {
	Header Header
	Type   MultipartType
	Flags  MultipartRequestFlags
	Body   []byte
}

type MultipartReply struct {
	Header Header
	Type   MultipartType
	Flags  MultipartReplyFlags
	Body   []byte
}

type ExperimenterMultipartHeader struct {
	Experimenter uint32
	ExpType      uint32
}
