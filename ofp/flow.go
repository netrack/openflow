package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// FlowModCommand represents a type of the flow table modification
// message.
type FlowModCommand uint8

const (
	// FlowAdd is a command used to add a new flow.
	FlowAdd FlowModCommand = iota

	// FlowModify is a command used to modify all matching flows.
	FlowModify

	// FlowModifyStrict is a command used to modify entry strictly
	// matching wildcards and priority.
	FlowModifyStrict

	// FlowDelete is a command used to delete all matching flows.
	FlowDelete

	// FlowDeleteStrict is a command used to delete entry strictly
	// matching wildcards and priority.
	FlowDeleteStrict
)

// FlowModFlag defines flags used in flow modification message.
type FlowModFlag uint16

const (
	// FlowFlagSendFlowRem instructs the switch to send a flow removed
	// message when the flow entry expires or is deleted.
	FlowFlagSendFlowRem FlowModFlag = 1 << iota

	// FlowFlagCheckOverlap instructs the switch to check that there are
	// no conflicting entries with the same priority prior to inserting it
	// in the flow table.
	//
	// If there is one, the flow mod fails and an error message is
	// returned.
	FlowFlagCheckOverlap

	// FlowFlagResetCounts instructs the switch to resets flow packet and
	// byte counts.
	FlowFlagResetCounts

	// FlowFlagNoPktCounts instructs the switch to not keep track of the
	// flow packet count.
	FlowFlagNoPktCounts

	// FlowFlagNoByteCounts instructs the switch to no keep track of the
	// flow byte count.
	FlowFlagNoByteCounts
)

// FlowMod represents a modification message to a flow table from the
// controller.
//
// For example, to create a flow entry to forward all packets arriving
// on the first port to the second port:
//
//	fmod := ofp.NewFlowMod(ofp.FlowAdd, nil)
//
//	// Match all packets that arrive on port number 1.
//	fmod.Match = ofputil.ExtendedMatch(ofputil.MatchInPort(1))
//
//	// Apply the output action, that will forward all
//	// matching packets to the port number 2.
//	fmod.Instructions = ofp.Instructions{
//		&InstructionApplyActions{ofp.Actions{
//			&ofp.ActionOutput{2, 0},
//		}},
//	}
//
//	// Create a request from the assembled message.
//	req := of.NewRequest(of.FlowMod, &fmod)
type FlowMod struct {
	// The Cookie is an opaque data value chosen by the controller.
	//
	// This value appears in flow removed messages and flow statistics,
	// and can also be used to filter flow statistics, flow modification
	// and flow deletion.
	Cookie uint64

	// The CookieMask is used with the cookie field to restrict flow
	// matching while modifying or deleting flow entries.
	//
	// This field is ignored by flow addition messages. A value of 0
	// indicates no restriction.
	CookieMask uint64

	// The Table is an id of the table to put the flow in.
	//
	// For flow deletion commands, TableAll can also be used to delete
	// matching flows from all tables.
	Table Table

	// Command specifies a flow modification command.
	Command FlowModCommand

	// The IdleTimeoute specifies time before discarding a flow entry
	// (in seconds).
	//
	// If the IdleTimeout is set and the HardTimeout is zero, the entry
	// must expire after IdleTimeout seconds with no received traffic.
	//
	// If the IdleTimeout is zero and the HardTimeout is set, the entry
	// must expire in HardTimeout seconds
	// regardless of whether or not packets are hitting the entry.
	IdleTimeout uint16

	// HardTimeout specifis max time before discarding a flow entry (in
	// seconds).
	//
	// If both IdleTimeout and HardTimeout are set, the flow entry will
	// timeout after IdleTimeout seconds with no traffic, or HardTimeout
	// seconds, whichever comes first.
	//
	// If both IdleTimeout and HardTimeout are zero, the entry is
	// considered permanent and will never time out.
	HardTimeout uint16

	// The Priority indicates priority within the specified flow table
	// table.
	//
	// Higher numbers indicate higher priorities. This field is used only
	// for flow addition messages when matching and adding flow entries,
	// and for flow modification and deletion messages when matching flow
	// entries.
	Priority uint16

	// The Buffer refers to a packet buffered at the switch and sent
	// to the controller by a packet-in message.
	//
	// If no buffered packet is associated with the flow mod, it must be
	// set to NoBuffer.
	//
	// A flow mod that includes a valid Buffer is effectively equivalent
	// to sending a two-message sequence of a flow mod and a packet-out to
	// PortTable, with the requirement that the switch must fully process
	// the flow mod before the packet out.
	Buffer uint32

	// For flow deletion commands, require matching entries to include
	// this as an output port. A value of PortAny indicates no restriction.
	OutPort PortNo

	// For flow deletion commands, require matching entries to include this
	// as an output group. A value of GroupAny indicates no restriction.
	OutGroup Group

	// Flags specifies a set of flow modification flags.
	Flags FlowModFlag

	// Match lists fields to match.
	Match Match

	// The Instructions contain the instruction set for the flow entry
	// when adding or modifying entries.
	//
	// If the instruction set is not valid or supported, the switch must
	// generate an error.
	Instructions Instructions
}

