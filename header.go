package openflow

import (
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// Immutable messages
	T_HELLO Type = iota
	T_ERROR
	T_ECHO_REQUEST
	T_ECHO_REPLY
	T_EXPERIMENTER

	// Switch configuration messages
	T_FEATURES_REQUEST
	T_FEATURES_REPLY
	T_GET_CONFIG_REQUEST
	T_GET_CONFIG_REPLY
	T_SET_CONFIG

	// Asynchronous messages
	T_PACKET_IN
	T_FLOW_REMOVED
	T_PORT_STATUS

	// Controller command messages
	T_PACKET_OUT
	T_FLOW_MOD
	T_GROUP_MOD
	T_PORT_MOD
	T_TABLE_MOD

	// Multipart messages
	T_MULTIPART_REQUEST
	T_MULTIPART_REPLY

	// Queue configuration messages
	T_QUEUE_GET_CONFIG_REQUEST
	T_QUEUE_GET_CONFIG_REPLY

	// Controller role change request messages
	T_ROLE_REQUEST
	T_ROLE_REPLY

	// Asynchronous message configuration
	T_ASYNC_REQUEST
	T_ASYNC_REPLY
	T_SET_ASYNC

	// Meters and rate limiters configuration messages
	T_METER_MOD
)

type Type uint8

// Each OpenFlow message begins with the OpenFlow header
type Header struct {
	Version uint8
	// One of the Type constants
	Type Type
	// Length including this header
	Length uint16
	// Transaction id associated with this packet.
	// Replies use the same id as was in the request
	// to facilitate pairing
	Xid uint32
}

func (h *Header) Len() int {
	return int(h.Length)
}

func (h *Header) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h)
}

func (h *Header) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, h)
}
