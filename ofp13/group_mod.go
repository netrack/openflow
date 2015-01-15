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

type GroupMod struct {
	Header  Header
	Command GroupModCommand
	Type    GroupType
	GroupId uint32
	Buckets []Bucket
}

type Bucket struct {
	Length     uint16
	Weight     uint16
	WatchPort  uint32
	WatchGroup uint32
	Actions    []ActionHeader
}
