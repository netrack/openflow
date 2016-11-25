package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestActionCopyTTLInOut(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionCopyTTLOut{}, []byte{
			0x00, 0xb, // Action type.
			0x00, 0x08, // Action lenght.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
		{&ActionCopyTTLIn{}, []byte{
			0x00, 0xc,
			0x00, 0x08,
			0x00, 0x00, 0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionOutput(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionOutput{Port: PortIn, MaxLen: 0}, []byte{
			0x0, 0x0, // Action type.
			0x0, 0x10, // Action length.
			0xff, 0xff, 0xff, 0xf8, // Port number.
			0x0, 0x0, // Maximum length.
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, // 6-byte padding.
		{&ActionOutput{Port: PortFlood, MaxLen: 0}, []byte{
			0x0, 0x0,
			0x0, 0x10,
			0xff, 0xff, 0xff, 0xfb,
			0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		{&ActionOutput{Port: PortController, MaxLen: 0x80}, []byte{
			0x0, 0x0,
			0x0, 0x10,
			0xff, 0xff, 0xff, 0xfd,
			0x0, 0x80,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionGroup(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionGroup{Group: GroupMax}, []byte{
			0x0, 0x16, // Action type.
			0x0, 0x08, // Action length.
			0xff, 0xff, 0xff, 0x00}}, // Group identifier.
		{&ActionGroup{Group: GroupAll}, []byte{
			0x0, 0x16,
			0x0, 0x08,
			0xff, 0xff, 0xff, 0xfc}},
		{&ActionGroup{Group: GroupAny}, []byte{
			0x0, 0x16,
			0x0, 0x08,
			0xff, 0xff, 0xff, 0xff}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionSetQueue(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionSetQueue{QueueID: QueueAll}, []byte{
			0x0, 0x15, // Action type.
			0x0, 0x08, // Action length.
			0xff, 0xff, 0xff, 0xff}}, // Queue identifier.
		{&ActionSetQueue{QueueID: 0x4200}, []byte{
			0x0, 0x15,
			0x0, 0x08,
			0x0, 0x0, 0x42, 0x00}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionMPLSTTL(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionSetMPLSTTL{TTL: 64}, []byte{
			0x0, 0x0f, // Action type.
			0x0, 0x08, // Action length.
			0x40,            // Time to live.
			0x0, 0x0, 0x0}}, // 3-bytes padding.
		{&ActionSetMPLSTTL{TTL: 32}, []byte{
			0x0, 0x0f,
			0x0, 0x08,
			0x20,
			0x0, 0x0, 0x0}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionSetNetworkTTL(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionSetNetworkTTL{TTL: 48}, []byte{
			0x0, 0x17, // Action type.
			0x0, 0x08, // Action length.
			0x30,            // Time to live.
			0x0, 0x0, 0x0}}, // 3-bytes padding.
	}

	encodingtest.RunMU(t, tests)
}

func TestActionPushPopVLAN(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionPushVLAN{EtherType: 1000}, []byte{
			0x0, 0x11, // Action type.
			0x0, 0x08, // Action length.
			0x03, 0xe8, // Ethernet type.
			0x0, 0x0}}, // 2-bytes padding.
		{&ActionPopVLAN{}, []byte{
			0x0, 0x12,
			0x0, 0x08,
			0x0, 0x0, 0x0, 0x0}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionPopMPLS(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionPopMPLS{EtherType: 1001}, []byte{
			0x0, 0x14, // Action type.
			0x0, 0x08, // Action length.
			0x03, 0xe9, // Ethernet type.
			0x0, 0x0}}, // 2-bytes padding.
		{&ActionPopMPLS{EtherType: 9}, []byte{
			0x0, 0x14,
			0x0, 0x8,
			0x0, 0x9,
			0x0, 0x0}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionSetField(t *testing.T) {
	xm1 := XM{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeInPort,
		Value: XMValue{0x00, 0x01},
		Mask:  XMValue{0x00, 0xff},
	}

	xm2 := XM{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeIPv4Src,
		Value: XMValue{172, 17, 0, 25},
	}

	tests := []encodingtest.MU{
		{&ActionSetField{Field: xm1}, []byte{
			0x00, 0x19, // Action type.
			0x00, 0x10, // Action length.
			0x80, 0x00, // OpenFlow basic.
			0x01,                   // Match field + Mask flag.
			0x04,                   // Payload length.
			0x00, 0x01, 0x00, 0xff, // Payload.
			0x0, 0x0, 0x0, 0x0, // 4-bytes padding.
		}},
		{&ActionSetField{Field: xm2}, []byte{
			0x00, 0x19,
			0x00, 0x10,
			0x80, 0x00,
			0x16,
			0x04,
			0xac, 0x11, 0x00, 0x19,
			0x0, 0x0, 0x0, 0x0,
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestActionExperimenter(t *testing.T) {
	tests := []encodingtest.MU{
		{&ActionExperimenter{41}, []byte{
			0xff, 0x0ff, // Action type.
			0x0, 0x08, // Action length.
			0x0, 0x0, 0x0, 0x29, // Experimeter.
		}},
		{&ActionExperimenter{42}, []byte{
			0xff, 0x0ff,
			0x0, 0x08,
			0x0, 0x0, 0x0, 0x2a,
		}},
	}

	encodingtest.RunMU(t, tests)
}
