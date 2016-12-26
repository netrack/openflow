package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// QueuePropType defines the type of the queue property.
type QueuePropType uint16

const (
	// QueuePropTypeMinRate indicates that queue guarantees minimum
	// datarate.
	QueuePropTypeMinRate QueuePropType = 1

	// QueuePropTypeMaxRate indicates that queue guarantees maximum
	// datarate.
	QueuePropTypeMaxRate QueuePropType = 2

	// QueuePropTypeExperimenter indicates an experimental queue
	// property.
	QueuePropTypeExperimenter QueuePropType = 0xffff
)

// queuePropTypeMap is a mapping used to decode the set of queue properties.
var queuePropTypeMap = map[QueuePropType]encoding.ReaderMaker{
	QueuePropTypeMinRate:      encoding.ReaderMakerOf(QueuePropMinRate{}),
	QueuePropTypeMaxRate:      encoding.ReaderMakerOf(QueuePropMaxRate{}),
	QueuePropTypeExperimenter: encoding.ReaderMakerOf(QueuePropExperimenter{}),
}

// Queue defines a queue number configured at the specific port.
type Queue uint32

const (
	// QueueAll refers to all queues configured at the specified port.
	QueueAll Queue = 0xffffffff
)

// QueueProp in an interface representing an OpenFlow queue property.
type QueueProp interface {
	encoding.ReadWriter

	// Type returns the queue property type.
	Type() QueuePropType
}

// QueueProps used to consolidate a set of queue properties.
type QueueProps []QueueProp

// WriteTo implements io.WriterTo interface. It serializes the set
// of queue properties into the wire format.
func (q QueueProps) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	for _, prop := range q {
		_, err := prop.WriteTo(&buf)
		if err != nil {
			return 0, err
		}
	}

	return encoding.WriteTo(w, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// set of queue properties from the wire format.
func (q QueueProps) ReadFrom(r io.Reader) (int64, error) {
	var queueType QueuePropType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := queuePropTypeMap[queueType]; ok {
			rd, err := rm.MakeReader()
			q = append(q, rd.(QueueProp))
			return rd, err
		}

		format := "ofp: unknown queue property type: '%x'"
		return nil, fmt.Errorf(format, queueType)
	}

	return encoding.ScanFrom(r, &queueType,
		encoding.ReaderMakerFunc(rm))
}

// packetQueueLen defines the length of the packet queue header.
const packetQueueLen = 16

// PacketQueue decribes the packet processing queue.
//
// An OpenFlow switch provides limited Quality-of-Service support (QoS)
// through a simple queuing mechanism. One (or more) queues can attach
// to a port and be used to map flow entries on it. Flow entries mapped
// to a specific queue will be treated according to that queue's
// configuration (e.g. min rate).
type PacketQueue struct {
	// Queue identifies the specified queue.
	Queue Queue

	// Port this queue attached to.
	Port PortNo

	// Properties is a list of queue properties.
	Properties QueueProps
}

// WriteTo implements io.WriterTo interface. It serializes the packet
// queue into the wire format.
func (q *PacketQueue) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	_, err := q.Properties.WriteTo(&buf)
	if err != nil {
		return 0, err
	}

	return encoding.WriteTo(w, q.Queue, q.Port,
		uint16(buf.Len()+packetQueueLen), pad6{}, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// packet queue from the wire format.
func (q *PacketQueue) ReadFrom(r io.Reader) (int64, error) {
	var length uint16
	n, err := encoding.ReadFrom(r, &q.Queue, &q.Port,
		&length, &defaultPad6)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(length-packetQueueLen))
	nn, err := q.Properties.ReadFrom(limrd)
	return n + nn, err
}

// queuePropLen defines the length of the queue property header length.
const queuePropLen = 16

// queueProp is a common header of the queue properties.
type queueProp struct {
	Type QueuePropType
	Len  uint16
}

const (
	// QueueMinRateUncfg indicates that minimum-rate queue property is not
	// configured.
	QueueMinRateUncfg uint16 = 0xffff

	// QueueMaxRateUncfg indicates that maximum-rate queue property is not
	// configured.
	QueueMaxRateUncfg uint16 = 0xffff
)

// QueuePropMinRate defines the minimum-rate queue property.
type QueuePropMinRate struct {
	// Rate in 1/10 of a percent. If value is more than 1000,
	// the rate is disabled.
	Rate uint16
}

// Type implements QueueProp interface. It returns the type of minimum-rate
// queue property.
func (q *QueuePropMinRate) Type() QueuePropType {
	return QueuePropTypeMinRate
}

