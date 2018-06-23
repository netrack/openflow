package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	// NoBuffer used as buffer identifier when there is no buffered packet
	// associated with the message.
	NoBuffer uint32 = 0xffffffff
)

const (
	// PacketInReasonNoMatch is set when there is no matching flow
	// (table-miss flow entry).
	PacketInReasonNoMatch PacketInReason = iota

	// PacketInReasonAction is set when flow entry explicitly outputs to
	// controller.
	PacketInReasonAction

	// PacketInReasonInvalidTTL is set when packet has invalid TTL.
	PacketInReasonInvalidTTL
)

// PacketInReason represents the reason why this packet have been sent
// to the controller.
type PacketInReason uint8

func (r PacketInReason) String() string {
	text, ok := packetInReasonText[r]
	if !ok {
		return fmt.Sprintf("PacketInReason(%d)", r)
	}
	return text
}

var packetInReasonText = map[PacketInReason]string{
	PacketInReasonNoMatch:    "PacketInReasonNoMatch",
	PacketInReasonAction:     "PacketInReasonAction",
	PacketInReasonInvalidTTL: "PacketInReasonInvalidTTL",
}

// PacketIn used by the datapath to sent the processing packet to
// controller.
//
// For example, to create rule used to flood all unknown packets through
// all ports of the switch, the following request can be sent:
//
//	var packet ofp.PacketIn
//	packet.ReadFrom(r.Body)
//
//	apply := ofp.InstructionApplyActions{
//		ofp.Actions{&of.ActionOutput{ofp.PortFlood, 0}},
//	}
//
//	// For each incoming packet-in request, create a
//	// respective flow modification command.
//	fmod := ofp.NewFlowMod(ofp.FlowAdd, packet)
//	fmod.Instructions = ofp.Instructions{apply}
//
//	req := of.NewRequest(TypeFlowMod, fmod)
type PacketIn struct {
	// Buffer is an identifier of the buffer, assigned by
	// datapath, that holds the processing packet.
	Buffer uint32

	// Length is the total length of the frame.
	Length uint16

	// Reason is the reason why packet is being sent.
	Reason PacketInReason

	// Table is an identifier of the table that was looked up.
	Table Table

	// Cookie of the flow entry that was looked up.
	Cookie uint64

	// Match is used to match the packet.
	Match Match

	// Data represents the original ethernet frame received by the datapath.
	Data []byte
}

// Cookies returns the cookie assigned to the rule, that triggered the
// packet-in message to controller.
func (p *PacketIn) Cookies() uint64 {
	return p.Cookie
}

// SetCookies sets the cookie to the packet-in message.
func (p *PacketIn) SetCookies(cookies uint64) {
	p.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the packet-in
// message into the wire format.
func (p *PacketIn) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.Buffer, p.Length,
		p.Reason, p.Table, p.Cookie, &p.Match, pad2{}, p.Data)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// packet-in message from wire format.
func (p *PacketIn) ReadFrom(r io.Reader) (int64, error) {
	// Read the Packet-In header, list of match rules
	// that used to match the processing packet,
	// then the original frame, if any.
	n, err := encoding.ReadFrom(r, &p.Buffer, &p.Length,
		&p.Reason, &p.Table, &p.Cookie, &p.Match, &defaultPad2)
	if err != nil {
		return n, err
	}
	p.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return n + int64(len(p.Data)), err
	}

	return n + int64(len(p.Data)), nil
}

// PacketOut used by the controller to send a packet out through the
// datapath.
type PacketOut struct {
	// Buffer is an identifier assigned by datapath (NoBuffer if none).
	// The Buffer is the same given in the PacketIn message.
	Buffer uint32

	// InPort is the ingress port that must be associated with the packet
	// for OpenFlow processing. It must be set to either a valid standard
	// switch port or PortController.
	InPort PortNo

	// Action field is an action list defining how the packet should be
	// processed by the switch. It may include packet modification, group
	// processing and an output port.
	Actions Actions

	// Followed by packet data. The length is inferred from the length
	// field in the header.
}

// WriteTo implements io.WriterTo interface. It serializes the message
// into the wire format.
func (p *PacketOut) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	// Simply serialize the list of actions into the buffer,
	// and then write the message with necessary paddings.
	_, err = p.Actions.WriteTo(&buf)
	if err != nil {
		return
	}

	return encoding.WriteTo(w, p.Buffer, p.InPort,
		uint16(buf.Len()), pad6{}, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// packet-out message from the wire format.
func (p *PacketOut) ReadFrom(r io.Reader) (int64, error) {
	var len uint16

	n, err := encoding.ReadFrom(r, &p.Buffer, &p.InPort,
		&len, &defaultPad6)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(len))
	p.Actions = nil

	nn, err := p.Actions.ReadFrom(limrd)
	return n + nn, err
}
