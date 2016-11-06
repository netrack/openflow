package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestInstructionGotoTable(t *testing.T) {
	tests := []encodingtest.M{
		{&InstructionGotoTable{TableID: 15}, []byte{
			0x00, 0x01, // Instruction type.
			0x00, 0x08, // Instruction length.
			0x0f,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
		}},
	}

	encodingtest.RunM(t, tests)
}

func TestIntructionWriteMetadata(t *testing.T) {
	tests := []encodingtest.MU{
		{&InstructionWriteMetadata{
			Metadata:     0x5091aedc9697445e,
			MetadataMask: 0x3ec894d841073494,
		}, []byte{
			0x00, 0x02, // Instruction type.
			0x00, 0x18, // Intruction length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x50, 0x91, 0xae, 0xdc, 0x96, 0x97, 0x44, 0x5e, // Metadata.
			0x3e, 0xc8, 0x94, 0xd8, 0x41, 0x07, 0x34, 0x94, // Metadata mask.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestInstructionActions(t *testing.T) {
	tests := []encodingtest.M{
		{&InstructionApplyActions{
			Actions: Actions{
				&ActionGroup{GroupID: GroupAll},
				&Action{Type: ActionTypeCopyTTLOut},
			},
		}, []byte{
			0x00, 0x04, // Instruction type.
			0x00, 0x18, // Instruction length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Actions.
			0x00, 0x16, // Action group.
			0x00, 0x08, // Action lenght.
			0xff, 0xff, 0xff, 0xfc,
			0x00, 0xb, // Action type.
			0x00, 0x08, // Action lenght.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunM(t, tests)
}

func TestIntructionMeter(t *testing.T) {
	tests := []encodingtest.M{
		{&InstructionMeter{MeterID: 0x6bb97a25}, []byte{
			0x00, 0x06, // Instruction type.
			0x00, 0x08, // Instruction length.
			0x6b, 0xb9, 0x7a, 0x25, // Meter identifier.
		}},
	}

	encodingtest.RunM(t, tests)
}