// WriteTo implements io.WriterTo interface. It serializes the queue property
// into the wire format.
func (q *QueuePropMinRate) WriteTo(w io.Writer) (int64, error) {
	header := queueProp{q.Type(), queuePropLen}
	return encoding.WriteTo(w, header, pad4{}, q.Rate, pad6{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the queue
// property from the wire format.
func (q *QueuePropMinRate) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8, &q.Rate, &defaultPad6)
}

// QueuePropMaxRate defines the maximum-rate queue property.
type QueuePropMaxRate struct {
	// Rate in 1/10 of a percent. If value is more than 1000,
	// the rate is disabled.
	Rate uint16
}

// Type implements QueueProp interface. It returns the type of
// maximum-rate queue property.
func (q *QueuePropMaxRate) Type() QueuePropType {
	return QueuePropTypeMaxRate
}

// WriteTo implements io.WriterTo interface. It serializes the
// maximum-rate property into the wire format.
func (q *QueuePropMaxRate) WriteTo(w io.Writer) (int64, error) {
	header := queueProp{q.Type(), queuePropLen}
	return encoding.WriteTo(w, header, pad4{}, q.Rate, pad6{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// maximum-rate queue property from the wire format.
func (q *QueuePropMaxRate) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8, &q.Rate, &defaultPad6)
}

// QueuePropExperimenter defines an experimental queue property.
type QueuePropExperimenter struct {
	// Experimenter identifier.
	Experimenter uint32

	// Experimenter-defined data.
	Data []byte
}

// Type implements QueueProp interface. It returns the type of the
// experimental queue property.
func (q *QueuePropExperimenter) Type() QueuePropType {
	return QueuePropTypeExperimenter
}

// WriteTo implements io.WriterTo interface. It serializes the
// experimental queue property into the wire format.
func (q *QueuePropExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := queueProp{q.Type(), queuePropLen + uint16(len(q.Data))}
	return encoding.WriteTo(w, header, pad4{},
		q.Experimenter, &pad4{}, q.Data)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// queue property from the wire format.
func (q *QueuePropExperimenter) ReadFrom(r io.Reader) (int64, error) {
	var header queueProp
	n, err := encoding.ReadFrom(r, &header, &defaultPad4,
		&q.Experimenter, &defaultPad4)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-queuePropLen))
	q.Data, err = ioutil.ReadAll(limrd)
	return n + int64(len(q.Data)), err
}

// QueueStatsRequest is a multipart request used to retrieve queue
// statistics for one or more ports and one or more queues.
//
// For example, to retrieve statistics of all queues configured on the
// second port, the following request can be created:
//
//	body := ofp.QueueStatsRequest{2, ofp.QueueAll}
//	req := of.NewRequest(of.TypeMultipartRequest,
//		ofp.NewMultipartRequest(ofp.MultipartTypeQueue, body))
type QueueStatsRequest struct {
	// Port identifier or PortAny for all ports.
	Port PortNo

	// Queue identifier or QueueAll for all queues.
	Queue Queue
}

// WriteTo implements io.WriterTo interface. It serializes the queue
// statistics request into the wire format.
func (q *QueueStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, q.Port, q.Queue)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// queue statistics request from the wire format.
func (q *QueueStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &q.Queue)
}

// QueueStats defines the statistics of the single queue. The list of
// QueueStats structures will be returned as a multipart response on
// queue statistics request.
type QueueStats struct {
	// Port uniquely identifies a port within a switch.
	Port PortNo

	// Queue identifier.
	Queue Queue

	// TxBytes is a number of transmitted bytes.
	TxBytes uint64

	// TxPackets is a number of transmitted bytes.
	TxPackets uint64

	// TxErrors is a number of transmitted errors.
	TxErrors uint64

	// DurationSec is a time queue has been alive in seconds.
	DurationSec uint32

	// DurationNSec is a time queue has been alive in nanoseconds
	// beyond DurationSec.
	DurationNSec uint32
}

// WriteTo implements io.WriterTo interface. It serializes the queue
// statistics into the wire format.
func (q *QueueStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *q)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// queue statistics from the wire format.
func (q *QueueStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &q.Queue, &q.TxBytes,
		&q.TxPackets, &q.TxErrors, &q.DurationSec, &q.DurationNSec)
}

// QueueGetConfigRequest is a message used to query the switch for
// configured queues on a port.
type QueueGetConfigRequest struct {
	// Port is a to be queried. Should refer to a valid physical
	// port (i.e. < PortMax), or PortAny to request all configured
	// queues.
	Port PortNo
}

// WriteTo implements io.WriterTo interface. It serializes the messages
// used to query the switch for configured queues into the wire format.
func (q *QueueGetConfigRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, q.Port, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// message used to query the switch for configured queues from the wire
// format.
func (q *QueueGetConfigRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &defaultPad4)
}

// QueueGetConfigReply is message used by switch to reply on query for
// configured queues. It contains a list of configured queues.
type QueueGetConfigReply struct {
	// Port identifies a port within a switch.
	Port PortNo

	// Queues is a list of packet queues.
	Queues []PacketQueue
}

// WriteTo implements io.WriterTo interface. It serializes the queue
// configuration reply into the wire format.
func (q *QueueGetConfigReply) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	for _, queue := range q.Queues {
		_, err := queue.WriteTo(&buf)
		if err != nil {
			return 0, err
		}
	}

	return encoding.WriteTo(w, q.Port, pad4{}, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// queue configuration reply from the wire format.
func (q *QueueGetConfigReply) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, &q.Port, &defaultPad4)
	if err != nil {
		return n, err
	}

	// The passed reader should be limited to the whole
	// OpenFlow message, thus return EOF error when the
	// messages ends. Otherwise this implementation will
	// read the packet queues indefinitely.
	queueMaker := encoding.ReaderMakerOf(PacketQueue{})
	nn, err := encoding.ReadSliceFrom(r, queueMaker, q.Queues)
	return n + nn, err
}
