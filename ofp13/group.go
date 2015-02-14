package ofp13

const (
	GC_ADD GroupModCommand = iota
	GC_MODIFY
	GC_DELETE
)

type GroupModCommand uint16

const (
	GT_ALL GroupType = iota
	GT_SELECT
	GT_INDIRECT
	GT_FF
)

type GroupType uint8

const (
	// Last usable group number
	G_MAX Group = 0xffffff00
	// Represents all groups for group delete commands
	G_ALL Group = 0xfffffffc
	// Wildcard group used only for flow stats requests.
	// Selects all flows regardless of group (including flows with no group)
	G_ANY Group = 0xffffffff
)

type Group uint32

type GroupMod struct {
	Command GroupModCommand
	Type    GroupType
	_       uint8
	GroupID uint32
	Buckets []Bucket
}

type Bucket struct {
	Length     uint16
	Weight     uint16
	WatchPort  uint32
	WatchGroup uint32
	_          pad4
	Actions    []ActionHeader
}

type GroupStatsRequest struct {
	GroupID Group
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
	BucketStast  []BucketCounter
}

type BucketCounter struct {
	PacketCount uint64
	ByteCount   uint64
}

type GroupDescStats struct {
	Length  uint16
	Type    GroupType
	_       pad1
	GroupID Group
	Buckets []Bucket
}

const (
	OFPGFC_SELECT_WEIGHT   GroupCapabilities = 1 << iota
	OFPGFC_SELECT_LIVENESS GroupCapabilities = 1 << iota
	OFPGFC_CHAINING        GroupCapabilities = 1 << iota
	OFPGFC_CHAINING_CHECKS GroupCapabilities = 1 << iota
)

type GroupCapabilities uint32

type GroupFeatures struct {
	Types        []GroupType
	Capabilities GroupCapabilities
	MaxGroups    []Group
	Actions      []ActionType
}