// NewFlowMod creates a flow modification message based on the specified
// packet-in message.
//
// It is responsibility of the caller to assign the missing instructions
// and the rest of parameters.
func NewFlowMod(c FlowModCommand, p *PacketIn) *FlowMod {
	var flags FlowModFlag

	switch c {
	case FlowDelete, FlowDeleteStrict:
		// For flow delete requests we cannot set the same flags
		// as for flow insertion and modification.
	default:
		// Use the overlap checking and flow removed notification
		// flags by default for generated message.
		flags = FlowFlagSendFlowRem
	}

	// When the packet-in message was not provided into
	// the constructor, we will use the default values.
	buffer := NoBuffer
	var match Match

	if p != nil {
		buffer, match = p.Buffer, p.Match
	}

	return &FlowMod{
		Command: c,
		Buffer:  buffer,
		Match:   match,
		Flags:   flags,

		// For FlowDelete command, define the output port and
		// group values as "any" to indicate no restrictions.
		OutPort:  PortAny,
		OutGroup: GroupAny,
	}
}

// Cookies implements CookieJar interface. It returns flow mod message
// cookies.
func (f *FlowMod) Cookies() uint64 {
	return f.Cookie
}

// SetCookies implements CookieJar. It sets cookies to flow mod message.
func (f *FlowMod) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the flow
// modification command into the wire format with a necessary padding.
func (f *FlowMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Cookie, f.CookieMask, f.Table,
		f.Command, f.IdleTimeout, f.HardTimeout, f.Priority,
		f.Buffer, f.OutPort, f.OutGroup, f.Flags, pad2{},
		&f.Match, &f.Instructions,
	)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the flow
// modification command from the wire format.
func (f *FlowMod) ReadFrom(r io.Reader) (int64, error) {
	// Set the list of instructions to nil, if someone will decide
	// to reuse the same exemplar for multiple deserializations.
	f.Instructions = nil

	return encoding.ReadFrom(r, &f.Cookie, &f.CookieMask, &f.Table,
		&f.Command, &f.IdleTimeout, &f.HardTimeout, &f.Priority,
		&f.Buffer, &f.OutPort, &f.OutGroup, &f.Flags, &defaultPad2,
		&f.Match, &f.Instructions,
	)
}

// FlowRemovedReason specifies the reason of the flow entry removal.
type FlowRemovedReason uint8

const (
	// FlowReasonIdleTimeout is set when flow was removed because of idle
	// time have been exceeded IdleTimeout.
	FlowReasonIdleTimeout FlowRemovedReason = iota

	// FlowReasonHardTimeout is set when flow was removed because of time
	// have been exceeded HardTimeout.
	FlowReasonHardTimeout

	// FlowReasonDelete is set when flow was evicted by a delete flow mod.
	FlowReasonDelete

	// FlowReasonGroupDelete is set when associated group was removed.
	FlowReasonGroupDelete
)

// FlowRemoved represents an OpenFlow message that is send if the
// controller has requested to be notified when flow entries are timed
// out or are deleted from tables.
type FlowRemoved struct {
	// The Cookie is an opaque data value chosen by the controller.
	Cookie uint64

	// The Priority indicates priority within the specified flow table.
	Priority uint16

	// The Reason specifies the reason of the flow entry removal.
	Reason FlowRemovedReason

	// Table is an identifier of the table.
	Table Table

	// DurationSec is a time flow was alive in seconds.
	DurationSec uint32

	// DurationNSec is a time flow was alive in nanoseconds beyond
	// DurationSec.
	DurationNSec uint32

	// The IdleTimeout specifies time before discarding a flow entry
	// (in seconds).
	IdleTimeout uint16

	// HardTimeout specifies max time before discarding a flow entry (in
	// seconds).
	HardTimeout uint16

	// PacketCount specifies a count of packets have been matched the
	// removed flow entry.
	PacketCount uint64

	// ByteCount specifies a count of packets in bytes have been matched
	// the removed flow entry.
	ByteCount uint64

	// Match lists fields to match.
	Match Match
}

// Cookies implements CookieJar interface.
func (f *FlowRemoved) Cookies() uint64 {
	return f.Cookie
}

// SetCookies implements CookieJar interface.
func (f *FlowRemoved) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the flow
// removed message into the wire format with necessary padding.
func (f *FlowRemoved) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Cookie, f.Priority, f.Reason,
		f.Table, f.DurationSec, f.DurationNSec, f.IdleTimeout,
		f.HardTimeout, f.PacketCount, f.ByteCount, &f.Match,
	)
}

