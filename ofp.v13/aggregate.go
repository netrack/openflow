package ofp

// Aggregate information about multiple flow entries is requested.
type AggregateStatsRequest struct {
	TableID    Table
	_          pad3
	OutPort    PortNo
	OutGroup   Group
	_          pad4
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

type AggreagateStatsReply struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
	_           pad4
}
