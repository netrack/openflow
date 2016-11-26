package ofp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	MeterCommandAdd MeterCommand = iota
	MeterCommandModify
	MeterCommandDelete
)

type MeterCommand uint16

const (
	MeterFlagKBitPerSec   MeterFlag = 1 << iota
	MeterFlagPacketPerSec MeterFlag = 1 << iota
	MeterFlagBurst        MeterFlag = 1 << iota
	MeterFlagStats        MeterFlag = 1 << iota
)

type MeterFlag uint16

const (
	MeterMax        Meter = 0xffff0000
	MeterSlowpath   Meter = 0xfffffffd
	MeterController Meter = 0xfffffffe
	MeterAll        Meter = 0xffffffff
)

type Meter uint32

const (
	MeterBandTypeDrop         MeterBandType = 1 + iota
	MeterBandTypeDSCPRemark   MeterBandType = 1 + iota
	MeterBandTypeExperimenter MeterBandType = 0xffff
)

type MeterBandType uint16

var meterBandMap = map[MeterBandType]encoding.ReaderMaker{
	MeterBandTypeDrop:         encoding.ReaderMakerOf(MeterBandDrop{}),
	MeterBandTypeDSCPRemark:   encoding.ReaderMakerOf(MeterBandDSCPRemark{}),
	MeterBandTypeExperimenter: encoding.ReaderMakerOf(MeterBandExperimenter{}),
}

type meterBand struct {
	Type MeterBandType
	Len  uint16
}

const meterBandLen = 16

type MeterBand interface {
	encoding.ReadWriter

	Type() MeterBandType
}

type MeterBands []MeterBand

func (m MeterBands) WriteTo(w io.Writer) (int64, error) {
	var n int64

	for _, meter := range m {
		nn, err := meter.WriteTo(w)
		n += nn

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (m MeterBands) ReadFrom(r io.Reader) (int64, error) {
	var meterBandType MeterBandType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := meterBandMap[meterBandType]; ok {
			rd, err := rm.MakeReader()
			m = append(m, rd.(MeterBand))
			return rd, err
		}

		format := "ofp: unknown meter band type: '%x'"
		return nil, fmt.Errorf(format, meterBandType)
	}

	return encoding.ScanFrom(r, &meterBandType,
		encoding.ReaderMakerFunc(rm))
}

type MeterBandDrop struct {
	Rate      uint32
	BurstSize uint32
}

func (m *MeterBandDrop) Type() MeterBandType {
	return MeterBandTypeDrop
}

func (m *MeterBandDrop) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate, m.BurstSize, pad4{})
}

func (m *MeterBandDrop) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{},
		&m.Rate, &m.BurstSize, &defaultPad4)
}

type MeterBandDSCPRemark struct {
	Rate      uint32
	BurstSize uint32
	PrecLevel uint8
}

func (m *MeterBandDSCPRemark) Type() MeterBandType {
	return MeterBandTypeDSCPRemark
}

func (m *MeterBandDSCPRemark) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate,
		m.BurstSize, m.PrecLevel, pad3{})
}

func (m *MeterBandDSCPRemark) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{},
		&m.Rate, &m.BurstSize, &m.PrecLevel, &defaultPad3)
}

type MeterBandExperimenter struct {
	Rate         uint32
	BurstSize    uint32
	Experimenter uint32
}

func (m *MeterBandExperimenter) Type() MeterBandType {
	return MeterBandTypeExperimenter
}

func (m *MeterBandExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate,
		m.BurstSize, m.Experimenter)
}

func (m *MeterBandExperimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{}, &m.Rate,
		&m.BurstSize, &m.Experimenter)
}

type MeterMod struct {
	Command MeterCommand
	Flags   MeterFlag
	Meter   Meter
	Bands   MeterBands
}

func (m *MeterMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Command, m.Flags, m.Meter, m.Bands)
}

func (m *MeterMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Command, &m.Flags, &m.Meter, m.Bands)
}

type MeterConfigRequest struct {
	Meter Meter
}

func (m *MeterConfigRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Meter, pad4{})
}

