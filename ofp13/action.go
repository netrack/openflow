package ofp13

import (
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	AT_OUTPUT       ActionType = iota
	AT_COPY_TTL_OUT ActionType = 10 + iota
	AT_COPY_TTL_IN  ActionType = 10 + iota
	AT_SET_MPLS_TTL ActionType = 12 + iota
	AT_DEC_MPLS_TTL ActionType = 12 + iota
	AT_PUSH_VLAN    ActionType = 12 + iota
	AT_POP_VLAN     ActionType = 12 + iota
	AT_PUSH_MPLS    ActionType = 12 + iota
	AT_POP_MPLS     ActionType = 12 + iota
	AT_SET_QUEUE    ActionType = 12 + iota
	AT_GROUP        ActionType = 12 + iota
	AT_SET_NW_TTL   ActionType = 12 + iota
	AT_DEC_NW_TTL   ActionType = 12 + iota
	AT_SET_FIELD    ActionType = 12 + iota
	AT_PUSH_PBB     ActionType = 12 + iota
	AT_POP_PBB      ActionType = 12 + iota
	AT_EXPERIMENTER ActionType = 0xffff
)

type ActionType uint16

const (
	// Maximum MaxLength value which can be used to request a specific byte Length.
	CML_MAX uint16 = 0xffe5
	// Indicates that no buffering should be applied and
	// the whole packet is to be sent to the controller.
	CML_NO_BUFFER uint16 = 0xffff
)

type actionhdr struct {
	Type ActionType
	// Length of action, including this header. This is the length of action,
	// including any padding to make it 64-bit aligned
	Len uint16
}

type action interface {
	io.WriterTo
}

type Actions []action

// Action header that is common to all actions. The length includes the
// header and any padding used to make the action 64-bit aligned.
// NB: The length of an action *must* always be a multiple of eight.
type Action struct {
	Type ActionType
}

func (a *Action) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{a.Type, 8}, pad4{},
	})
}

// Action structure for PAT_OUTPUT, which sends packets out port.
// When the port is the PP_CONTROLLER, MaxLen indicates the max
// number of bytes to send. A MaxLen of zero means no bytes of the
// packet should be sent. A MaxLen of CML_NO_BUFFER means that
// the packet is not buffered and the complete packet is to be sent to
// the controller.
type ActionOutput struct {
	// Output port
	Port PortNo
	// Max length to send to controller
	MaxLen uint16
}

func (a ActionOutput) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_OUTPUT, 16}, a, pad6{},
	})
}

// Action structure for PAT_GROUP
type ActionGroup struct {
	// The group_id indicates the group used to process this packet.
	// The set of buckets to apply depends on the group type.
	GroupID uint32
}

func (a ActionGroup) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_GROUP, 8}, a,
	})
}

// ActionSetQueue sets the queue id that will be used to map a
// flow entry to an already-configured queue on a port,
// regardless of the ToS and VLAN PCP bits. The packet should not change as a result
// of a Set-Queue action. If the switch needs to set the
// ToS/PCP bits for internal handling, the original values should
// be restored before sending the packet out.
type ActionSetQueue struct {
	QueueID uint32
}

func (a ActionSetQueue) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_SET_QUEUE, 8}, a,
	})
}

// Action structure for PAT_SET_MPLS_TTL
type ActionMPLSTTL struct {
	// The TTL field is the MPLS TTL to set
	TTL uint8
}

func (a ActionMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_SET_MPLS_TTL, 8}, a, pad3{},
	})
}

// Action structure for PAT_SET_NW_TTL
type ActionSetNetworkTTL struct {
	// The TTL field is the TTL address to set in the IP header.
	TTL uint8
}

func (a ActionSetNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_SET_NW_TTL, 8}, a, pad3{},
	})
}

// Action structure for PAT_PUSH_VLAN/MPLS/PBB
type ActionPush struct {
	// One of PAT_PUSH_VLAN/MPLS/PBB
	Type ActionType
	// The EtherType indicates the Ethertype of the new tag.
	// It is used when pushing a new VLAN tag, new
	// MPLS header or PBB service header.
	EtherType uint16
}

func (a ActionPush) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{a.Type, 8}, a.EtherType, pad2{},
	})
}

// Action structure for OFPAT_POP_MPLS
type ActionPopMPLS struct {
	// The EtherType indicates the Ethertype of the payload.
	EtherType uint16
}

func (a ActionPopMPLS) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_POP_MPLS, 8}, a, pad2{},
	})
}

type ActionSetField struct {
	Fields []OXM
}

// Action header for PAT_EXPERIMENTER.
// The rest of the body is experimenter-defined
type ActionExperimenter struct {
	Experimenter uint32
}

func (a *ActionExperimenter) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		actionhdr{AT_EXPERIMENTER, 8}, a,
	})
}
