package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// ActionTypeOutput outputs the packet to the switch port.
	ActionTypeOutput ActionType = iota

	// ActionTypeCopyTTLOut copies the TTL from the next-to-outermost header
	// to outermost header with TTL.
	ActionTypeCopyTTLOut ActionType = 10 + iota

	// ActionTypeCopyTTLIn copies the TTL from the outermost header to the
	// next-to-outermost header with TTL.
	ActionTypeCopyTTLIn ActionType = 10 + iota

	// ActionTypeSetMPLSTTL replaces the existing MTPL TTL. This applies
	// only to the packets with existing MPLS shim header.
	ActionTypeSetMPLSTTL ActionType = 12 + iota

	// ActionTypeDecMPLSTTL decrements the MTPS TTL. This applies only to
	// the packets with existing MPLS shim header.
	ActionTypeDecMPLSTTL ActionType = 12 + iota

	// ActionTypePushVLAN pushes a new VLAN header onto the packet.
	ActionTypePushVLAN ActionType = 12 + iota

	// ActionTypePopVLAN pops the outer-most VLAN header from the packet.
	ActionTypePopVLAN ActionType = 12 + iota

	// ActionTypePushMPLS pushes a new MPLS shim header onto the packet.
	ActionTypePushMPLS ActionType = 12 + iota

	// ActionTypePopMPLS pops the outer-most MPLS tag or shim header from
	// the packet.
	ActionTypePopMPLS ActionType = 12 + iota

	// ActionTypeSetQueue specifies on which queue attached to the port
	// should be used to queue and forward packet.
	ActionTypeSetQueue ActionType = 12 + iota

	// ActionTypeGroup specifies that action should be set to the group
	// action, when a packet has to be processed by the group table.
	ActionTypeGroup ActionType = 12 + iota

	// ActionTypeSetNwTTL replaces the existing IPv4 TTL or IPv6 Hop Limit
	// and updates the IP checksum.
	ActionTypeSetNwTTL ActionType = 12 + iota

	// ActionTypeDecNwTTL decrements the IPv4 TTL or IPv6 Hop Limit and
	// and updates the IP checksum.
	ActionTypeDecNwTTL ActionType = 12 + iota

	// ActionTypeSetField sets the value to the packet field.
	ActionTypeSetField ActionType = 12 + iota

	// ActionTypePushPBB pushes a new PBB service instance header (I-TAG)
	// onto the packet.
	ActionTypePushPBB ActionType = 12 + iota

	// ActionTypePopPBB pops the outer-most PBB service instance header
	// (I-TAG) from the packet.
	ActionTypePopPBB ActionType = 12 + iota

	// ActionTypeExperimenter applies the experimental action.
	ActionTypeExperimenter ActionType = 0xffff
)

// ActionType specifies the action type.
type ActionType uint16

const (
	// ContentLenMax defines the maximum length of the bytes, that should
	// be submitted to the controller on output action type.
	ContentLenMax uint16 = 0xffe5

	// ContentLenNoBuffer indicates that no buffering should be applied and
	// the whole packet is to be sent to the controller on output action type.
	ContentLenNoBuffer uint16 = 0xffff
)

type actionhdr struct {
	Type ActionType

	// Length of action, including this header.
	Len uint16
}

type action interface {
	io.WriterTo
}

// Actions groups the set of actions.
type Actions []action

// WriteTo writes the list of action to the given writer instance.
func (a Actions) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, action := range a {
		_, err = action.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	return encoding.WriteTo(w, buf.Bytes())
}

// Action is header that is common to all actions. The length includes
// the header and any padding used to make the action 64-bit aligned.
type Action struct {
	// Type specified type of the action.
	Type ActionType
}

// WriteTo implements io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *Action) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type, 8}, pad4{})
}

// ActionOutput is an action used to output the packets to the switch port.
//
// When the port is the PortController, MaxLen indicates the max number of
// bytes to send. A MaxLen of zero means no bytes of the packet should be
// sent.
//
// A MaxLen of ContentLenNoBuffer means that the packet is not buffered and
// the complete packet is to be sent to the controller.
type ActionOutput struct {
	// Output port.
	Port PortNo

	// Max length to send to controller.
	MaxLen uint16
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionOutput) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeOutput, 16}, *a, pad6{})
}

// ActionGroup is an action that specifis the group used to process
// the packet.
type ActionGroup struct {
	// The GroupID indicates the group used to process this packet.
	// The set of buckets to apply depends on the group type.
	GroupID Group
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionGroup) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeGroup, 8}, *a)
}

// ActionSetQueue sets the queue ID that will be used to map a flow entry
// to an already-configured queue on a port, regardless of the ToS and VLAN
// PCP bits.
//
// The packet should not change as a result of a Set-Queue action. If the
// switch needs to set the ToS/PCP bits for internal handling, the original
// values should be restored before sending the packet out.
type ActionSetQueue struct {
	// The QueueID indicates the queue used to forward the packet.
	QueueID Queue
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetQueue) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeSetQueue, 8}, *a)
}

// ActionSetMPLSTTL is an action used to replace the MPLS TTL
// value of the processing packet.
type ActionSetMPLSTTL struct {
	// The TTL field is the MPLS TTL to set
	TTL uint8
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeSetMPLSTTL, 8}, *a, pad3{})
}

// ActionSetNetworkTTL is an action used to replace the network
// TTL of the processing packet.
type ActionSetNetworkTTL struct {
	// The TTL field is the TTL address to set in the IP header.
	TTL uint8
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeSetNwTTL, 8}, *a, pad3{})
}

// ActionPush is an action used to push the VLAN, MPLS or PBB
// tag onto the processing packet.
type ActionPush struct {
	// Type is the Push-Action type. It should be one of
	// VLAN, MPLS or PBB.
	Type ActionType

	// The EtherType indicates the Ethertype of the new tag.
	//
	// It is used when pushing a new VLAN tag, new MPLS header
	// or PBB service header.
	EtherType uint16
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionPush) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type, 8}, a.EtherType, pad2{})
}

// ActionPopMPLS is an action used to extract the outer-most MPLS tag
// or shim header from the processing packet.
type ActionPopMPLS struct {
	// The EtherType indicates the Ethertype of the payload.
	EtherType uint16
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionPopMPLS) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypePopMPLS, 8}, *a, pad2{})
}

// ActionSetField is an action used to set the value of the packet field.
type ActionSetField struct {
	// Field contains a header field described using
	// a single OXM TLV structure.
	Field XM
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetField) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := a.Field.WriteTo(&buf)
	if err != nil {
		return 0, err
	}

	// Length is padded to 64 bits
	length := buf.Len() + 4

	if length%8 != 0 {
		_, err := buf.Write(make(pad, 8-length%8))
		if err != nil {
			return 0, err
		}
	}

	header := actionhdr{ActionTypeSetField, uint16(buf.Len() + 4)}
	return encoding.WriteTo(w, header, buf.Bytes())
}

// ActionExperimenter is an experimenter action.
type ActionExperimenter struct {
	// The Experimenter identifies the experimental feature.
	Experimenter uint32
}

// WriteTo implements the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionExperimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{ActionTypeExperimenter, 8}, *a)
}
