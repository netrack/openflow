package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestFlowMod(t *testing.T) {
	flags := FlowFlagSendFlowRem | FlowFlagCheckOverlap

	match := Match{MatchTypeXM, []XM{{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeInPort,
		Value: XMValue{0x00, 0x00, 0x00, 0x03},
	}}}

	instr := Instructions{&InstructionActions{
		Type: InstructionTypeClearActions,
	}}

	tests := []encodingtest.M{
		{&FlowMod{
			Cookie:       0xdbf7525e57bd7eef,
			CookieMask:   0x44d8b8f011090dcb,
			TableID:      TableMax,
			Command:      FlowAdd,
			IdleTimeout:  45,
			HardTimeout:  90,
			Priority:     10,
			BufferID:     NoBuffer,
			OutPort:      PortFlood,
			OutGroup:     GroupAny,
			Flags:        flags,
			Match:        match,
			Instructions: instr,
		}, []byte{
			0xdb, 0xf7, 0x52, 0x5e, 0x57, 0xbd, 0x7e, 0xef, // Cookie.
			0x44, 0xd8, 0xb8, 0xf0, 0x11, 0x09, 0x0d, 0xcb, // Cookie mask.
			0xfe,       // Table identifier.
			0x00,       // Command.
			0x00, 0x2d, // IDLE timeout.
			0x00, 0x5a, // Hard timeout.
			0x00, 0x0a, // Prioriry.
			0xff, 0xff, 0xff, 0xff, // Buffer identifier.
			0xff, 0xff, 0xff, 0xfb, // Out port.
			0xff, 0xff, 0xff, 0xff, // Out group.
			0x00, 0x03, // Flags.
			0x00, 0x00, // 2-byte padding.

			// Match.
			0x00, 0x01, // Match type.
			0x00, 0x0c, // Match length.
			0x80, 0x00, // OpenFlow basic.
			0x00,                   // Match field + Mask flag.
			0x04,                   // Payload length.
			0x00, 0x00, 0x00, 0x03, // Payload.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Instructions.
			0x00, 0x05, // Instruction type.
			0x00, 0x08, // Intruction length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunM(t, tests)
}

func TestFlowRemoved(t *testing.T) {
	match := Match{MatchTypeXM, []XM{{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeInPort,
		Value: XMValue{0x00, 0x00, 0x00, 0x03},
	}}}

	tests := []encodingtest.MU{
		{&FlowRemoved{
			Cookie:       0xf22884334a8def04,
			Priority:     11,
			Reason:       FlowReasonHardTimeout,
			TableID:      TableMax,
			DurationSec:  929584189,
			DurationNSec: 1244051003,
			IdleTimeout:  46,
			HardTimeout:  91,
			PacketCount:  8005984375916722949,
			ByteCount:    3104105491404993109,
			Match:        match,
		}, []byte{
			0xf2, 0x28, 0x84, 0x33, 0x4a, 0x8d, 0xef, 0x04, // Cookie.
			0x00, 0x0b, // Priority,
			0x01,                   // Reason.
			0xfe,                   // Table identifier.
			0x37, 0x68, 0x54, 0x3d, // Duration seconds.
			0x4a, 0x26, 0xb6, 0x3b, // Duration nanoseconds.
			0x00, 0x2e, // IDLE timeout.
			0x00, 0x5b, // Hard timeout.
			0x6f, 0x1a, 0xf8, 0x5f, 0x53, 0xd7, 0xfb, 0x05, // Packet count.
			0x2b, 0x13, 0xff, 0x7f, 0x88, 0x88, 0xb2, 0x55, // Byte count.

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

func TestFlowStatsRequest(t *testing.T) {
}

func TestFlowStatsReply(t *testing.T) {
}
