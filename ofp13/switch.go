package ofp13

const (
	C_FLOW_STATS   Capabilities = 1 << iota
	C_TABLE_STATS  Capabilities = 1 << iota
	C_PORT_STATS   Capabilities = 1 << iota
	C_GROUP_STATS  Capabilities = 1 << iota
	C_IP_REASM     Capabilities = 1 << iota
	C_QUEUE_STATS  Capabilities = 1 << iota
	C_PORT_BLOCKED Capabilities = 1 << iota
)

type Capabilities uint32

const (
	C_FRAG_NORMAL ConfigFlags = iota
	C_FRAG_DROP   ConfigFlags = 1 << 0
	C_FRAG_REASM  ConfigFlags = 1 << 1
	C_FRAG_MASK   ConfigFlags = iota
)

type ConfigFlags uint16

type SwitchFeatures struct {
	Header       Header
	DatapathId   uint64
	NumBuffers   uint32
	NumTables    uint8
	AuxiliaryId  uint8
	Capabilities Capabilities
	Reserved     uint32
}

type SwitchConfig struct {
	Header         Header
	Flags          ConfigFlags
	MissSendLength uint16
}
