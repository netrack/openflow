package ofp

import (
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

	for _, bucket := range g.Buckets {
		nn, err := bucket.WriteTo(w)
		n += nn

		if err != nil {
			return n, err
		}
	}

	return n, err
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

type GroupStats struct {
	Length       uint16
	_            pad2
	GroupID      Group
	RefCount     uint32
	_            pad4
	PacketCount  uint64
	ByteCount    uint64
	DurationSec  uint32
	DurationNSec uint32
	BucketStats  []BucketCounter
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

type GroupDescStats struct {
	Length  uint16
	Type    GroupType
	_       pad1
	GroupID Group
	Buckets []Bucket
}

const (
	GroupCapabilitySelectWeight   GroupCapability = 1 << iota
	GroupCapabilitySelectLiveness GroupCapability = 1 << iota
	GroupCapabilityChaining       GroupCapability = 1 << iota
	GroupCapabilityChainingChecks GroupCapability = 1 << iota
)

type GroupCapability uint32

type GroupFeatures struct {
	Types        [4]GroupType
	Capabilities GroupCapability
	MaxGroups    []Group
	Actions      []ActionType
}
