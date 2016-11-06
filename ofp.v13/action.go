package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

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

type actionhdr struct {
	Type ActionType

	// Length of action, including this header.
	Len uint16
}

const actionLen uint16 = 8

type Action interface {
	encoding.ReadWriter

	Type() ActionType
}

// Actions groups the set of actions.
type Actions []Action

func (a Actions) bytes() ([]byte, error) {
	var buf bytes.Buffer

	for _, action := range a {
		_, err := action.WriteTo(&buf)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// WriteTo writes the list of action to the given writer instance.
func (a Actions) WriteTo(w io.Writer) (int64, error) {
	buf, err := a.bytes()
	if err != nil {
		return int64(len(buf)), err
	}

	return encoding.WriteTo(w, buf)
}

func (a Actions) ReadFrom(r io.Reader) (int64, error) {
	var actionType ActionType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := actionMap[actionType]; ok {
			rd, err := rm.MakeReader()
			a = append(a, rd.(Action))
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

func (a *ActionOutput) Type() ActionType {
	return ActionTypeOutput
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionOutput) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), 16}, *a, pad6{})
}

func (a *ActionOutput) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &actionhdr{}, &a.Port, &a.MaxLen, &defaultPad6)
}

type ActionCopyTTLOut struct{}

func (a *ActionCopyTTLOut) Type() ActionType {
	return ActionTypeCopyTTLOut
}

func (a *ActionCopyTTLOut) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, pad4{})
}

func (a *ActionCopyTTLOut) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

type ActionCopyTTLIn struct{}

func (a *ActionCopyTTLIn) Type() ActionType {
	return ActionTypeCopyTTLIn
}

func (a *ActionCopyTTLIn) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, pad4{})
}

func (a *ActionCopyTTLIn) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionSetMPLSTTL is an action used to replace the MPLS TTL
// value of the processing packet.
type ActionSetMPLSTTL struct {
	// The TTL field is the MPLS TTL to set
	TTL uint8
}

func (a *ActionSetMPLSTTL) Type() ActionType {
	return ActionTypeSetMPLSTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, a.TTL, pad3{})
}

func (a *ActionSetMPLSTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &actionhdr{}, &a.TTL, &defaultPad3)
}

type ActionDecMPLSTTL struct{}

func (a *ActionDecMPLSTTL) Type() ActionType {
	return ActionTypeDecMPLSTTL
}

func (a *ActionDecMPLSTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, pad4{})
}

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

func (a *ActionPushVLAN) Type() ActionType {
	return ActionTypePushVLAN
}

func (a *ActionPushVLAN) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen},
		a.EtherType, pad2{})
}

func (a *ActionPushVLAN) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

type ActionPopVLAN struct{}

func (a *ActionPopVLAN) Type() ActionType {
	return ActionTypePopVLAN
}

func (a *ActionPopVLAN) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, pad4{})
}

func (a *ActionPopVLAN) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

type ActionPushMPLS struct {
	EtherType uint16
}

func (a *ActionPushMPLS) Type() ActionType {
	return ActionTypePushMPLS
}

func (a *ActionPushMPLS) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen},
		a.EtherType, pad2{})
}

func (a *ActionPushMPLS) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

// ActionPopMPLS is an action used to extract the outer-most MPLS tag
// or shim header from the processing packet.
type ActionPopMPLS struct {
	// The EtherType indicates the Ethertype of the payload.
	EtherType uint16
}

func (a *ActionPopMPLS) Type() ActionType {
	return ActionTypePopMPLS
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionPopMPLS) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen},
		a.EtherType, pad2{})
}

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

func (a *ActionSetQueue) Type() ActionType {
	return ActionTypeSetQueue
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetQueue) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, a.QueueID)
}

func (a *ActionSetQueue) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.QueueID)
}

// ActionGroup is an action that specifis the group used to process
// the packet.
type ActionGroup struct {
	// The GroupID indicates the group used to process this packet.
	// The set of buckets to apply depends on the group type.
	GroupID Group
}

func (a *ActionGroup) Type() ActionType {
	return ActionTypeGroup
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionGroup) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, a.GroupID)
}

func (a *ActionGroup) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.GroupID)
}

// ActionSetNetworkTTL is an action used to replace the network
// TTL of the processing packet.
type ActionSetNetworkTTL struct {
	// The TTL field is the TTL address to set in the IP header.
	TTL uint8
}

func (a *ActionSetNetworkTTL) Type() ActionType {
	return ActionTypeSetNwTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionSetNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, a.TTL, pad3{})
}

func (a *ActionSetNetworkTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.TTL, &defaultPad3)
}

type ActionDecNetworkTTL struct{}

func (a *ActionDecNetworkTTL) Type() ActionType {
	return ActionTypeDecNwTTL
}

// WriteTo implement the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionDecNetworkTTL) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, pad4{})
}

func (a *ActionDecNetworkTTL) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// ActionSetField is an action used to set the value of the packet field.
type ActionSetField struct {
	// Field contains a header field described using
	// a single OXM TLV structure.
	Field XM
}

func (a *ActionSetField) Type() ActionType {
	return ActionTypeSetField
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

	header := actionhdr{a.Type(), uint16(buf.Len() + 4)}
	return encoding.WriteTo(w, header, buf.Bytes())
}

func (a *ActionSetField) ReadFrom(r io.Reader) (int64, error) {
	var header actionhdr
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

type ActionPushPBB struct {
	EtherType uint16
}

func (a *ActionPushPBB) Type() ActionType {
	return ActionTypePushPBB
}

func (a *ActionPushPBB) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen},
		a.EtherType, pad2{})
}

func (a *ActionPushPBB) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.EtherType, &defaultPad2)
}

type ActionPopPBB struct{}

// ActionExperimenter is an experimenter action.
type ActionExperimenter struct {
	// The Experimenter identifies the experimental feature.
	Experimenter uint32
}

func (a *ActionExperimenter) Type() ActionType {
	return ActionTypeExperimenter
}

// WriteTo implements the io.WriterTo interface. It serializes
// the action with a necessary padding.
func (a *ActionExperimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, actionhdr{a.Type(), actionLen}, a.Experimenter)
}

func (a *ActionExperimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad4, &a.Experimenter)
}
