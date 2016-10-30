package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// New flow.
	FlowAdd FlowModCommand = iota

	// Modify all matching flows.
	FlowModity

	// Modify entry strictly matching wildcards and priority.
	FlowModifyStrict

	// Delete all matching flows.
	FlowDelete

	// Delete entry strictly matching wildcards and priority.
	FlowDeleteStrict
)

// FlowModCommand represents a type of the flow table modification
// message.
type FlowModCommand uint8

const (
	// When the FF_SEND_FLOW_REM flag is set, the switch must send a
	// flow removed message when the flow entry expires or is deleted.
	FlowFlagSendFlowRem FlowModFlag = 1 << iota

	// When the FF_CHECK_OVERLAP flag is set, the switch must check that
	// there are no conflicting entries with the same priority prior to
	// inserting it in the flow table.
	//
	// If there is one, the flow mod fails and an error message is
	// returned.
	FlowFlagCheckOverlap FlowModFlag = 1 << iota

	// Reset flow packet and byte counts.
	FlowFlagResetCounts FlowModFlag = 1 << iota

	// When the FF_NO_PKT_COUNTS flag is set, the switch does not need to
	// keep track of the flow packet count.
	FlowFlagNoPktCounts FlowModFlag = 1 << iota

	// When the FF_NO_BYT_COUNTS flag is set, the switch does not need to
	// keep track of the flow byte count.
	FlowFlagNoBytCounts FlowModFlag = 1 << iota
)

type FlowModFlag uint16

// FlowMod represents a modification message to a flow table from the
// controller.
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
	// This field is ignored by FC_ADD messages. A value of 0 indicates
	// no restriction
	CookieMask uint64

	// The TableID is an id of the table to put the flow in.
	//
	// For FC_DELETE_* commands, TT_ALL can also be used to delete matching
	// flows from all tables.
	TableID Table

	// Command specifies a flow modifications command.
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
	// for FC_ADD messages when matching and adding flow entries, and for
	// FC_MODIFY_STRICT or FC_DELETE_STRICT messages when matching flow
	// entries.
	Priority uint16

	// The BufferID refers to a packet buffered at the switch and sent
	// to the controller by a packet-in message.
	//
	// If no buffered packet is associated with the flow mod, it must be
	// set to NO_BUFFER.
	//
	// A flow mod that includes a valid BufferID is effectively equivalent
	// to sending a two-message sequence of a flow mod and a packet-out to
	// P_TABLE, with the requirement that the switch must fully process
	// the flow mod before the packet out.
	BufferID uint32

	// For FC_DELETE* commands, require matching entries to
	// include this as an output port. A value of P_ANY
	// indicates no restriction
	OutPort PortNo

	// For FC_DELETE* commands, require matching entries to
	// include this as an output group. A value of G_ANY
	// indicates no restriction.
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

// Cookies implements CookieJar interface. It returns flow mod message
// cookies.
func (f *FlowMod) Cookies() uint64 {
	return f.Cookie
}

// Cookies implements CookieJar. It sets cookies to flow mod message.
func (f *FlowMod) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

func (f *FlowMod) Bytes() []byte {
	return Bytes(f)
}

// WriteTo implements WriterTo interface.
func (f *FlowMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Cookie, f.CookieMask, f.TableID,
		f.Command, f.IdleTimeout, f.HardTimeout, f.Priority,
		f.BufferID, f.OutPort, f.OutGroup, f.Flags, pad2{},
		&f.Match, &f.Instructions,
	)
}

func (f *FlowMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.Cookie, &f.CookieMask, &f.TableID,
		&f.Command, &f.IdleTimeout, &f.HardTimeout, &f.Priority,
		&f.BufferID, &f.OutPort, &f.OutGroup, &f.Flags, &defaultPad2,
		&f.Match, &f.Instructions,
	)
}

const (
	// Flow idle time exceeded IdleTimeout.
	FlowReasonIdleTimeout FlowRemovedReason = iota

	// Time exceeded HardTimeout.
	FlowReasonHardTimeout

	// Evicted by a delete flow mod.
	FlowReasonDelete

	// Group was removed.
	FlowReasonGroupDelete
)

// FlowRemovedReason specifies the reason of the flow entry removal.
type FlowRemovedReason uint8

// FlowRemoved represents an OpenFlow message that is send if the
// controller has requested to be notified when flow entries time out or
// are deleted from tables.
type FlowRemoved struct {
	// The Cookie is an opaque data value chosen by the controller.
	Cookie uint64

	// The Priority indicates priority within the specified flow table
	// table.
	Priority uint16

	// The Reason stores specifies the reason of the flow entry removal.
	Reason FlowRemovedReason

	// TableID is an id of the table.
	TableID Table

	// DurationSec is a time flow was alive in seconds.
	DurationSec uint32

	// DurationNSec is a time flow was alive in nanoseconds beyond
	// DurationSec.
	DurationNSec uint32

	// The IdleTimeoute specifies time before discarding a flow entry
	// (in seconds).
	IdleTimeout uint16

	// HardTimeout specifis max time before discarding a flow entry (in
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

func (f *FlowRemoved) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Cookie, f.Priority, f.Reason,
		f.TableID, f.DurationSec, f.DurationNSec, f.IdleTimeout,
		f.HardTimeout, f.PacketCount, f.ByteCount, &f.Match,
	)
}

// ReadFrom implements ReaderFrom interface.
func (f *FlowRemoved) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.Cookie, &f.Priority, &f.Reason,
		&f.TableID, &f.DurationSec, &f.DurationNSec, &f.IdleTimeout,
		&f.HardTimeout, &f.PacketCount, &f.ByteCount, &f.Match,
	)
}

type FlowStatsRequest struct {
	TableID Table

	OutPort  PortNo
	OutGroup Group

	Cookie     uint64
	CookieMask uint64
	Match      Match
}

func (f *FlowStatsRequest) Cookies() uint64 {
	return f.Cookie
}

func (f *FlowStatsRequest) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

func (f *FlowStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.TableID, pad3{}, f.OutPort,
		f.OutGroup, pad4{}, f.Cookie, f.CookieMask, &f.Match,
	)
}

func (f *FlowStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.TableID, &defaultPad3, &f.OutPort,
		&f.OutGroup, &defaultPad4, &f.Cookie, &f.CookieMask, &f.Match,
	)
}

type FlowStats struct {
	TableID Table

	DurationSec  uint32
	DurationNSec uint32

	Priority    uint16
	IdleTimeout uint16
	HardTimeout uint16
	Flags       FlowModFlag

	Cookie      uint64
	PacketCount uint64
	ByteCount   uint64

	Match        Match
	Instructions Instructions
}

func (f *FlowStats) Cookies() uint64 {
	return f.Cookie
}

func (f *FlowStats) SetCookies(cookies uint64) {
	f.Cookie = cookies
}

func (f *FlowStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := encoding.WriteTo(&buf, f.TableID, pad1{}, f.DurationSec,
		f.DurationNSec, f.Priority, f.IdleTimeout, f.HardTimeout,
		f.Flags, pad4{}, f.Cookie, f.PacketCount, f.ByteCount,
		&f.Match, &f.Instructions,
	)

	if err != nil {
		return 0, err
	}

	return encoding.WriteTo(w, uint16(buf.Len()), buf.Bytes())
}
