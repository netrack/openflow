package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestPacketIn(t *testing.T) {
	tests := []encodingtest.U{
		{&PacketIn{
			BufferID: NoBuffer,
			Length:   0x20,
			Reason:   PacketInReasonAction,
			TableID:  Table(2),
			Cookie:   0xdeadbeef,
			Match: Match{MatchTypeXM, []XM{{
				Class: XMClassOpenflowBasic,
				Type:  XMTypeInPort,
				Value: XMValue{0x00, 0x00, 0x00, 0x03},
			}}},
		}, []byte{
			0xff, 0xff, 0xff, 0xff, // Buffer identifier.
			0x00, 0x20, // Total frame length.
			0x01, // Packet-in submission reason.
			0x02, // Table identifier.
			0x00, 0x00, 0x00, 0x00,
			0xde, 0xad, 0xbe, 0xef, // Cookie.

			0x00, 0x01, // Match type.
			0x00, 0x0c, // Match length.

			// Match.
			0x80, 0x00, // OpenFlow basic.
			0x00,                   // Match field + Mask flag.
			0x04,                   // Payload length.
			0x00, 0x00, 0x00, 0x03, // Payload.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x00, // 2-byte padding.
		}},
	}

	encodingtest.RunU(t, tests)
}

func TestPacketOut(t *testing.T) {
	tests := []encodingtest.M{
		{&PacketOut{
			BufferID: NoBuffer,
			InPort:   PortController,
			Actions:  Actions{&ActionGroup{GroupID: GroupAll}},
		}, []byte{
			0xff, 0xff, 0xff, 0xff, // Buffer identifier.
			0xff, 0xff, 0xff, 0xfd, // Port number.
			0x00, 0x08, // Actions list length in bytes.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.

			// Actions.
			0x00, 0x16, // Action group.
			0x00, 0x08,
			0xff, 0xff, 0xff, 0xfc,
		}},
	}

	encodingtest.RunM(t, tests)
}
