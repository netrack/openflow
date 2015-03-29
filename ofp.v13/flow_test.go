package ofp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOXMWrite(t *testing.T) {
}

func TestMatchRead(t *testing.T) {
}

func TestMatchWrite(t *testing.T) {
	var buf bytes.Buffer

	m := Match{MT_OXM, []OXM{{
		Class: XMC_OPENFLOW_BASIC,
		Field: XMT_OFB_IN_PORT,
		Value: []byte{0, 0, 0, 3},
	}}}

	_, err := m.WriteTo(&buf)
	if err != nil {
		t.Fatal("Failed to marshal match:", err)
	}

	hexstr := fmt.Sprintf("%x", buf.Bytes())
	if hexstr != "0001000c800000040000000300000000" {
		t.Fatal("Marshaled match data is incorrect:", hexstr)
	}
}

func TestFlowModWrite(t *testing.T) {
}

func TestOXMValue(t *testing.T) {
	value := OXMValue{0xef}
	if value.UInt8() != 0xef {
		t.Fatal("Failed to return right uint8 value:", value.UInt8())
	}

	value = OXMValue{0x10, 0xab}
	if value.UInt16() != 0x10ab {
		t.Fatal("Failed to return right uin16 value:", value.UInt16())
	}

	value = OXMValue{0xde, 0x12, 0x15, 0x70}
	if value.UInt32() != 0xde121570 {
		t.Fatal("Failed to return right uin32 value:", value.UInt32())
	}
}
