package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestSwitchFeatures(t *testing.T) {
	caps := CapabilityFlowStats | CapabilityPortStats

	tests := []encodingtest.MU{
		{&SwitchFeatures{
			DatapathID:   0xd824a830a6001794,
			NumBuffers:   255,
			NumTables:    127,
			AuxiliaryID:  0xfc,
			Capabilities: caps,
			Reserved:     0x0,
		}, []byte{
			0xd8, 0x24, 0xa8, 0x30, 0xa6, 0x00, 0x17, 0x94, // Datapath.
			0x00, 0x00, 0x00, 0xff, // Number of buffers.
			0x7f,       // Number of tables.
			0xfc,       // Auxiliary identifier.
			0x00, 0x00, // 2-byte padding.
			0x00, 0x00, 0x00, 0x05, // Capabilities.
			0x00, 0x00, 0x00, 0x00, // Reserved.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestSwitchConfig(t *testing.T) {
	tests := []encodingtest.MU{
		{&SwitchConfig{
			Flags:          ConfigFlagFragDrop,
			MissSendLength: 65535,
		}, []byte{
			0x00, 0x01, // Flags.
			0xff, 0xff, // Miss send length.
		}},
	}

	encodingtest.RunMU(t, tests)

}
