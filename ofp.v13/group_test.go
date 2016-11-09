package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestBucket(t *testing.T) {
	actions := Actions{
		&ActionCopyTTLIn{},
		&ActionOutput{Port: PortNo(3), MaxLen: 65535},
	}

	tests := []encodingtest.MU{
		{&Bucket{
			Weight:     42,
			WatchPort:  PortNo(5),
			WatchGroup: Group(7),
			Actions:    actions,
		}, []byte{
			0x00, 0x28, // Length.
			0x00, 0x2a, // Wight.
			0x00, 0x00, 0x00, 0x05, // Watch port.
			0x00, 0x00, 0x00, 0x07, // Watch group.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Actions.
			0x00, 0x0c, // Copy TTL in.
			0x00, 0x08, // Action length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			0x00, 0x00, // Output.
			0x00, 0x10, // Action length.
			0x00, 0x00, 0x00, 0x03, // Port number.
			0xff, 0xff, // Maximum length.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding
		}},
	}

	gob.Register(ActionCopyTTLIn{})
	gob.Register(ActionOutput{})
	encodingtest.RunMU(t, tests)
}

func TestBucketCounter(t *testing.T) {
	tests := []encodingtest.MU{
		{&BucketCounter{
			PacketCount: 10838451347809794865,
			ByteCount:   9634678394999596076,
		}, []byte{
			0x96, 0x69, 0xea, 0x13, 0x85, 0x8e, 0x7b, 0x31,
			0x85, 0xb5, 0x40, 0xf4, 0x1b, 0x14, 0xe4, 0x2c,
		}},
		{&BucketCounter{
			PacketCount: 5523660708591555761,
			ByteCount:   14713084686717327527,
		}, []byte{
			0x4c, 0xa7, 0xfc, 0x26, 0x1b, 0x61, 0xd4, 0xb1,
			0xcc, 0x2f, 0x60, 0x91, 0xbe, 0x06, 0x98, 0xa7,
		}},
	}

	encodingtest.RunMU(t, tests)
}
