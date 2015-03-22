package of

import (
	"bytes"
	"net"
	"testing"
)

func TestBytes(t *testing.T) {
	_, netw, err := net.ParseCIDR("10.0.1.1/32")
	if err != nil {
		t.Fatal("Failed to resolve valid ip address")
	}

	b := Bytes(netw.IP)
	if !bytes.Equal(b, []byte{10, 0, 1, 1}) {
		t.Fatal("Failed to marhal bytes:", b)
	}
}
