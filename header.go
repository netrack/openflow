package of

import (
	"errors"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// Immutable messages.
	TypeHello Type = iota
	TypeError
	TypeEchoRequest
	TypeEchoReply
	TypeExperiment

	// Switch configuration messages.
	TypeFeaturesRequest
	TypeFeaturesReply
	TypeGetConfigRequest
	TypeGetConfigReply
	TypeSetConfig

	// Asynchronous messages.
	TypePacketIn
	TypeFlowRemoved
	TypePortStatus

	// Controller command messages.
	TypePacketOut
	TypeFlowMod
	TypeGroupMod
	TypePortMod
	TypeTableMod

	// Multipart messages
	TypeMultipartRequest
	TypeMultipartReply

	// Queue configuration messages.
	TypeQueueGetConfigRequest
	TypeQueueGetConfigReply

	// Controller role change request messages.
	TypeRoleRequest
	TypeRoleReply

	// Asynchronous message configuration.
	TypeAsynchRequest
	TypeAsyncReply
	TypeSetAsync

	// Meters and rate limiters configuration messages.
	TypeMeterMod
)

type Type uint8

const (
	VersionHeaderKey HeaderKey = iota
	TypeHeaderKey
	XIDHeaderKey
)

type HeaderKey int

type Header interface {
	io.WriterTo
	io.ReaderFrom

	Set(k HeaderKey, v interface{}) error
	Get(k HeaderKey) interface{}
	Len() int
}

// Each OpenFlow message begins with the OpenFlow header
type header struct {
	Version uint8
	// One of the Type constants
	Type Type
	// Length including this header
	Length uint16
	// Transaction id associated with this packet.
	// Replies use the same id as was in the request
	// to facilitate pairing
	XID uint32
}

func (h *header) Set(k HeaderKey, v interface{}) error {
	switch k {
	case VersionHeaderKey:
		return h.setVersion(v)
	case TypeHeaderKey:
		return h.setType(v)
	case XIDHeaderKey:
		return h.setXID(v)
	default:
		return errors.New("header: unsettable field")
	}

	return nil
}

func (h *header) setVersion(v interface{}) error {
	version, ok := v.(uint8)
	if !ok {
		return errors.New("header: Version must be uint8")
	}

	h.Version = version
	return nil
}

func (h *header) setType(v interface{}) error {
	typ, ok := v.(Type)
	if !ok {
		return errors.New("header: Type must be uint8")
	}

	h.Type = Type(typ)
	return nil
}

func (h *header) setXID(v interface{}) error {
	xid, ok := v.(uint32)
	if !ok {
		return errors.New("header: XID must be uint32")
	}

	h.XID = xid
	return nil
}

func (h *header) Get(k HeaderKey) (v interface{}) {
	switch k {
	case VersionHeaderKey:
		v = h.Version
	case TypeHeaderKey:
		v = h.Type
	case XIDHeaderKey:
		v = h.XID
	}

	return
}

func (h *header) Len() int {
	return int(h.Length)
}

func (h *header) WriteTo(w io.Writer) (int64, error) {
	return binary.Write(w, binary.BigEndian, h)
}

func (h *header) ReadFrom(r io.Reader) (int64, error) {
	return binary.Read(r, binary.BigEndian, h)
}
