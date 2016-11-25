package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestTableMod(t *testing.T) {
	tests := []encodingtest.MU{
		{&TableMod{
			Table:  TableMax,
			Config: TableConfigDeprecatedMask,
		}, []byte{
			0xfe,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
			0x00, 0x00, 0x00, 0x03, // Configuration.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTableStats(t *testing.T) {
	tests := []encodingtest.MU{
		{&TableStats{
			Table:        TableMax,
			ActiveCount:  267,
			LookupCount:  132,
			MatchedCount: 54,
		}, []byte{
			0xfe,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
			0x00, 0x00, 0x01, 0x0b, // Active count.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x84, // Lookup count.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x36, // Matched count.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropInstructions(t *testing.T) {
	instr := Instructions{
		&InstructionMeter{Meter: Meter(4)},
		&InstructionGotoTable{Table: Table(15)},
	}

	tests := []encodingtest.MU{
		{&TablePropInstructions{
			Instructions: instr,
		}, []byte{
			0x00, 0x00, // Property type.
			0x00, 0x14, // Property length.

			// Instructions.
			0x00, 0x06, // Instruction type.
			0x00, 0x08, // Instruction length.
			0x00, 0x00, 0x00, 0x04, // Meter identifier.

			0x00, 0x01, // Instruction type.
			0x00, 0x08, // Instruction length.
			0x0f,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.

			// Alignment.
			0x00, 0x00, 0x00, 0x00,
		}},
	}

	gob.Register(InstructionMeter{})
	gob.Register(InstructionGotoTable{})

	encodingtest.RunMU(t, tests)
}

func TestTablePropNextTables(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropNextTables{
			NextTables: []Table{1, 2, 3},
		}, []byte{
			0x00, 0x02, // Property type.
			0x00, 0x07, // Property length.

			// Next tables.
			0x01, 0x02, 0x03,

			// Alignment.
			0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropWriteActions(t *testing.T) {
	actions := Actions{
		&ActionCopyTTLOut{},
		&ActionCopyTTLIn{},
	}

	tests := []encodingtest.MU{
		{&TablePropWriteActions{
			Actions: actions,
		}, []byte{
			0x00, 0x04, // Property type.
			0x00, 0x14, // Property length.

			// Actions.
			0x00, 0xb, // Action type.
			0x00, 0x08, // Action lenght.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			0x00, 0xc, // Action type.
			0x00, 0x08, // Action length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Alignment.
			0x00, 0x00, 0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropApplyActions(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropApplyActions{
			Miss: true,
			Actions: Actions{
				&ActionGroup{Group: GroupAll},
			},
		}, []byte{
			0x00, 0x07, // Property type.
			0x00, 0x0c, // Property length.

			// Actions.
			0x00, 0x16, // Action type.
			0x00, 0x08, // Action length.
			0xff, 0xff, 0xff, 0xfc, // Group identifier.

			// Alignment.
			0x00, 0x00, 0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}

var fields = []XM{
	{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeUDPSrc,
		Value: XMValue{0x00, 0x35},
		Mask:  XMValue{0xff, 0xff},
	},
	{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeInPort,
		Value: XMValue{0x00, 0x00, 0x00, 0x03},
	},
}

var fieldsBytes = []byte{
	0x80, 0x00, // OpenFlow basic.
	0x1f,                   // Match field + Mask flag.
	0x04,                   // Payload length.
	0x00, 0x35, 0xff, 0xff, // Payload.

	0x80, 0x00, // OpenFlow basic.
	0x00,                   // Match field + Mask flag.
	0x04,                   // Payload length.
	0x00, 0x00, 0x00, 0x03, // Payload.
}

func TestTablePropMatch(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropMatch{Fields: fields}, append([]byte{
			0x00, 0x08, // Property type.
			0x00, 0x14, // Property length.
		}, fieldsBytes...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropWildcards(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropWildcards{Fields: fields}, append([]byte{
			0x00, 0x09, // Property type.
			0x00, 0x14, // Property length.
		}, fieldsBytes...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropWriteSetField(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropWriteSetField{Fields: fields}, append([]byte{
			0x00, 0x0a, // Property type.
			0x00, 0x14, // Property length.
		}, fieldsBytes...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropApplySetField(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropApplySetField{
			Miss:   true,
			Fields: fields,
		}, append([]byte{
			0x00, 0x0d, // Property type.
			0x00, 0x14, // Property length.
		}, fieldsBytes...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestTablePropExperimenter(t *testing.T) {
	tests := []encodingtest.MU{
		{&TablePropExperimenter{
			Experimenter: 42,
			ExpType:      43,
			Data:         []byte{0x11, 0x22},
		}, []byte{
			0xff, 0xfe, // Property type.
			0x00, 0x0e, // Property length.
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
			0x00, 0x00, 0x00, 0x2b, // Experimenter type.
			0x11, 0x22, // Experimenter data.

			// Alignment.
			0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTableFeatures(t *testing.T) {
	name := make([]byte, maxTableNameLen)
	copy(name, []byte("table-1"))

	properties := []TableProp{
		&TablePropApplyActions{Actions: Actions{
			&ActionGroup{Group: Group(3)},
		}},
		&TablePropInstructions{Instructions: Instructions{
			&InstructionMeter{Meter: Meter(4)},
		}},
	}

	bytes := []byte{
		0x00, 0x60, // Length.
		0x03,                         // Table identifier.
		0x00, 0x00, 0x00, 0x00, 0x00, // 5-byte padding.
	}

	bytes = append(bytes, name...)
	bytes = append(bytes,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, // Metadata match.
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, // Metadata write.
		0x00, 0x00, 0x00, 0x00, // Table configuration.
		0x00, 0x00, 0x02, 0x00, // Maximum entries.

		// Properties.
		0x00, 0x06, // Property type.
		0x00, 0x0c, // Property length.
		// Actions.
		0x00, 0x16, // Action type.
		0x00, 0x08, // Action length.
		0x00, 0x00, 0x00, 0x03, // Group identi0ier.
		// Alignment.
		0x00, 0x00, 0x00, 0x00,

		0x00, 0x00, // Property type.
		0x00, 0x0c, // Property length.
		// Instructions.
		0x00, 0x06, // Instruction type.
		0x00, 0x08, // Instruction length.
		0x00, 0x00, 0x00, 0x04, // Meter identifier.
		// Alignment.
		0x00, 0x00, 0x00, 0x00,
	)

	tests := []encodingtest.MU{
		{&TableFeatures{
			Table:         Table(3),
			Name:          string(name),
			MetadataMatch: 0xff,
			MetadataWrite: 0xfe,
			MaxEntries:    512,
			Properties:    properties,
		}, bytes},
	}

	gob.Register(TablePropApplyActions{})
	gob.Register(TablePropInstructions{})

	encodingtest.RunMU(t, tests)
}
