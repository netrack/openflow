package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

const (
	QueuePropTypeMinRate      QueuePropType = 1 + iota
	QueuePropTypeMaxRate      QueuePropType = 1 + iota
	QueuePropTypeExperimenter QueuePropType = 0xffff
)

type QueuePropType uint16

var queuePropTypeMap = map[QueuePropType]encoding.ReaderMaker{
	QueuePropTypeMinRate:      encoding.ReaderMakerOf(QueuePropMinRate{}),
	QueuePropTypeMaxRate:      encoding.ReaderMakerOf(QueuePropMaxRate{}),
	QueuePropTypeExperimenter: encoding.ReaderMakerOf(QueuePropExperimenter{}),
}

const (
	//TODO: Q_ANY Queue
	QueueAll Queue = 0xffffffff
)

type Queue uint32

type QueueProp interface {
	encoding.ReadWriter

	Type() QueuePropType
}

type QueueProps []QueueProp

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

const packetQueueLen = 16

type PacketQueue struct {
	Queue      Queue
	Port       PortNo
	Properties QueueProps
}

func (q *PacketQueue) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	_, err := q.Properties.WriteTo(&buf)
	if err != nil {
		return 0, err
	}

	return encoding.WriteTo(w, q.Queue, q.Port,
		uint16(buf.Len()+packetQueueLen), pad6{}, buf.Bytes())
}

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

const queuePropLen = 16

type queuePropHdr struct {
	Type QueuePropType
	Len  uint16
}

type QueuePropMinRate struct {
	Rate uint16
}

func (q *QueuePropMinRate) Type() QueuePropType {
	return QueuePropTypeMinRate
}

func (q *QueuePropMinRate) WriteTo(w io.Writer) (int64, error) {
	header := queuePropHdr{q.Type(), queuePropLen}
	return encoding.WriteTo(w, header, pad4{}, q.Rate, pad6{})
}

func (q *QueuePropMinRate) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8, &q.Rate, &defaultPad6)
}

type QueuePropMaxRate struct {
	Rate uint16
}

func (q *QueuePropMaxRate) Type() QueuePropType {
	return QueuePropTypeMaxRate
}

func (q *QueuePropMaxRate) WriteTo(w io.Writer) (int64, error) {
	header := queuePropHdr{q.Type(), queuePropLen}
	return encoding.WriteTo(w, header, pad4{}, q.Rate, pad6{})
}

func (q *QueuePropMaxRate) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8, &q.Rate, &defaultPad6)
}

type QueuePropExperimenter struct {
	Experimenter uint32
	Data         []byte
}

func (q *QueuePropExperimenter) Type() QueuePropType {
	return QueuePropTypeExperimenter
}

func (q *QueuePropExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := queuePropHdr{q.Type(), queuePropLen + uint16(len(q.Data))}
	return encoding.WriteTo(w, header, pad4{},
		q.Experimenter, &pad4{}, q.Data)
}

func (q *QueuePropExperimenter) ReadFrom(r io.Reader) (int64, error) {
	var header queuePropHdr
	n, err := encoding.ReadFrom(r, &header, &defaultPad4,
		&q.Experimenter, &defaultPad4)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-queuePropLen))
	q.Data, err = ioutil.ReadAll(limrd)
	return n + int64(len(q.Data)), err
}

type QueueStatsRequest struct {
	Port  PortNo
	Queue Queue
}

func (q *QueueStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, q.Port, q.Queue)
}

func (q *QueueStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &q.Queue)
}

type QueueStats struct {
	Port         PortNo
	Queue        Queue
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	DurationSec  uint32
	DurationNSec uint32
}

func (q *QueueStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *q)
}

func (q *QueueStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &q.Queue, &q.TxBytes,
		&q.TxPackets, &q.TxErrors, &q.DurationSec, &q.DurationNSec)
}

type QueueGetConfigRequest struct {
	Port PortNo
}

func (q *QueueGetConfigRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, q.Port, pad4{})
}

func (q *QueueGetConfigRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &q.Port, &defaultPad4)
}

type QueueGetConfigReply struct {
	Port   PortNo
	Queues []PacketQueue
}

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

func (q *QueueGetConfigReply) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, &q.Port, &defaultPad4)
	if err != nil {
		return n, err
	}

	// The passed reader should be limited to the whole
	// OpenFlow message, thus return EOF error when the
	// messages ends. Otherwise this implementation will
	// read the packet queues indefinitely.
	var queues []PacketQueue
	for {
		var queue PacketQueue
		nn, err := queue.ReadFrom(r)
		n += nn

		if err != nil {
			return n, encoding.SkipEOF(err)
		}

		queues = append(queues, queue)
	}

	// Assign the list of queues back to the reply to
	// make the unmarshaling reproducable.
	q.Queues = queues
	return n, nil
}