func (m *MeterConfigRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Meter, &defaultPad4)
}

const meterConfigLen = 8

type MeterConfig struct {
	Flags MeterFlag
	Meter Meter
	Bands MeterBands
}

func (m *MeterConfig) WriteTo(w io.Writer) (int64, error) {
	// Write the list of meter bands to the temporary
	// buffer to calculate the total length of the message.
	var buf bytes.Buffer
	_, err := m.Bands.WriteTo(&buf)
	if err != nil {
		return 0, err
	}

	length := uint16(meterConfigLen + buf.Len())
	return encoding.WriteTo(w, length, m.Flags, m.Meter, buf.Bytes())
}

func (m *MeterConfig) ReadFrom(r io.Reader) (int64, error) {
	var length uint16

	// Read the header of the message, so we could use
	// the rest of the reader to unmarshal list of bands.
	n, err := encoding.ReadFrom(r, &length, &m.Flags, &m.Meter)
	if err != nil {
		return n, err
	}

	// Use the rest of bytes to decode the bands.
	limrd := io.LimitReader(r, int64(length-meterConfigLen))
	nn, err := m.Bands.ReadFrom(limrd)
	return n + nn, err
}

// meterFeaturesBandTypesLen is a length of the list of band
// types bitmap of meter features.
const meterFeaturesBandTypesLen = 2

type MeterFeatures struct {
	MaxMeter     uint32
	BandTypes    []MeterBandType
	Capabilities MeterFlag
	MaxBands     uint8
	MaxColor     uint8
}

func (m *MeterFeatures) init() []MeterBandType {
	return make([]MeterBandType, meterFeaturesBandTypesLen)
}

func (m *MeterFeatures) WriteTo(w io.Writer) (int64, error) {
	types := m.init()
	copy(types, m.BandTypes)

	return encoding.WriteTo(w, m.MaxMeter, types,
		m.Capabilities, m.MaxBands, m.MaxColor, pad2{})
}

func (m *MeterFeatures) ReadFrom(r io.Reader) (int64, error) {
	m.BandTypes = m.init()

	return encoding.ReadFrom(r, &m.MaxMeter, &m.BandTypes,
		&m.Capabilities, &m.MaxBands, &m.MaxColor, &defaultPad2)
}

type MeterBandStats struct {
	PacketBandCount uint64
	ByteBandCount   uint64
}

func (m *MeterBandStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.PacketBandCount, m.ByteBandCount)
}

func (m *MeterBandStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.PacketBandCount, &m.ByteBandCount)
}

// meterStatsLen is a length of the meter stats message
// without the list of meter band statistic.
const meterStatsLen = 40

type MeterStats struct {
	Meter         Meter
	FlowCount     uint32
	PacketInCount uint64
	ByteInCount   uint64
	DurationSec   uint32
	DurationNSec  uint32
	BandStats     []MeterBandStats
}

func (m *MeterStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of meter band statistics into
	// the termporary buffer to calculate the total
	// message length.
	_, err := encoding.WriteSliceTo(&buf, m.BandStats)
	if err != nil {
		return 0, err
	}

	length := uint16(meterStatsLen + buf.Len())

	return encoding.WriteTo(w, m.Meter, length, pad6{},
		m.FlowCount, m.PacketInCount, m.ByteInCount,
		m.DurationSec, m.DurationNSec, buf.Bytes())
}

func (m *MeterStats) ReadFrom(r io.Reader) (int64, error) {
	var length uint16

	// Unmarshal all elements from the reader except
	// of the list of meter band statistics.
	n, err := encoding.ReadFrom(r, &m.Meter, &length,
		&defaultPad6, &m.FlowCount, &m.PacketInCount,
		&m.ByteInCount, &m.DurationSec, &m.DurationNSec)

	limrd := io.LimitReader(r, int64(length-meterStatsLen))
	statsMaker := encoding.ReaderMakerOf(MeterBandStats{})

	nn, err := encoding.ReadSliceFrom(limrd, statsMaker, m.BandStats)
	return n + nn, err
}
