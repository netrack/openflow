package ofp13

const (
	MC_ADD MeterModCommands = iota
	MC_MODIFY
	MC_DELETE
)

type MeterModCommands uint16

const (
	MF_KBPS MeterFlags = iota
	MF_PKTPS
	MF_BURST
	MF_STATS
)

type MeterFlags uint16

const (
	M_MAX        Meter = 0xffff0000
	M_SLOWPATH   Meter = 0xfffffffd
	M_CONTROLLER Meter = 0xfffffffe
	M_ALL        Meter = 0xffffffff
)

type Meter uint32

type MeterMod struct {
	Command MeterModCommands
	Flags   MeterFlags
	MeterID Meter
	Bands   []MeterBandHeader
}

const (
	MBT_DROP         MeterBandType = 1
	MBT_DSCP_REMARK  MeterBandType = 2
	MBT_EXPERIMENTER MeterBandType = 0xFFFF
)

type MeterBandType uint16

type MeterBandHeader struct {
	Type      MeterBandType
	Length    uint16
	Rate      uint32
	BurstSize uint32
}

type MeterBandDrop struct {
	MeterBandHeader
	_ pad4
}

type MeterBandDscpRemart struct {
	MeterBandHeader
	PrecLevel uint8
	_         pad3
}

type MeterBandExperimenter struct {
	MeterBandHeader
	Experimenter uint32
}

type MeterMultipartRequest struct {
	MeterID Meter
	_       pad4
}

type MeterStats struct {
	MeterID       Meter
	Length        uint16
	_             pad6
	FlowCount     uint32
	PacketInCount uint64
	ByteInCount   uint64
	DurationSec   uint32
	DurationNSec  uint32
	BandStats     []MeterBandStats
}

type MeterBandStats struct {
	PacketBandCount uint64
	ByteBandCount   uint64
}

type MeterConfig struct {
	Length  uint16
	Flags   MeterModCommands
	MeterID Meter
	Bands   []MeterBandHeader
}

type MeterFeatures struct {
	MaxMeter     uint32
	BandTypes    MeterBandType
	Capabilities MeterFlags
	MaxBands     uint8
	MaxColor     uint8
	_            pad2
}
