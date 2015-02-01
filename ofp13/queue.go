package ofp13

const (
	QT_MIN_RATE     QueueProperties = iota
	QT_MAX_RATE     QueueProperties = iota
	QT_EXPERIMENTER QueueProperties = 0xffff
)

type QueueProperties uint16

const (
	//TODO: Q_ANY Queue
	Q_ALL Queue = 0xffffffff
)

type Queue uint32

type PacketQueue struct {
	QueueId    Queue
	Port       PortNo
	Length     uint16
	Properties []QueuePropHeader
}

type QueuePropHeader struct {
	Property QueueProperties
	Length   uint16
}

type QueuePropMinRate struct {
	PropHeader QueuePropHeader
	Rate       uint16
}

type QueuePropMaxRate struct {
	PropHeader QueuePropHeader
	Rate       uint16
}

type QueuePropExperimenter struct {
	PropHeader   QueuePropHeader
	Experimenter uint32
	date         []byte
}

type QueueStatsRequest struct {
	PornNo  PortNo
	QueueId Queue
}

type QueueStatus struct {
	PortNo       PortNo
	QueueId      Queue
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	DurationSec  uint32
	DurationNSec uint32
}

type QueueGetConfigRequest struct {
	Header Header
	Port   PortNo
}

type QueueGetConfigReply struct {
	Header Header
	Port   PortNo
	Queues []PacketQueue
}
