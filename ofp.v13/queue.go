package ofp

const (
	QT_MIN_RATE     QueueProperties = iota
	QT_MAX_RATE     QueueProperties = iota
	QT_EXPERIMENTER QueueProperties = 0xffff
)

type QueueProperties uint16

const (
	//TODO: Q_ANY Queue
	QueueAll Queue = 0xffffffff
)

type Queue uint32

type PacketQueue struct {
	QueueID    Queue
	Port       PortNo
	Length     uint16
	_          pad6
	Properties []QueuePropHeader
}

type QueuePropHeader struct {
	Property QueueProperties
	Length   uint16
	_        pad3
}

type QueuePropMinRate struct {
	PropHeader QueuePropHeader
	Rate       uint16
	_          pad6
}

type QueuePropMaxRate struct {
	PropHeader QueuePropHeader
	Rate       uint16
	_          pad6
}

type QueuePropExperimenter struct {
	PropHeader   QueuePropHeader
	Experimenter uint32
	_            pad4
	Data         []byte
}

type QueueStatsRequest struct {
	PornNo  PortNo
	QueueID Queue
}

type QueueStatus struct {
	PortNo       PortNo
	QueueID      Queue
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	DurationSec  uint32
	DurationNSec uint32
}

type QueueGetConfigRequest struct {
	Port PortNo
	_    pad4
}

type QueueGetConfigReply struct {
	Port   PortNo
	_      pad4
	Queues []PacketQueue
}
