package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestMeterMod(t *testing.T) {
	bands := MeterBands{
		&MeterBandDrop{64, 128},
		&MeterBandDSCPRemark{128, 256, 6},
		&MeterBandExperimenter{256, 512, 42},
	}

	tests := []encodingtest.MU{
		{&MeterMod{
			Command: MeterModify,
			Flags:   MeterFlagStats | MeterFlagBurst,
			Meter:   Meter(42),
			Bands:   bands,
		}, []byte{
			0x00, 0x01, // Meter command.
			0x00, 0x0c, // Flags.
			0x00, 0x00, 0x00, 0x2a, // Meter identifier.

			// Band drop.
			0x00, 0x01, // Meter type.
			0x00, 0x10, // Meter length.
			0x00, 0x00, 0x00, 0x40, // Rate.
			0x00, 0x00, 0x00, 0x80, // Burst size.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Band DSCP remark.
			0x00, 0x02, // Meter type.
			0x00, 0x10, // Meter length.
			0x00, 0x00, 0x00, 0x80, // Rate.
			0x00, 0x00, 0x01, 0x00, // Burst size.
			0x06,             // Precedence level.
			0x00, 0x00, 0x00, // 3-byte padding.

			// Band experimenter.
			0xff, 0xff, // Meter type.
			0x00, 0x10, // Meter length.
			0x00, 0x00, 0x01, 0x00, // Rate.
			0x00, 0x00, 0x02, 0x00, // Burst size.
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
		}},
	}

	gob.Register(MeterBandDrop{})
	gob.Register(MeterBandDSCPRemark{})
	gob.Register(MeterBandExperimenter{})

	encodingtest.RunMU(t, tests)
}

func TestMeterConfigRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{&MeterConfigRequest{Meter(2)}, []byte{
			0x00, 0x00, 0x00, 0x02, // Meter idenfitier.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestMeterConfig(t *testing.T) {
	tests := []encodingtest.MU{
		{&MeterConfig{
			Flags: MeterFlagKBitPerSec | MeterFlagBurst,
			Meter: Meter(42),
			Bands: MeterBands{&MeterBandDrop{64, 128}},
		}, []byte{
			0x00, 0x18, // Length.
			0x00, 0x05, // Flags.
			0x00, 0x00, 0x00, 0x2a, // Meter identifier.

			// Band drop.
			0x00, 0x01, // Meter type.
			0x00, 0x10, // Meter length.
			0x00, 0x00, 0x00, 0x40, // Rate.
			0x00, 0x00, 0x00, 0x80, // Burst size.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestMeterFeatures(t *testing.T) {
	types := uint32(1<<MeterBandTypeDrop) |
		uint32(1<<MeterBandTypeDSCPRemark)

	tests := []encodingtest.MU{
		{&MeterFeatures{
			MaxMeter:     45,
			BandTypes:    types,
			Capabilities: uint32(1 << MeterFlagBurst),
			MaxBands:     128,
			MaxColor:     16,
		}, []byte{
			0x00, 0x00, 0x00, 0x2d, // Max meter.
			0x00, 0x00, 0x00, 0x06, // Band types.
			0x00, 0x00, 0x00, 0x10, // Capabilities.
			0x80,       // Max bands.
			0x10,       // Max color.
			0x00, 0x00, // 2-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestMeterStats(t *testing.T) {
	stats := []MeterBandStats{
		{1413263059007179439, 7830709349700751879},
		{15621460444570393343, 13154395619072477107},
	}

	tests := []encodingtest.MU{
		{&MeterStats{
			Meter:         Meter(42),
			FlowCount:     2716600054,
			PacketInCount: 3600438613393559849,
			ByteInCount:   7110296996057607002,
			DurationSec:   50,
			DurationNSec:  10,
			BandStats:     stats,
		}, []byte{
			0x00, 0x00, 0x00, 0x2a, // Meter identifier.
			0x00, 0x48, // Length.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.
			0xa1, 0xec, 0x06, 0xf6, // Flow count.
			0x31, 0xf7, 0x53, 0xd7, 0xca, 0xec, 0x51, 0x29, // Packet-in count.
			0x62, 0xac, 0xd9, 0x72, 0x29, 0x8a, 0x6b, 0x5a, // Byte-in count.
			0x00, 0x00, 0x00, 0x32, // Duration seconds.
			0x00, 0x00, 0x00, 0x0a, // Duration nanoseconds.

			// Band statistics.
			0x13, 0x9c, 0xeb, 0x43, 0xae, 0x4e, 0x0a, 0xaf, // Packet count.
			0x6c, 0xac, 0x44, 0xaa, 0x28, 0x3e, 0x52, 0x07, // Byte count.

			0xd8, 0xca, 0x93, 0x82, 0x1f, 0x6f, 0x1a, 0xff, // Packet count.
			0xb6, 0x8d, 0xcd, 0x1e, 0xdd, 0xc5, 0x07, 0xb3, // Byte count.
		}},
	}

	encodingtest.RunMU(t, tests)
}
