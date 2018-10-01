package openflow

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	// TypeHello is used by either controller or switch during connection
	// setup. It is used for version negotiation. When the connection
	// is established, each side must immediately send a Hello message
	// with the version field set to the highest version supported by
	// the sender. If the version negotiation fails, an Error message
	// is sent with type HelloFailed and code Incompatible.
	TypeHello Type = iota

	// TypeError can be sent by either the switch or the controller and
	// indicates the failure of an operation. The simplest failure pertain
	// to malformed messages or failed version negotiation, while more
	// complex scenarios desbie some failure in state change at the switch.
	TypeError

	// TypeEchoRequest is used to exchange information about latency,
	// bandwidth and liveness. Echo request timeout indicates disconnection.
	TypeEchoRequest

	// TypeEchoReply is used to exchange information about latency,
	// bandwidth and liveness. Echo reply is sent as a response to Echo
	// request.
	TypeEchoReply

	// TypeExperiment is a mechanism for proprietary messages within the
	// protocol.
	TypeExperiment

	// TypeFeaturesRequest is used when a transport channel (TCP, SCTP,
	// TLS) is established between the switch and controller, the first
	// activity is feature determination. The controller will send a
	// feature request to the switch over the transport channel.
	TypeFeaturesRequest

	// TypeFeaturesReply is the switch's reply to the controller
	// enumerating its abilities.
	TypeFeaturesReply

	// TypeGetConfigRequest sequence is used to query and set the
	// fragmentation handling properties of the packet processing pipeline.
	TypeGetConfigRequest

	// TypeGetConfigReply is the switch's reply to the controller
	// that acknowledges configuration requests.
	TypeGetConfigReply

	// TypeSetConfig is used by the controller to alter the switch's
	// configuration. This message is unacknowledged.
	TypeSetConfig

	// TypePacketIn message is a way for the switch to send a captured
	// packet to the controller.
	TypePacketIn

	// TypeFlowRemoved is sent to the controller by the switch when a
	// flow entry in a flow table is removed. It happens when a timeout
	// occurs, either due to inactivity or hard timeout. An idle timeout
	// happens when no packets are matched in a period of time. A hard
	// timeout happens when a certain period of time elapses, regardless
	// of the number of matching packets. Whether the switch sends a
	// TypeFlowRemoved message after a timeout is specified by the
	// TypeFlowMod. Flow entry removal with a TypeFlowMod message from
	// the controller can also lead to a TypeFlowRemoved message.
	TypeFlowRemoved

	// TypePortStatus messages are asynchronous events sent from the
	// switch to the controller indicating a change of status for the
	// indicated port.
	TypePortStatus

	// TypePacketOut is used by the controller to inject packets into the
	// data plane of a particular switch. Such message can either carry a
	// raw packet to inject into the switch, or indicate a local buffer
	// on the switch containing a raw packet to release. Packets injected
	// into the data plane of a switch using this method are not treated
	// the same as packets that arrive on standard ports. The packet jumps
	// to the action set application stage in the standard packet processing
	// pipeline. If the action set indicates table processing is necessary
	// then the input port id is used as the arrival port of the raw packet.
	TypePacketOut

	// TypeFlowMod is one of the main messages, it allows the controller
	// to modify the state of switch.
	TypeFlowMod

	// TypeGroupMod is used by controller to modify group tables.
	TypeGroupMod

	// TypePortMod is used by the controller to modify the state of port.
	TypePortMod

	// TypeTableMod is used to determine a packet's fate when it misses
	// in the table. It can be forwarded to the controller, dropped, or
	// sent to the next table.
	TypeTableMod

	// TypeMultipartRequest is used by the controller to request the
	// state of the datapath.
	TypeMultipartRequest

	// TypeMultipartReply are the replies from the switch to controller
	// on TypeMultipartRequest messages.
	TypeMultipartReply

	// TypeBarrierRequest can be used by the controller to set a
	// synchronization point, ensuring that all previous state messages
	// are completed before the barrier response is sent back to the
	// controller.
	TypeBarrierRequest

	// TypeBarrierReply is a response from the switch to controller
	// on TypeBarrierRequest messages.
	TypeBarrierReply

	// TypeQueueGetConfigRequest can be used by the controller to
	// query the state of queues associated with various ports on switch.
	TypeQueueGetConfigRequest

	// TypeQueueGetConfigReply is a response from the switch to controller
	// on TypeQueueGetConfigReply messages.
	TypeQueueGetConfigReply

	// TypeRoleRequest is the message used by the controller to
	// modify its role among multiple controllers on a switch.
	TypeRoleRequest

	// TypeRoleReply is a response from the switch to controller on
	// TypeRoleRequest.
	TypeRoleReply

	// TypeGetAsyncRequest is used by the controller to request the switch
	// which asynchronous events are enabled on the switch for the
	// communication channel.
	TypeGetAsyncRequest

	// TypeGetAsyncReply is used by the switch as response to controller
	// on TypeAsyncRequest messages.
	TypeGetAsyncReply

	// TypeSetAsync is used by the controller to set which asynchronous
	// messages it should send, as well as to query the switch for which
	// asynchronous messages it will send.
	TypeSetAsync

	// TypeMeterMod used by the controller to modify the meter.
	TypeMeterMod
)

// Type is an OpenFlow message type.
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
	TypeGetAsyncRequest:       "TypeGetAsyncRequest",
	TypeGetAsyncReply:         "TypeGetAsyncReply",
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

// Copy returns a copy of the request header.
func (h *Header) Copy() *Header {
	return &Header{h.Version, h.Type, h.Length, h.Transaction}
}

// Len of the packet payload including header.
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
