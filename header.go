package openflow

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/netrack/openflow/internal/encoding"
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

	// Barrier messages
	TypeBarrierRequest
	TypeBarrierReply

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

func (t Type) String() string {
	text, ok := typeText[t]
	if !ok {
		return fmt.Sprintf("Type(%d)", t)
	}

	return text
}

var typeText = map[Type]string{
	TypeHello:                 "TypeHello",
	TypeError:                 "TypeError",
	TypeEchoRequest:           "TypeEchoRequest",
	TypeEchoReply:             "TypeEchoReply",
	TypeExperiment:            "TypeExperiment",
	TypeFeaturesRequest:       "TypeFeaturesRequest",
	TypeFeaturesReply:         "TypeFeaturesReply",
	TypeGetConfigRequest:      "TypeGetConfigRequest",
	TypeGetConfigReply:        "TypeGetConfigReply",
	TypeSetConfig:             "TypeSetConfig",
	TypePacketIn:              "TypePacketIn",
	TypeFlowRemoved:           "TypeFlowRemoved",
	TypePortStatus:            "TypePortStatus",
	TypePacketOut:             "TypePacketOut",
	TypeFlowMod:               "TypeFlowMod",
	TypeGroupMod:              "TypeGroupMod",
	TypePortMod:               "TypePortMod",
	TypeTableMod:              "TypeTableMod",
	TypeMultipartRequest:      "TypeMultipartRequest",
	TypeMultipartReply:        "TypeMultipartReply",
	TypeBarrierRequest:        "TypeBarrierRequest",
	TypeBarrierReply:          "TypeBarrierReply",
	TypeQueueGetConfigRequest: "TypeQueueGetConfigRequest",
	TypeQueueGetConfigReply:   "TypeQueueGetConfigReply",
	TypeRoleRequest:           "TypeRoleRequest",
	TypeRoleReply:             "TypeRoleReply",
	TypeAsynchRequest:         "TypeAsynchRequest",
	TypeAsyncReply:            "TypeAsyncReply",
	TypeSetAsync:              "TypeSetAsync",
	TypeMeterMod:              "TypeMeterMod",
}

// The Header is a response header. It contains the negotiated
// version of the OpenFlow, a type and length of the message.
type Header struct {
	// Version specifies the version of the protocol.
	Version uint8

	// Type defines a type of the message.
	Type Type

	// Length including this Header.
	Length uint16

	// Transaction is an transaction ID associated with this packet.
	//
	// Replies use the same id as was in the request to facilitate pairing.
	Transaction uint32
}

func (h *Header) Copy() *Header {
	return &Header{h.Version, h.Type, h.Length, h.Transaction}
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
	return encoding.ReadFrom(r, &h.Version, &h.Type, &h.Length, &h.Transaction)
}

// TransactionMatcher creates a new matcher that matches the request
// by the transaction identifier.
//
// If the header has non-zero transaction identifier, it will be used
// to create a new matcher, otherwise a random number will be generated.
func TransactionMatcher(h *Header) Matcher {
	// When the trasnaction is not defined, we will generate a
	// random number, that will be used to match the response.
	if h.Transaction == 0 {
		h.Transaction = rand.Uint32()
	}

	transaction := h.Transaction
	matcher := func(r *Request) bool {
		return r.Header.Transaction == transaction
	}

	// Return a function wrapped into the function adapter.
	return &MatcherFunc{matcher}
}
