package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	GroupCommandAdd GroupCommand = iota
	GroupCommandModify
	GroupCommandDelete
)

type GroupCommand uint16

const (
	GroupTypeAll GroupType = iota
	GroupTypeSelect
	GroupTypeIndirect
	GroupTypeFastFailover
)

type GroupType uint8

const (
	// Last usable group number
	GroupMax Group = 0xffffff00

	// Represents all groups for group delete commands
	GroupAll Group = 0xfffffffc

	// Wildcard group used only for flow stats requests.
	// Selects all flows regardless of group (including flows with no group)
	GroupAny Group = 0xffffffff
)

type Group uint32

type GroupMod struct {
	Command GroupCommand
	Type    GroupType
	Group   Group
	Buckets []Bucket
}

func (g *GroupMod) WriteTo(w io.Writer) (int64, error) {
	n, err := encoding.WriteTo(w, g.Command, g.Type, pad1{}, g.Group)
	if err != nil {
		return n, err
	}

	nn, err := encoding.WriteSliceTo(w, g.Buckets)
	return n + nn, err
}

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

type Bucket struct {
	Weight     uint16
	WatchPort  PortNo
	WatchGroup Group
	Actions    Actions
}

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
	nn, err := b.Actions.ReadFrom(limrd)
	return n + nn, err
}

type GroupStatsRequest struct {
	Group Group
}

func (g *GroupStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, g.Group, pad4{})
}

func (g *GroupStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &g.Group, &defaultPad4)
}

// groupStatsLen is a length of the group statistics header,
// it does not include the length of the bucket counters.
const groupStatsLen = 40

type GroupStats struct {
	Group        Group
	RefCount     uint32
	PacketCount  uint64
	ByteCount    uint64
	DurationSec  uint32
	DurationNSec uint32
	BucketStats  []BucketCounter
}

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

type BucketCounter struct {
	PacketCount uint64
	ByteCount   uint64
}

func (b *BucketCounter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, b.PacketCount, b.ByteCount)
}

func (b *BucketCounter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &b.PacketCount, &b.ByteCount)
}

const groupDescStatsLen = 8

type GroupDescStats struct {
	Type    GroupType
	Group   Group
	Buckets []Bucket
}

func (g *GroupDescStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of buckets to the temporary buffer,
	// so we could set an appropriate lenght of the message.
	_, err := encoding.WriteSliceTo(&buf, g.Buckets)
	if err != nil {
		return 0, err
	}

	length := uint16(groupDescStatsLen + buf.Len())
	return encoding.WriteTo(w, length, g.Type, pad1{},
		g.Group, buf.Bytes())
}

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

const (
	GroupCapabilitySelectWeight   GroupCapability = 1 << iota
	GroupCapabilitySelectLiveness GroupCapability = 1 << iota
	GroupCapabilityChaining       GroupCapability = 1 << iota
	GroupCapabilityChainingChecks GroupCapability = 1 << iota
)

type GroupCapability uint32

const (
	// groupTypeLen is a length of the group types list
	// of group features message.
	groupFeaturesTypeLen = 4

	// maxGroupsLen is a length of the list with maximum numbers
	// of groups for each type.
	groupFeaturesMaxGroupsLen = 4

	// actionsLen is a length of the groups features actions types.
	groupFeaturesActionsLen = 8
)

type GroupFeatures struct {
	Types        []GroupType
	Capabilities GroupCapability
	MaxGroups    []uint32
	Actions      []ActionType
}

// init allocates the slices of the size (mentioned in the OpenFlow
// protocol specification).
func (g *GroupFeatures) init() ([]GroupType, []uint32, []ActionType) {
	groupTypes := make([]GroupType, groupFeaturesTypeLen)
	maxGroups := make([]uint32, groupFeaturesMaxGroupsLen)
	actions := make([]ActionType, groupFeaturesActionsLen)

	return groupTypes, maxGroups, actions
}

func (g *GroupFeatures) WriteTo(w io.Writer) (int64, error) {
	// For each list of features we will allocate the fixed-length
	// slices and copy the user-defined data into them.
	//
	// This actions is required in order to allow to the user define
	// the features, less than the list could fit. So it is just for
	// the sake of convenience.
	types, groups, actions := g.init()

	copy(types, g.Types)
	copy(groups, g.MaxGroups)
	copy(actions, g.Actions)

	// TODO: probably, we need to generate an error, when the list
	// of features exceed the defined in the protocol length, instead
	// of silently ommiting them.
	return encoding.WriteTo(w, types, g.Capabilities,
		groups, actions)
}

func (g *GroupFeatures) ReadFrom(r io.Reader) (int64, error) {
	// Allocate the memory for fixed-size lists of features,
	// so we could read the complete message.
	g.Types, g.MaxGroups, g.Actions = g.init()

	return encoding.ReadFrom(r, &g.Types, &g.Capabilities,
		&g.MaxGroups, &g.Actions)
}
