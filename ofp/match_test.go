package ofp

import (
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestXM(t *testing.T) {
	tests := []encodingtest.MU{
		{&XM{Class: XMClassOpenflowBasic,
			Type:  XMTypeUDPSrc,
			Value: XMValue{0x00, 0x35},
			Mask:  XMValue{0xff, 0xff}},
			[]byte{
				0x80, 0x00, // OpenFlow basic.
				0x1f,                   // Match field + Mask flag.
				0x04,                   // Payload length.
				0x00, 0x35, 0xff, 0xff, // Payload.
			}},
	}

	encodingtest.RunMU(t, tests)
}

func TestMatch(t *testing.T) {
	m := &Match{MatchTypeXM, []XM{{
		Class: XMClassOpenflowBasic,
		Type:  XMTypeInPort,
		Value: XMValue{0x00, 0x00, 0x00, 0x03},
	}}}

	tests := []encodingtest.MU{{m, []byte{
		0x00, 0x01, // Match type.
		0x00, 0x0c, // Match length.
		0x80, 0x00, // OpenFlow basic.
		0x00,                   // Match field + Mask flag.
		0x04,                   // Payload length.
		0x00, 0x00, 0x00, 0x03, // Payload.
		0x00, 0x00, 0x00, 0x00, // 4-byte padding.
	}}}

	encodingtest.RunMU(t, tests)
}

func TestXMValue(t *testing.T) {
	value := XMValue{0xef}
	if value.UInt8() != 0xef {
		t.Fatal("Failed to return right uint8 value:", value.UInt8())
	}

	value = XMValue{0x10, 0xab}
	if value.UInt16() != 0x10ab {
		t.Fatal("Failed to return right uin16 value:", value.UInt16())
	}

	value = XMValue{0xde, 0x12, 0x15, 0x70}
	if value.UInt32() != 0xde121570 {
		t.Fatal("Failed to return right uin32 value:", value.UInt32())
	}
}
