package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// New flow.
	FC_ADD FlowModCommand = iota

	// Modify all matching flows.
	FC_MODIFY

	// Modify entry strictly matching wildcards and priority.
	FC_MODIFY_STRICT

	// Delete all matching flows.
	FC_DELETE

	// Delete entry strictly matching wildcards and priority.
	FC_DELETE_STRICT
)

// FlowModCommand represents a type of the flow table modification
// message.
type FlowModCommand uint8

const (
	// When the FF_SEND_FLOW_REM flag is set, the switch must send a
	// flow removed message when the flow entry expires or is deleted.
	FF_SEND_FLOW_REM FlowModFlags = 1 << iota

	// When the FF_CHECK_OVERLAP flag is set, the switch must check that
	// there are no conflicting entries with the same priority prior to
	// inserting it in the flow table.
	//
	// If there is one, the flow mod fails and an error message is
	// returned.
	FF_CHECK_OVERLAP FlowModFlags = 1 << iota

	// Reset flow packet and byte counts.
	FF_RESET_COUNTS FlowModFlags = 1 << iota

	// When the FF_NO_PKT_COUNTS flag is set, the switch does not need to
	// keep track of the flow packet count.
	FF_NO_PKT_COUNTS FlowModFlags = 1 << iota

	// When the FF_NO_BYT_COUNTS flag is set, the switch does not need to
	// keep track of the flow byte count.
	FF_NO_BYT_COUNTS FlowModFlags = 1 << iota
)

type FlowModFlags uint16

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
	Flags FlowModFlags

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
func (f *FlowMod) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	_, err = f.Match.WriteTo(&buf)
	if err != nil {
		return
	}

	_, err = f.Instructions.WriteTo(&buf)
	if err != nil {
		return
	}

	return encoding.WriteTo(w,
		f.Cookie,
		f.CookieMask,
		f.TableID,
		f.Command,
		f.IdleTimeout,
		f.HardTimeout,
		f.Priority,
		f.BufferID,
		f.OutPort,
		f.OutGroup,
		f.Flags,
		pad2{},
		buf.Bytes(),
	)
}

const (
	// Flow idle time exceeded IdleTimeout.
	RR_IDLE_TIMEOUT FlowRemovedReason = iota

	// Time exceeded HardTimeout.
	RR_HARD_TIMEOUT

	// Evicted by a delete flow mod.
	RR_DELETE

	// Group was removed.
	RR_GROUP_DELETE
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

// ReadFrom implements ReaderFrom interface.
func (f *FlowRemoved) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = encoding.ReadFrom(r,
		&f.Cookie,
		&f.Priority,
		&f.Reason,
		&f.TableID,
		&f.DurationSec,
		&f.DurationNSec,
		&f.IdleTimeout,
		&f.HardTimeout,
		&f.PacketCount,
		&f.ByteCount,
	)

	if err != nil {
		return
	}

	nn, err := f.Match.ReadFrom(r)
	return n + nn, err
}

type FlowStatsRequest struct {
	TableID    Table
	_          pad3
	OutPort    PortNo
	OutGroup   Group
	_          pad4
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

type FlowStats struct {
	Length       uint16
	TableID      Table
	_            pad1
	DurationSec  uint32
	DurationNSec uint32

	Priority    uint16
	IdleTimeout uint16
	HardTimeout uint16
	Flags       FlowModFlags
	_           pad4
	Cookie      uint64
	PacketCount uint64
	ByteCount   uint64
	Match       Match
}
