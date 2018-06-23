package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestPacketIn(t *testing.T) {
	tests := []encodingtest.MU{
		{&PacketIn{
			Buffer: NoBuffer,
			Length: 0x38,
			Reason: PacketInReasonAction,
			Table:  Table(2),
			Cookie: 0xdeadbeef,
			Match: Match{MatchTypeXM, []XM{{
				Class: XMClassOpenflowBasic,
				Type:  XMTypeInPort,
				Value: XMValue{0x00, 0x00, 0x00, 0x03},
			}}},
			Data: []byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
				0x08, 0x06,
			},
		}, []byte{
			0xff, 0xff, 0xff, 0xff, // Buffer identifier.
			0x00, 0x38, // Total frame length.
			0x01, // Packet-in submission reason.
			0x02, // Table identifier.
			0x00, 0x00, 0x00, 0x00,
			0xde, 0xad, 0xbe, 0xef, // Cookie.

			0x00, 0x01, // Match type.
			0x00, 0x0c, // Match length.

			// Match.
			0x80, 0x00, // OpenFlow basic.
			0x00,                   // Match field + Mask flag.
			0x04,                   // Match field length.
			0x00, 0x00, 0x00, 0x03, // Match field value.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x00, // 2-byte padding.

			// Original ethernet frame.
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // Destination MAC.
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // Source MAC.
			0x08, 0x06, // Ether-Type
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestPacketOut(t *testing.T) {
	tests := []encodingtest.MU{
		{&PacketOut{
			Buffer:  NoBuffer,
			InPort:  PortController,
			Actions: Actions{&ActionGroup{Group: GroupAll}},
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

	gob.Register(ActionGroup{})
	encodingtest.RunMU(t, tests)
}
