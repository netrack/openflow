package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
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

func TestGroupMod(t *testing.T) {
	actions1 := Actions{&ActionCopyTTLIn{}}
	actions2 := Actions{&ActionCopyTTLOut{}}

	buckets := []Bucket{
		{10, PortNo(2), Group(5), actions1},
		{20, PortNo(3), Group(6), actions2},
	}

	tests := []encodingtest.MU{
		{&GroupMod{
			Command: GroupModify,
			Type:    GroupTypeIndirect,
			Group:   Group(3),
			Buckets: buckets,
		}, []byte{
			0x00, 0x01, // Group command.
			0x02,                   // Group type.
			0x00,                   // 1-byte padding.
			0x00, 0x00, 0x00, 0x03, // Group identifier.

			// Buckets.
			0x00, 0x18, // Length.
			0x00, 0x0a, // Weight.
			0x00, 0x00, 0x00, 0x02, // Watch port.
			0x00, 0x00, 0x00, 0x05, // Watch group.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Actions.
			0x00, 0x0c, // Copy TTL in.
			0x00, 0x08, // Action length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			0x00, 0x18, // Length.
			0x00, 0x14, // Weight.
			0x00, 0x00, 0x00, 0x03, // Watch port.
			0x00, 0x00, 0x00, 0x06, // Watch group.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Actions.
			0x00, 0x0b, // Copy TTL out.
			0x00, 0x08, // Action length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	gob.Register(ActionCopyTTLOut{})
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

func TestGroupStatsRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{&GroupStatsRequest{Group(7)}, []byte{
			0x00, 0x00, 0x00, 0x07, // Group identifier.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestGroupDescStats(t *testing.T) {
	actions := Actions{&ActionCopyTTLIn{}}
	buckets := []Bucket{{1, PortNo(2), Group(42), actions}}

	tests := []encodingtest.MU{
		{&GroupDescStats{
			Type:    GroupTypeSelect,
			Group:   Group(42),
			Buckets: buckets,
		}, []byte{
			0x00, 0x20, // Length.
			0x01,                   // Group type.
			0x00,                   // 1-byte padding.
			0x00, 0x00, 0x00, 0x2a, // Group identifier.

			// Buckets.
			0x00, 0x18, // Length.
			0x00, 0x01, // Weight.
			0x00, 0x00, 0x00, 0x02, // Watch port.
			0x00, 0x00, 0x00, 0x2a, // Watch group.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Actions.
			0x00, 0x0c, // Copy TTL in.
			0x00, 0x08, // Action length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestGroupStats(t *testing.T) {
	counters := []BucketCounter{{
		PacketCount: 16024454786759363784,
		ByteCount:   9638326023866500405,
	}}

	tests := []encodingtest.MU{
		{&GroupStats{
			Group:        Group(4),
			RefCount:     7,
			PacketCount:  15724173823290642489,
			ByteCount:    17210585649678657194,
			DurationSec:  2042699544,
			DurationNSec: 2073841368,
			BucketStats:  counters,
		}, []byte{
			0x00, 0x38, // Length.
			0x00, 0x00, // 2-byte padding.
			0x00, 0x00, 0x00, 0x04, // Group identifier.
			0x00, 0x00, 0x00, 0x07, // Reference counter.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0xda, 0x37, 0x7c, 0xc8, 0x33, 0xbc, 0x10, 0x39, // Packet count.
			0xee, 0xd8, 0x48, 0x76, 0x79, 0x88, 0x26, 0xaa, // Byte count.
			0x79, 0xc1, 0x1f, 0x18, // Duration seconds.
			0x7b, 0x9c, 0x4e, 0xd8, // Duration nanoseconds.

			// Bucket counters.
			0xde, 0x62, 0x4c, 0xba, 0x34, 0x19, 0x8c, 0xc8,
			0x85, 0xc2, 0x36, 0x73, 0xe1, 0xf7, 0x45, 0x35,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestGroupFeatures(t *testing.T) {
	one := uint32(1)

	types := one<<uint32(GroupTypeAll) |
		one<<uint32(GroupTypeSelect) |
		one<<uint32(GroupTypeIndirect) |
		one<<uint32(GroupTypeFastFailover)

	capabilities := one<<uint32(GroupCapabilityChaining) |
		uint32(GroupCapabilitySelectWeight)

	actions := one<<uint32(ActionTypePopMPLS) |
		one<<uint32(ActionTypePushMPLS)

	tests := []encodingtest.MU{
		{&GroupFeatures{
			Types:        types,
			Capabilities: capabilities,
			MaxGroups:    [4]uint32{4, 5, 6, 7},
			Actions:      [4]uint32{actions, 0, 0, 0},
		}, []byte{
			0x00, 0x00, 0x00, 0x0f, // Group types.
			0x00, 0x00, 0x00, 0x11, // Capabilities.

			// Maximum groups.
			0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x05,
			0x00, 0x00, 0x00, 0x06,
			0x00, 0x00, 0x00, 0x07,

			// Actions.
			0x00, 0x18, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}
