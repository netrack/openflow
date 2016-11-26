package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	// NoBuffer used, when there is no buffered packet
	// associated with the message.
	NoBuffer uint32 = 0xffffffff
)

const (
	// PacketInReasonNoMatch used when there is no matching flow
	// (table-miss flow entry).
	PacketInReasonNoMatch PacketInReason = iota

	// PacketInreasonAction used, when flow entry explicitly
	// outputs to controller.
	PacketInReasonAction

	// PacketInReasonInvalidTTL used when packet has invalid TTL
	PacketInReasonInvalidTTL
)

// PacketInReason represents the reason why this packet have
// been sent to the controller.
type PacketInReason uint8

// PacketIn used by the datapath to sent the processing packet to
// controller.
type PacketIn struct {
	// BufferID is an identifier of the buffer, assigned by
	// datapath, that holds the processing packet.
	BufferID uint32

	// Length is the total length of the frame.
	Length uint16

	// Reason is the reason why packet is being sent.
	Reason PacketInReason

	// TableID is an identifier of the table that was looked up.
	TableID Table

	// Cookie of the flow entry that was looked up.
	Cookie uint64

	// Match is used to match the packet.
	Match Match
}

// Cookies returns the cookie assigned to the rule, that triggered
// the packet-in message to controller.
func (p *PacketIn) Cookies() uint64 {
	return p.Cookie
}

// SetCookies sets the cookie to the packet-in message.
func (p *PacketIn) SetCookies(cookies uint64) {
	p.Cookie = cookies
}

func (p *PacketIn) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.BufferID, p.Length,
		p.Reason, p.TableID, p.Cookie, &p.Match, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It de-serializes the
// packet-in message from binary format.
func (p *PacketIn) ReadFrom(r io.Reader) (int64, error) {
	// Read the packet-in header, then the list of match
	// rules, that used to match the processing packet.
	return encoding.ReadFrom(r, &p.BufferID, &p.Length,
		&p.Reason, &p.TableID, &p.Cookie, &p.Match, &defaultPad2)
}

// PacketOut used by the controller to send a packet out through the
// datapath.
type PacketOut struct {
	// BufferID is an identifier assigned by datapath (NoBuffer if none).
	// The BufferID is the same given in the PacketIn message.
	BufferID uint32

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

// Bytes returns the message in the binary format.
func (p *PacketOut) Bytes() (b []byte) { return Bytes(p) }

// WriteTo implements io.WriterTo interface. It serializes the message
// in the binary format.
func (p *PacketOut) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	// Simply serialize the list of actions into the buffer,
	// and then write the message with necessary paddings.
	_, err = p.Actions.WriteTo(&buf)
	if err != nil {
		return
	}

	return encoding.WriteTo(w, p.BufferID, p.InPort,
		uint16(buf.Len()), pad6{}, buf.Bytes())
}

func (p *PacketOut) ReadFrom(r io.Reader) (int64, error) {
	var len uint16

	n, err := encoding.ReadFrom(r, &p.BufferID, &p.InPort,
		&len, &defaultPad6)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(len))
	p.Actions = nil

	nn, err := p.Actions.ReadFrom(limrd)
	return n + nn, err
}
