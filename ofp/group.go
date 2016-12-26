package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// GroupCommand represents a type of the group modification message.
type GroupCommand uint16

const (
	// GroupAdd is a command used to a add a new group.
	GroupAdd GroupCommand = iota

	// GroupModify is a command used to modify all matching groups.
	GroupModify

	// GroupDelete is a command used to delete all matching groups.
	GroupDelete
)

// GroupType represents a type of the group. Values in range [128, 255]
// are reserved for experimental use.
type GroupType uint8

const (
	// GroupTypeAll defines multicast/broadcast group.
	GroupTypeAll GroupType = iota

	// GroupTypeSelect defines a select group.
	GroupTypeSelect

	// GroupTypeIndirect defines an indirect group.
	GroupTypeIndirect

	// GroupTypeFastFailover defines fast failover group.
	GroupTypeFastFailover
)

// Group uniquely identifies the group in the switch.
type Group uint32

const (
	// GroupMax is the last usable group number.
	GroupMax Group = 0xffffff00

	// GroupAll represents all groups for group delete commands.
	GroupAll Group = 0xfffffffc

	// GroupAny is a wildcard group used only for flow stats requests.
	// Selects all flows regardless of group (including flows with no
	// group)
	GroupAny Group = 0xffffffff
)

// GroupMod is a message used to modify the group table from the
// controller.
//
// For example, to create a fast-failover group from first to the second
// port, the following messages can be sent:
//
//	mod := &GroupMod{
//		Command: GroupAdd,
//		Type:    GroupTypeFastFailover,
//		Buckets: []Bucket{{
//			WatchPort: 1, Actions: Actions{
//				&ActionOutput{Port: 2}
//		}}},
//	}
//
//	req := of.NewRequest(of.TypeGroupMod, mod)
type GroupMod struct {
	// Command specified a group modification command.
	Command GroupCommand

	// Type is a type of the group.
	Type GroupType

	// Group identifier.
	Group Group

	// Buckets is an array of buckets. For indirect group type, the
	// array must contain exactly one bucket, other group types may
	// have multiple buckets in the array.
	//
	// For fast failover group, the bucket order does not define the
	// bucket priorities, and bucket order can be changed by modifying
	// the group.
	Buckets []Bucket
}

// WriteTo implements io.WriterTo interface. It serializes the group
// modification message into the wire format.
func (g *GroupMod) WriteTo(w io.Writer) (int64, error) {
	n, err := encoding.WriteTo(w, g.Command, g.Type, pad1{}, g.Group)
	if err != nil {
		return n, err
	}

	nn, err := encoding.WriteSliceTo(w, g.Buckets)
	return n + nn, err
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// groups modification message from the wire format.
func (g *GroupMod) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, &g.Command, &g.Type,
		&defaultPad1, &g.Group)
	if err != nil {
		return n, err
	}

	bucketMaker := encoding.ReaderMakerOf(Bucket{})
	nn, err := encoding.ReadSliceFrom(r, bucketMaker, g.Buckets)
	return n + nn, err
}

// bucketLen is a length of the bucket header, it does not
// include the length of the actions list.
const bucketLen = 16

// Bucket consolidates the parameters and a list of actions that can
// be applied to entering packets. The exact behavior depends on the
// group type.
type Bucket struct {
	// Weight is a relative weight of bucket. Only defined for
	// select groups.
	Weight uint16

	// WatchPort is a port whose state affects whether this bucket
	// is live. Only required for fast failover groups.
	WatchPort PortNo

	// WatchGroup is a group whose state affects whether this bucket
	// is live. Only require for fast failover group.
	WatchGroup Group

	// Actions is an action set associated with the bucket. When
	// bucket is selected for a packet, its action set is applied to
	// the packet.
	Actions Actions
}

