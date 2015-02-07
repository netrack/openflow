package ofp13

type AggregateStatsRequest struct {
	TableId    Table
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