// ReadFrom implements ReaderFrom interface. It serializes the flow
// removed message from the wire format.
func (f *FlowRemoved) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.Cookie, &f.Priority, &f.Reason,
		&f.Table, &f.DurationSec, &f.DurationNSec, &f.IdleTimeout,
		&f.HardTimeout, &f.PacketCount, &f.ByteCount, &f.Match,
	)
}

// FlowStatsRequest is a multipart request used to retrieve information
// about individual flow entries.
//
// For example, to retrieve information about the flow entries that are
// matching the second ingress port, the following request could be
// sent:
//
//	body := &ofp.FlowStatsRequest{
//		Table:    ofp.TableAll,
//		OutPort:  ofp.PortAny,
//		OutGroup: ofp.GropAny,
//		Match:    ofputil.ExtendedMatch(
//			ofputil.MatchInPort(2),
//		),
//	}
//
//	req := ofp.NewMultipartRequest(
//		ofp.MultipartTypeFlow, body)
type FlowStatsRequest struct {
	// Table is an identifier of the table to read or TableAll for
	// inspect all tables of the datapath.
	Table Table

	// Require matching entries to include this as an output port.
	// A value PortAny indicates no restrictions.
	OutPort PortNo

	// Require matching entries to include this as an output group.
	// A value GroupAny indicates no restrictions.
	OutGroup Group

	// Require matching entries to contain this cookie value.
	Cookie uint64

	// Mask used to restrict the cookie bits that must match. A zero
	// value indicates no restrictions.
	CookieMask uint64

	// Fields to match.
	Match Match
}

// Cookies implements openflow.CookieJar interface. It returns the
// cookie assigned to the flow statistic request.
func (f *FlowStatsRequest) Cookies() uint64 {
	return f.Cookie
}

// SetCookies implements openflow.CookieJar interface. It sets the
// specified cookie to the flow statistics request.
func (f *FlowStatsRequest) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the flow
// statistics request into the wire format.
func (f *FlowStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Table, pad3{}, f.OutPort,
		f.OutGroup, pad4{}, f.Cookie, f.CookieMask, &f.Match,
	)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the flow
// statistics request from the wire format.
func (f *FlowStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.Table, &defaultPad3, &f.OutPort,
		&f.OutGroup, &defaultPad4, &f.Cookie, &f.CookieMask, &f.Match,
	)
}

// FlowStats is a body returned within the multipart flow statistics
// reply.
type FlowStats struct {
	// Idenfitier of the table this statistics came from.
	Table Table

	// DurationSec is a time flow has been alive in seconds.
	DurationSec uint32

	// DurationNSec is a time flow has been alive in nanoseconds
	// beyond DurationSec.
	DurationNSec uint32

	// Priority of the entry.
	Priority uint16

	// IdleTimeout is a number of seconds IDLE before expiration.
	IdleTimeout uint16

	// HardTimeout is a number of seconds before expiration.
	HardTimeout uint16

	// Flags configured for the returned flow entry.
	Flags FlowModFlag

	// Opaque controller-issued identifier.
	Cookie uint64

	// Number of packets in the flow.
	PacketCount uint64

	// Number of bytes in the flow.
	ByteCount uint64

	// Description of the fields.
	Match Match

	// The set of instructions associated with a flow entry.
	Instructions Instructions
}

// Cookies implements openflow.CookieJar interface. It returns the
// assigned cookies to the flow statistics.
func (f *FlowStats) Cookies() uint64 {
	return f.Cookie
}

// SetCookies implements openflow.CookieJar interface. It sets the
// specified cookies into the flow statistics message.
func (f *FlowStats) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the flow
// statistics into the wire format.
func (f *FlowStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := encoding.WriteTo(&buf, f.Table, pad1{}, f.DurationSec,
		f.DurationNSec, f.Priority, f.IdleTimeout, f.HardTimeout,
		f.Flags, pad4{}, f.Cookie, f.PacketCount, f.ByteCount,
		&f.Match, &f.Instructions,
	)

	if err != nil {
		return 0, err
	}

	return encoding.WriteTo(w, uint16(buf.Len()+2), buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the flow statistics from the wire format.
func (f *FlowStats) ReadFrom(r io.Reader) (int64, error) {
	var len uint16

	n, err := encoding.ReadFrom(r, &len, &f.Table, &defaultPad1,
		&f.DurationSec, &f.DurationNSec, &f.Priority, &f.IdleTimeout,
		&f.HardTimeout, &f.Flags, &defaultPad4, &f.Cookie, &f.PacketCount,
		&f.ByteCount, &f.Match)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(len)-n)
	f.Instructions = nil

	nn, err := f.Instructions.ReadFrom(limrd)
	return n + nn, err
}