// WriteTo implements io.WriterTo interface. It serializes the bucket
// into the wire format.
func (b *Bucket) WriteTo(w io.Writer) (int64, error) {
	// Serialize the list of actions first, to set the
	// valid length into the bucket header.
	actions, err := b.Actions.bytes()
	if err != nil {
		return 0, err
	}

	// The length of the bucket header consist of the
	// header length itself and the length of actions.
	return encoding.WriteTo(w, uint16(bucketLen+len(actions)),
		b.Weight, b.WatchPort, b.WatchGroup, pad4{}, actions)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// bucket from the wire format.
func (b *Bucket) ReadFrom(r io.Reader) (int64, error) {
	// Read the header of the bucket to limit the count
	// of bytes used to unmarshal the list of actions.
	var length uint16
	n, err := encoding.ReadFrom(r, &length, &b.Weight,
		&b.WatchPort, &b.WatchGroup, &defaultPad4)
	if err != nil {
		return n, err
	}

	// Created a limited reader to not read more bytes
	// that it is allocated for the list of actions.
	limrd := io.LimitReader(r, int64(length-bucketLen))
	b.Actions = nil

	nn, err := b.Actions.ReadFrom(limrd)
	return n + nn, err
}

// GroupStatsRequest is a multipart request used to collect statistics
// for one or more groups.
//
// For example, to retrieve the statistics about the first group the
// following request can be created:
//
//	req := ofp.NewMultipartRequest(
//		ofp.MultipartGroup,
//		&ofp.GroupStatsRequest{Group: 1},
//	)
type GroupStatsRequest struct {
	// Group identifier.
	Group Group
}

// WriteTo implements io.WriterTo interface. It serializes the group
// statistics request into the wire format.
func (g *GroupStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, g.Group, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// groups statistics from the wire format.
func (g *GroupStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &g.Group, &defaultPad4)
}

// groupStatsLen is a length of the group statistics header,
// it does not include the length of the bucket counters.
const groupStatsLen = 40

// GroupStats defines the statistics for the single group. An array
// of GroupStats is returned on multipart group request.
type GroupStats struct {
	// Group identifier.
	Group Group

	// RefCount is the number of flows or groups that directly forward
	// to this group.
	RefCount uint32

	// PacketCount is the number of packets processed by this group.
	PacketCount uint64

	// ByteCount is the number of bytes processed by this group.
	ByteCount uint64

	// DurationSec is the time group has been alive in seconds.
	DurationSec uint32

	// DurationNSec is the time group has been alive in nanoseconds
	// beyond DurationSec.
	DurationNSec uint32

	// BucketStats is an array of bucket statistics.
	BucketStats []BucketCounter
}

// WriteTo implements io.WriterTo interface. It serializes the
// group statistics into the wire format.
func (g *GroupStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := encoding.WriteSliceTo(&buf, g.BucketStats)
	if err != nil {
		return 0, err
	}

	return encoding.WriteTo(w, uint16(groupStatsLen+buf.Len()), pad2{},
		g.Group, g.RefCount, pad4{}, g.PacketCount, g.ByteCount,
		g.DurationSec, g.DurationNSec, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the group statistics from the wire format.
func (g *GroupStats) ReadFrom(r io.Reader) (int64, error) {
	var length uint16

	n, err := encoding.ReadFrom(r, &length, &defaultPad2,
		&g.Group, &g.RefCount, &defaultPad4, &g.PacketCount, &g.ByteCount,
		&g.DurationSec, &g.DurationNSec)

	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(length-groupStatsLen))
	counterMaker := encoding.ReaderMakerOf(BucketCounter{})

	nn, err := encoding.ReadSliceFrom(limrd, counterMaker, g.BucketStats)
	return n + nn, err
}

// BucketCounter consolidates the statistics about the single bucket.
type BucketCounter struct {
	// PacketCount is the number of packets processed by bucket.
	PacketCount uint64

	// ByteCount is the number of bytes processed by bucket.
	ByteCount uint64
}

// WriteTo implements io.WriterTo interface. It serializes the bucket
// counter into the wire format.
func (b *BucketCounter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, b.PacketCount, b.ByteCount)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// bucket counter from the wire format.
func (b *BucketCounter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &b.PacketCount, &b.ByteCount)
}

// groupDescStatsLen is the length of the group description
// statistics header (without the length of buckets).
const groupDescStatsLen = 8

// GroupDescStats is a multipart reply that holds a set of groups
// on a switch, along with their corresponding bucket actions.
//
// The respective multipart request body is empty, while the body
// reply is an array of GroupDescStats structures.
type GroupDescStats struct {
	// Type is a group type.
	Type GroupType

	// Group identifier.
	Group Group

	// Buckets is an array of buckets configured on the selected
	// group.
	Buckets []Bucket
}

// WriteTo implements io.WriterTo interface. It serialiezes the
// groups description statistics into the wire format.
func (g *GroupDescStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of buckets to the temporary buffer,
	// so we could set an appropriate length of the message.
	_, err := encoding.WriteSliceTo(&buf, g.Buckets)
	if err != nil {
		return 0, err
	}

	length := uint16(groupDescStatsLen + buf.Len())
	return encoding.WriteTo(w, length, g.Type, pad1{},
		g.Group, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the group description statistics from the wire format.
func (g *GroupDescStats) ReadFrom(r io.Reader) (int64, error) {
	var length uint16

	// Read the header of the messages to retrieve the
	// length of the buckets list.
	n, err := encoding.ReadFrom(r, &length, &g.Type,
		&defaultPad1, &g.Group)

	limrd := io.LimitReader(r, int64(length-groupDescStatsLen))
	bucketMaker := encoding.ReaderMakerOf(Bucket{})

	nn, err := encoding.ReadSliceFrom(limrd, bucketMaker, g.Buckets)
	return n + nn, err
}

// GroupCapability defines the group configuration flags.
type GroupCapability uint32

const (
	// GroupCapabilitySelectWeight if set, means group supports weight
	// (available for select groups).
	GroupCapabilitySelectWeight GroupCapability = 1 << iota

	// GroupCapabilitySelectLiveness if set, means group supports
	// liveness (available for select groups).
	GroupCapabilitySelectLiveness

	// GroupCapabilityChaining if set, means group supports chaining
	// groups.
	GroupCapabilityChaining

	// GroupCapabilityChainingChecks if set, checks chaining for loops
	// and delete.
	GroupCapabilityChainingChecks
)

const (
	// maxGroupsLen is a length of the list with maximum numbers
	// of groups for each type.
	groupFeaturesMaxGroupsLen = 4

	// actionsLen is a length of the groups features actions types.
	groupFeaturesActionsLen = 8
)

// GroupFeatures is a multipart reply that provides the capabilities
// of groups on switch.
//
// The request body is empty, while the reply body is the GroupFeatures
// structure.
//
// For example, to create a new group features message:
//
//	features := &ofp.GroupFeatures{
//		// Bitmap of supported group types.
//		Types: ofputil.GroupBitmap(
//			GroupTypeIndirect,
//			GroupTypeFastFailover),
//
//		// Bitmap of group capabilities.
//		Capabilities: ofputil.GroupCapabilitiesBitmap(
//			GroupCapabilitySelectWeight,
//			GroupCapabilityChaining),
//
//		// Maximum number of groups for each type.
//		MaxGroups: [4]uint32{
//			GroupTypeAll:          255,
//			GroupTypeSelect:       255,
//			GroupTypeIndirect:     127,
//			GroupTypeFastFailover: 127,
//		},
//
//		// Bitmap of supported actions.
//		Actions: ofputil.Bitmap128(
//			ofputil.ActionBitmap(
//				ActionTypeOutput,
//				ActionTypeCopyTTLOut,
//				ActionTypeCopyTTLIn,
//			),
//			0,
//			0,
//			0,
//		),
//	}
type GroupFeatures struct {
	// Types is a bitmap of group types supported.
	Types uint32

	// Capabilities is a bitmap of capabilities supported.
	Capabilities uint32

	// MaxGroups is a maximum number of groups for each type.
	MaxGroups [4]uint32

	// Actions is a bitmap of actions that are supported.
	Actions [4]uint32
}

// WriteTo implemenst io.WriterTo interface. It serializes the
// group features message into the wire format.
func (g *GroupFeatures) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, g.Types, g.Capabilities,
		g.MaxGroups, g.Actions)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the group features message from the wire format.
func (g *GroupFeatures) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &g.Types, &g.Capabilities,
		&g.MaxGroups, &g.Actions)
}
