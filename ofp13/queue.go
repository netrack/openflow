package ofp13

const (
	QT_MIN_RATE     QueueProperties = iota
	QT_MAX_RATE     QueueProperties = iota
	QT_EXPERIMENTER QueueProperties = 0xffff
)

type QueueProperties uint16

type PacketQueue struct {
	QueueId    uint32
	Port       uint32
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
