package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// ActionType specifies the action type.
type ActionType uint16

// String returns a string representation of the action type.
func (a ActionType) String() string {
	text, ok := actionText[a]
	// If action is now known just say it.
	if !ok {
		return fmt.Sprintf("Action(%d)", a)
	}

	return text
}

const (
	// ActionTypeOutput outputs the packet to the switch port.
	ActionTypeOutput ActionType = iota

	// ActionTypeCopyTTLOut copies the TTL from the next-to-outermost
	// header to outermost header with TTL.
	ActionTypeCopyTTLOut ActionType = 10 + iota

	// ActionTypeCopyTTLIn copies the TTL from the outermost header to
	// the next-to-outermost header with TTL.
	ActionTypeCopyTTLIn ActionType = 10 + iota

	// ActionTypeSetMPLSTTL replaces the existing MTPL TTL. This applies
	// only to the packets with existing MPLS shim header.
	ActionTypeSetMPLSTTL ActionType = 12 + iota

	// ActionTypeDecMPLSTTL decrements the MTPS TTL. This applies only
	// to the packets with existing MPLS shim header.
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

var actionText = map[ActionType]string{
	ActionTypeOutput:       "ActionOutput",
	ActionTypeCopyTTLOut:   "ActionCopyTTLOut",
	ActionTypeCopyTTLIn:    "ActionCopyTTLIn",
	ActionTypeSetMPLSTTL:   "ActionSetMPLSTTL",
	ActionTypeDecMPLSTTL:   "ActionDecMPLSTTL",
	ActionTypePushVLAN:     "ActionPushVLAN",
	ActionTypePopVLAN:      "ActionPopVLAN",
	ActionTypePushMPLS:     "ActionPushMPLS",
	ActionTypePopMPLS:      "ActionPopMPLS",
	ActionTypeSetQueue:     "ActionSetQueue",
	ActionTypeGroup:        "ActionGroup",
	ActionTypeSetNwTTL:     "ActionSetNwTTL",
	ActionTypeDecNwTTL:     "ActionDecNwTTL",
	ActionTypeSetField:     "ActionSetField",
	ActionTypePushPBB:      "ActionPushPBB",
	ActionTypePopPBB:       "ActionPopPBB",
	ActionTypeExperimenter: "ActionExperimenter",
}

var actionMap = map[ActionType]encoding.ReaderMaker{
	ActionTypeOutput:       encoding.ReaderMakerOf(ActionOutput{}),
	ActionTypeCopyTTLOut:   encoding.ReaderMakerOf(ActionCopyTTLOut{}),
	ActionTypeCopyTTLIn:    encoding.ReaderMakerOf(ActionCopyTTLIn{}),
	ActionTypeSetMPLSTTL:   encoding.ReaderMakerOf(ActionSetMPLSTTL{}),
	ActionTypeDecMPLSTTL:   encoding.ReaderMakerOf(ActionDecMPLSTTL{}),
	ActionTypePushVLAN:     encoding.ReaderMakerOf(ActionPushVLAN{}),
	ActionTypePopVLAN:      encoding.ReaderMakerOf(ActionPopVLAN{}),
	ActionTypePushMPLS:     encoding.ReaderMakerOf(ActionPushMPLS{}),
	ActionTypePopMPLS:      encoding.ReaderMakerOf(ActionPopMPLS{}),
	ActionTypeSetQueue:     encoding.ReaderMakerOf(ActionSetQueue{}),
	ActionTypeGroup:        encoding.ReaderMakerOf(ActionGroup{}),
	ActionTypeSetNwTTL:     encoding.ReaderMakerOf(ActionSetNetworkTTL{}),
	ActionTypeDecNwTTL:     encoding.ReaderMakerOf(ActionDecNetworkTTL{}),
	ActionTypeSetField:     encoding.ReaderMakerOf(ActionSetField{}),
	ActionTypePushPBB:      encoding.ReaderMakerOf(ActionPushPBB{}),
	ActionTypePopPBB:       encoding.ReaderMakerOf(ActionPopPBB{}),
	ActionTypeExperimenter: encoding.ReaderMakerOf(ActionExperimenter{}),
}

const (
	// ContentLenMax defines the maximum length of the bytes, that should
	// be submitted to the controller on output action type.
	ContentLenMax uint16 = 0xffe5

	// ContentLenNoBuffer indicates that no buffering should be applied and
	// the whole packet is to be sent to the controller on output action type.
	ContentLenNoBuffer uint16 = 0xffff
)

// action defines a header of each action, it will be used for
// marshalling and unmarshalling actions.
type action struct {
	Type ActionType

	// Length of action, including this header.
	Len uint16
}

func (a *action) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &a.Type, &a.Len)
}

const (
	// actionLen is a minimum length of the action.
	actionLen uint16 = 8

	// actionHeaderLen is a length of the action header.
	actionHeaderLen uint16 = 4
)

