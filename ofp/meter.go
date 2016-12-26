package ofp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// MeterCommand represents a type of meter modification message.
type MeterCommand uint16

const (
	// MeterAdd is a command used to add a new meter.
	MeterAdd MeterCommand = iota

	// MeterModify is a command used to modify specified meter.
	MeterModify

	// MeterDelete is a command used to delete specified meter.
	MeterDelete
)

// MeterFlag represents a meter configuration flag.
type MeterFlag uint16

const (
	// MeterFlagKBitPerSec if set, rates value in kilo-bits per second.
	MeterFlagKBitPerSec MeterFlag = 1 << iota

	// MeterFlagPacketPerSec if set, rates value in packets per second.
	MeterFlagPacketPerSec

	// MeterFlagBurst if set, makes burst size.
	MeterFlagBurst

	// MeterFlagStats if set, collects statistics.
	MeterFlagStats
)

// Meter uniquely identifies a meter within a switch. Meters are defined
// starting with 1 up to the maximum number of meters that the switch
// can support. The OpenFlow protocol also defines some virtual meters
// that can not be associated with flows.
type Meter uint32

const (
	// MeterMax defines the last usable meter.
	MeterMax Meter = 0xffff0000

	// MeterSlowpath defines meter for slow datapath.
	MeterSlowpath Meter = 0xfffffffd

	// MeterController defines meter for controller connection.
	MeterController Meter = 0xfffffffe

	// MeterAll represents all meters for statistics request command.
	MeterAll Meter = 0xffffffff
)

// MeterBandType represents a type of meter band.
type MeterBandType uint16

const (
	// MeterBandTypeDrop is used to drop packets.
	MeterBandTypeDrop MeterBandType = 1 + iota

	// MeterBandTypeDSCPRemark is used to remark DSCP in IP header.
	MeterBandTypeDSCPRemark MeterBandType = 1 + iota

	// MeterBandTypeExperimenter is used as an experimenter meter band.
	MeterBandTypeExperimenter MeterBandType = 0xffff
)

var meterBandMap = map[MeterBandType]encoding.ReaderMaker{
	MeterBandTypeDrop:         encoding.ReaderMakerOf(MeterBandDrop{}),
	MeterBandTypeDSCPRemark:   encoding.ReaderMakerOf(MeterBandDSCPRemark{}),
	MeterBandTypeExperimenter: encoding.ReaderMakerOf(MeterBandExperimenter{}),
}

// meterBand is a header of meter band. It holds the type of the meter
// band and its length.
type meterBand struct {
	Type MeterBandType
	Len  uint16
}

// meterBandLen is a length of the meter band.
const meterBandLen = 16

// MeterBand is an interface representing an OpenFlow meter band.
type MeterBand interface {
	encoding.ReadWriter

	// Type returns the type of the meter band.
	Type() MeterBandType
}

// MeterBands group the set of meter bands.
type MeterBands []MeterBand

// WriteTo implements io.WriterTo interface. It serializes the
// list of meter bands into the wire format.
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

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the list of meter bands from the wire format.
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

// MeterBandDrop defines a simple rate limiter that drops packets that
// exceed the band rate value.
type MeterBandDrop struct {
	// Rate for dropping packets.
	Rate uint32

	// BurstSize is a size of bursts.
	BurstSize uint32
}

// Type implements MeterBand interface. It returns the type of meter
// band.
func (m *MeterBandDrop) Type() MeterBandType {
	return MeterBandTypeDrop
}

// WriteTo implements io.WriterTo interface. It serializes the rate
// limiter meter band into the wire format.
func (m *MeterBandDrop) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate, m.BurstSize, pad4{})
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// rate limiter meter band from the wire format.
func (m *MeterBandDrop) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{},
		&m.Rate, &m.BurstSize, &defaultPad4)
}

// MeterBandDSCPRemark defines a simple differentiated services police
// that remark the drop procedure of the DSCP field in the IP header of
// the packets that exceed the band rate value.
type MeterBandDSCPRemark struct {
	// Rate for remarking packets.
	Rate uint32

	// BurstSize is a size of bursts.
	BurstSize uint32

	// PrecLevel indicates by which amount the drop precedence of
	// the packet should be increased if the band is exceeded.
	PrecLevel uint8
}

// Type implements MeterBand interface. It returns the type of the
// meter band.
func (m *MeterBandDSCPRemark) Type() MeterBandType {
	return MeterBandTypeDSCPRemark
}

// WriteTo implements io.WriterTo interface. It serializes the
// DSCP remark meter band into the wire format.
func (m *MeterBandDSCPRemark) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate,
		m.BurstSize, m.PrecLevel, pad3{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// DSCP remark meter band from the wire format.
func (m *MeterBandDSCPRemark) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{},
		&m.Rate, &m.BurstSize, &m.PrecLevel, &defaultPad3)
}

// MeterBandExperimenter defines an experimental meter band.
type MeterBandExperimenter struct {
	// Rate for this band.
	Rate uint32

	// Size of bursts.
	BurstSize uint32

	// Experimenter identifier.
	Experimenter uint32
}

// Type implements MeterBand interface. It returns the type of the
// meter band.
func (m *MeterBandExperimenter) Type() MeterBandType {
	return MeterBandTypeExperimenter
}

// WriteTo implements io.WriterTo interface. It serializes the
// experimental meter band into the wire format.
func (m *MeterBandExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := meterBand{m.Type(), meterBandLen}
	return encoding.WriteTo(w, header, m.Rate,
		m.BurstSize, m.Experimenter)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// experimental meter band from the wire format.
func (m *MeterBandExperimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &meterBand{}, &m.Rate,
		&m.BurstSize, &m.Experimenter)
}

