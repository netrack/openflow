package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestHello(t *testing.T) {
	elems := HelloElems{&HelloElemVersionBitmap{
		[]uint32{0x10, 0x13, 0x14, 0x15}},
	}

	tests := []encodingtest.MU{
		{ReadWriter: &Hello{}, Bytes: []byte{}},
		{ReadWriter: &Hello{elems}, Bytes: []byte{
			0x00, 0x01, // Hello element type.
			0x00, 0x18, // Hello element length.
			0x00, 0x00, 0x00, 0x10, // OpenFlow versions.
			0x00, 0x00, 0x00, 0x13,
			0x00, 0x00, 0x00, 0x14,
			0x00, 0x00, 0x00, 0x15,
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	gob.Register(HelloElemVersionBitmap{})
	encodingtest.RunMU(t, tests)
}

func TestExperimenter(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &Experimenter{
			Experimenter: 42,
			ExpType:      43,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
			0x00, 0x00, 0x00, 0x2b, // Experimenter type.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestRoleRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &RoleRequest{
			Role:         ControllerRoleMaster,
			GenerationID: 0x22e92b72b39cab3a,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x02, // Controller role.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x22, 0xe9, 0x2b, 0x72, 0xb3, 0x9c, 0xab, 0x3a, // Generation.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestAsyncConfig(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &AsyncConfig{
			[2]uint32{1 << PacketInReasonAction, 0},
			[2]uint32{0, 1 << PortReasonModify},
			[2]uint32{0, 1 << FlowReasonGroupDelete},
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, // Packet-in.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, // Port.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Flow.
		}},
	}

	encodingtest.RunMU(t, tests)
}
