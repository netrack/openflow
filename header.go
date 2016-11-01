package of

import (
	"io"

	"github.com/netrack/openflow/encoding"
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

// The Header is a response header. It contains the negotiated
// version of the OpenFlow, a type and length of the message.
type Header struct {
	// Version specifies the version of the protocol.
	Version uint8

	// Type defines a type of the message.
	Type Type

	// Length including this Header.
	Length uint16

	// XId is an transaction id associated with this packet.
	//
	// Replies use the same id as was in the request to facilitate pairing.
	XID uint32
}

func (h *Header) Copy() *Header {
	return &Header{h.Version, h.Type, h.Length, h.XID}
}

// Length of the packet payload including header.
func (h *Header) Len() int {
	return int(h.Length)
}

// WriteTo writes the header in the write format to the given writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *h)
}

// ReadFrom reads the header from the given reader in the wire format.
func (h *Header) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &h.Version, &h.Type, &h.Length, &h.XID)
}