// Action is an interface representing an OpenFlow action.
type Action interface {
	encoding.ReadWriter

	// Type returns the type of the single action.
	Type() ActionType
}

// Actions group the set of actions.
type Actions []Action

func (a *Actions) bytes() ([]byte, error) {
	var buf bytes.Buffer

	for _, action := range *a {
		_, err := action.WriteTo(&buf)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// WriteTo writes the list of action to the given writer instance.
func (a *Actions) WriteTo(w io.Writer) (int64, error) {
	buf, err := a.bytes()
	if err != nil {
		return int64(len(buf)), err
	}

	return encoding.WriteTo(w, buf)
}

// ReadFrom decodes the list of actions from the wire format into
// the list of types that implement Action interface.
func (a *Actions) ReadFrom(r io.Reader) (int64, error) {
	var actionType ActionType
	*a = nil

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := actionMap[actionType]; ok {
			rd, err := rm.MakeReader()
			*a = append(*a, rd.(Action))
			return rd, err
		}

		format := "ofp: unknown action type: '%x'"
		return nil, fmt.Errorf(format, actionType)
	}

	return encoding.ScanFrom(r, &actionType,
		encoding.ReaderMakerFunc(rm))
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

// Type retuns the type of the action.
func (a *ActionOutput) Type() ActionType {
	return ActionTypeOutput
}

// WriteTo implements the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionOutput) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), 16}, *a, pad6{})
}

// ReadFrom implements the io.ReaderFrom interface. It deserialized
// the output action from a wire format.
func (a *ActionOutput) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &action{}, &a.Port, &a.MaxLen, &defaultPad6)
}

// ActionCopyTTLOut is an action used to copy TTL from next-to-outermost
// to outermost header.
type ActionCopyTTLOut struct{}

// Type returns type of the action.
func (a *ActionCopyTTLOut) Type() ActionType {
	return ActionTypeCopyTTLOut
}

// WriteTo implements io.WriterTo interface. It serializes
// the "copy TTL out" action with a necessary padding.
func (a *ActionCopyTTLOut) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "copy TTL out" action from a wire format.
func (a *ActionCopyTTLOut) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionCopyTTLIn is an action used to copy TTL from outermost to
// next-to-outermost header.
type ActionCopyTTLIn struct{}

// Type returns type of the action.
func (a *ActionCopyTTLIn) Type() ActionType {
	return ActionTypeCopyTTLIn
}

// WriteTo implements io.WriterTo interface. It serializes
// the "copy TTL in" action with a necessary padding.
func (a *ActionCopyTTLIn) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "copy TTL in" action from a wire format.
func (a *ActionCopyTTLIn) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionSetMPLSTTL is an action used to replace the MPLS TTL
// value of the processing packet.
type ActionSetMPLSTTL struct {
	// The TTL field is the MPLS time-to-live value to set.
	TTL uint8
}

