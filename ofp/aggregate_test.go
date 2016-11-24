package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestAggregateStatsRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{&AggregateStatsRequest{
			Table:      TableMax,
			OutPort:    PortNormal,
			OutGroup:   GroupAll,
			Cookie:     0xaabbccdd,
			CookieMask: 0xff00ff00,
			Match: Match{MatchTypeXM, []XM{{
				Class: XMClassOpenflowBasic,
				Type:  XMTypeInPort,
				Value: XMValue{0x00, 0x00, 0x00, 0x03},
			}}},
		}, []byte{
			0xfe,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
			0xff, 0xff, 0xff, 0xfa, // Out port.
			0xff, 0xff, 0xff, 0xfc, // Out group.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, // Cookie.
			0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, // Cookie mask.

			// Match.
			0x00, 0x01, // Match type.
			0x00, 0x0c, // Match length.
			0x80, 0x00, // OpenFlow basic.
			0x00,                   // Match field + Mask flag.
			0x04,                   // Payload length.
			0x00, 0x00, 0x00, 0x03, // Payload.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestAggregateStats(t *testing.T) {
	tests := []encodingtest.MU{
		{&AggregateStats{
			PacketCount: 0x0906caed7a9289a1,
			ByteCount:   0x202bba4297c31ce6,
			FlowCount:   0xcd348340,
		}, []byte{
			0x09, 0x06, 0xca, 0xed, 0x7a, 0x92, 0x89, 0xa1, // Packet count.
			0x20, 0x2b, 0xba, 0x42, 0x97, 0xc3, 0x1c, 0xe6, // Byte count.
			0xcd, 0x34, 0x83, 0x40, // Flow count.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}
