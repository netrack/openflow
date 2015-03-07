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
