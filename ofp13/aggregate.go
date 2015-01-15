package ofp13

type AggregateStatsRequest struct {
	TableId Table
	OutPort PortNo
	//OutGroup GroupNo
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

type AggreagateStatsReply struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
}