// Type returns type of the action.
func (a *ActionSetMPLSTTL) Type() ActionType {
	return ActionTypeSetMPLSTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "set MPLS TTL" action with a necessary padding.
func (a *ActionSetMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, a.TTL, pad3{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "set MPLS TTL" action from a wire format.
func (a *ActionSetMPLSTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &action{}, &a.TTL, &defaultPad3)
}

// ActionDecMPLSTTL is an actions used to decrement time to live value
// of the MPLS header of the processing packet.
type ActionDecMPLSTTL struct{}

// Type returns type of the action.
func (a *ActionDecMPLSTTL) Type() ActionType {
	return ActionTypeDecMPLSTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "decrement MPLS TTL" action with a necessary padding.
func (a *ActionDecMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "decrement MPLS TTL" action from a wire format.
func (a *ActionDecMPLSTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionPushVLAN is an action used to push the VLAN tag onto the
// processing packet.
type ActionPushVLAN struct {
	// The EtherType indicates the Ethertype of the new tag.
	//
	// It is used when pushing a new VLAN tag, new MPLS header
	// or PBB service header.
	EtherType uint16
}

// Type returns type of the action.
func (a *ActionPushVLAN) Type() ActionType {
	return ActionTypePushVLAN
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "push VLAN" action with a necessary padding.
func (a *ActionPushVLAN) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen},
		a.EtherType, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "push VLAN" action from a wire format.
func (a *ActionPushVLAN) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

// ActionPopVLAN is an action used to pop the VLAN tag from the
// processing packet.
type ActionPopVLAN struct{}

// Type returns type of the action.
func (a *ActionPopVLAN) Type() ActionType {
	return ActionTypePopVLAN
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "pop VLAN" action with a necessary padding.
func (a *ActionPopVLAN) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "pop VLAN" action from a wire format.
func (a *ActionPopVLAN) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionPushMPLS is an action used to push the MPLS tag onto the
// processing packet.
type ActionPushMPLS struct {
	EtherType uint16
}

// Type returns type of the action.
func (a *ActionPushMPLS) Type() ActionType {
	return ActionTypePushMPLS
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "push MPLS" action with a necessary padding.
func (a *ActionPushMPLS) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen},
		a.EtherType, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "push VLAN" action from a wire format.
func (a *ActionPushMPLS) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

// ActionPopMPLS is an action used to extract the outer-most MPLS tag
// or shim header from the processing packet.
type ActionPopMPLS struct {
	// The EtherType indicates the Ethertype of the payload.
	EtherType uint16
}

// Type returns type of the action.
func (a *ActionPopMPLS) Type() ActionType {
	return ActionTypePopMPLS
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "pop MPLS" action with a necessary padding.
func (a *ActionPopMPLS) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen},
		a.EtherType, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "pop MPLS" action from a wire format.
func (a *ActionPopMPLS) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
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

// Type results type of the action.
func (a *ActionSetQueue) Type() ActionType {
	return ActionTypeSetQueue
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "set queue" action with a necessary padding.
func (a *ActionSetQueue) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, a.QueueID)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "set queue" action from a wire format.
func (a *ActionSetQueue) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.QueueID)
}

// ActionGroup is an action that specifies the group used to process
// the packet.
type ActionGroup struct {
	// The Group indicates the group used to process this packet.
	// The set of buckets to apply depends on the group type.
	Group Group
}

// Type returns type of the action.
func (a *ActionGroup) Type() ActionType {
	return ActionTypeGroup
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "group" action with a necessary padding.
func (a *ActionGroup) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, a.Group)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "group" action from a wire format.
func (a *ActionGroup) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.Group)
}

// ActionSetNetworkTTL is an action used to replace the network
// TTL of the processing packet.
type ActionSetNetworkTTL struct {
	// The TTL field is the TTL address to set in the IP header.
	TTL uint8
}

// Type returns type of the action.
func (a *ActionSetNetworkTTL) Type() ActionType {
	return ActionTypeSetNwTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "set network TTL" action with a necessary padding.
func (a *ActionSetNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, a.TTL, pad3{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "set network TTL" action from a wire format.
func (a *ActionSetNetworkTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.TTL, &defaultPad3)
}

// ActionDecNetworkTTL is an actions used to decrement time to live value
// of the network-layer header of the processing packet.
type ActionDecNetworkTTL struct{}

// Type returns type of the action.
func (a *ActionDecNetworkTTL) Type() ActionType {
	return ActionTypeDecNwTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "decrement network TTL" action with a necessary padding.
func (a *ActionDecNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "decrement network TTL" action from a wire format.
func (a *ActionDecNetworkTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionSetField is an action used to set the value of the packet field.
type ActionSetField struct {
	// Field contains a header field described using
	// a single OXM TLV structure.
	Field XM
}

// Type returns type of the action.
func (a *ActionSetField) Type() ActionType {
	return ActionTypeSetField
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "set field" action with a necessary padding.
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

	header := action{a.Type(), uint16(buf.Len() + 4)}
	return encoding.WriteTo(w, header, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "set field" action from a wire format.
func (a *ActionSetField) ReadFrom(r io.Reader) (int64, error) {
	var header action
	var num int64

	n, err := encoding.ReadFrom(r, &header)
	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-4))
	num, err = a.Field.ReadFrom(limrd)
	n += num

	if err != nil {
		return n, err
	}

	b, err := ioutil.ReadAll(limrd)
	n += int64(len(b))
	return n, err
}

// ActionPushPBB is an action used to push a new PBB service tag
// (I-TAG) onto the processing packet.
type ActionPushPBB struct {
	EtherType uint16
}

// Type returns type of the action.
func (a *ActionPushPBB) Type() ActionType {
	return ActionTypePushPBB
}

// WriteTo implement the io.WriterTo interface. It serializes
// the "push PBB" action with a necessary padding.
func (a *ActionPushPBB) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen},
		a.EtherType, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "push PBB" action from a wire format.
func (a *ActionPushPBB) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

// ActionPopPBB is an action used to pop the outer PBB service tag
// (I-TAG) from the processing packet.
type ActionPopPBB struct{}

// ActionExperimenter is an experimenter action.
type ActionExperimenter struct {
	// The Experimenter identifies the experimental feature.
	Experimenter uint32
}

// Type returns type of the action.
func (a *ActionExperimenter) Type() ActionType {
	return ActionTypeExperimenter
}

// WriteTo implements the io.WriterTo interface. It serializes
// the "experimenter" action with a necessary padding.
func (a *ActionExperimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, action{a.Type(), actionLen}, a.Experimenter)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the "experimenter" action from a wire format.
func (a *ActionExperimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.Experimenter)
}