// MeterMod is a message used to modify the meter from the
// controller.
//
// For example, to modify the third meter to drop the packets with a
// rate more than one hundred packets per second, the following request
// can be created:
//
//	mod := &ofp.MeterMod{
//		Command: MeterModify,
//		Flags:   MeterPacketPerSec,
//		Meter:   3,
//		Bands:   MeterBands{
//			&MeterBandDrop{Rate: 100, BurstSize: 150},
//		},
//	}
//
//	req := of.NewRequest(of.TypeMeterMod, mod)
type MeterMod struct {
	// Command specifies a meter modification command.
	Command MeterCommand

	// Flags specifies the meter modification flags.
	Flags MeterFlag

	// Meter is a meter instance.
	Meter Meter

	// Bands is a list of meter bands. It can contain any number of
	// bands, and each band type can be repeated when it make sense.
	//
	// Only a single band is used at a time, if the current rate of
	// packets exceed the rate of multiple bands, the band with the
	// highest configured rate is used.
	Bands MeterBands
}

// WriteTo implements io.WriterTo interface. It serializes the meter
// modification message into the wire format.
func (m *MeterMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Command, m.Flags, m.Meter, m.Bands)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// meter modfiication message from the wire format.
func (m *MeterMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Command, &m.Flags, &m.Meter, m.Bands)
}

// MeterConfigRequest is a multipart request used to retrieve
// configuration for one or more meter.
type MeterConfigRequest struct {
	// Meter instance or MeterAll.
	Meter Meter
}

// WriteTo implements io.WriterTo interface. It serializes the meter
// configuration request into the wire format.
func (m *MeterConfigRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Meter, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// meter configuration request from the wire format.
func (m *MeterConfigRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Meter, &defaultPad4)
}

// meterConfigLen is a length of the meter configuration header
// excluding the length of list of meter bands.
const meterConfigLen = 8

// MeterConfig is meter configuration. This message is returned
// within a body of multipart request.
type MeterConfig struct {
	// Flags is a bitmap of flags that apply.
	Flags MeterFlag

	// Meter instance.
	Meter Meter

	// Bands is a list of associated meter bands.
	Bands MeterBands
}

// WriteTo implements io.WriterTo interface. It serializes the
// meter configuration into the wire format.
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

// ReadFrom implements io.ReadFrom interface. It deserializes the
// meter configuration from the wire format.
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

// MeterFeatures is a message returned within a body of multipart reply.
// It provides the set of features of the metering subsystem.
type MeterFeatures struct {
	// MaxMeter is a maximum number of meters.
	MaxMeter uint32

	// BandTypes is a bitmap of meter band types.
	BandTypes uint32

	// Capabilities is a bitmap of meter flags.
	Capabilities uint32

	// MaxBands is a maximum bands per meters.
	MaxBands uint8

	// MaxColor is a maximum color value.
	MaxColor uint8
}

// WriteTo implements io.WriterTo interface. It serializes the
// meter features into the wire format.
func (m *MeterFeatures) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.MaxMeter, m.BandTypes,
		m.Capabilities, m.MaxBands, m.MaxColor, pad2{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the meter features from the wire format.
func (m *MeterFeatures) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.MaxMeter, &m.BandTypes,
		&m.Capabilities, &m.MaxBands, &m.MaxColor, &defaultPad2)
}

// MeterBandStats consolidates the number of processed bytes and
// packets by the single band.
type MeterBandStats struct {
	// PacketBandCount is a number of packets in band.
	PacketBandCount uint64

	// ByteBandCount is a number of bytes in band.
	ByteBandCount uint64
}

// WriteTo implements io.WriterTo interface. It serializes the meter
// band statistics into the fire format.
func (m *MeterBandStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.PacketBandCount, m.ByteBandCount)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// meter band statistics from the wire format.
func (m *MeterBandStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.PacketBandCount, &m.ByteBandCount)
}

// MeterStatsRequest is a multipart request used to retrieve statistics
// for one or more meters.
type MeterStatsRequest struct {
	// Meter instance or MeterAll.
	Meter Meter
}

// WriteTo implements io.WriterTo interface. It serializes the
// meter statistics request into the wire format.
func (m *MeterStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, m.Meter, pad4{})
}

// ReadFrom implements io.WriterFrom interface. It deserializes the
// meter statistics request from the wire format.
func (m *MeterStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &m.Meter, &defaultPad4)
}

// meterStatsLen is a length of the meter stats message
// without the list of meter band statistic.
const meterStatsLen = 40

// MeterStats is a meter statistics. It is a message returned within
// the body of multipart reply. It provides the set of features of
// the single meter.
type MeterStats struct {
	// Meter is a meter instance.
	Meter Meter

	// FlowCount is a number of flows bound to meter.
	FlowCount uint32

	// PacketInCount is a number of packets in input.
	PacketInCount uint64

	// ByteInCount is a number of bytes in input.
	ByteInCount uint64

	// DurationSec is a time meter has been alive in seconds.
	DurationSec uint32

	// DurationNSec is a time meter has been alive in nanoseconds
	// beyond DurationSec.
	DurationNSec uint32

	// BandStats is a list of meter band statistics.
	BandStats []MeterBandStats
}

// WriteTo implements io.WriterTo interface. It serializes the meter
// statistics into the wire format.
func (m *MeterStats) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of meter band statistics into
	// the temporary buffer to calculate the total
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

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// meter statistics from the wire format.
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
