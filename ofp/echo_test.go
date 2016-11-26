package ofp

import (
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestEchoRequest(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	tests := []encodingtest.MU{
		{&EchoRequest{Data: data}, data},
	}

	encodingtest.RunMU(t, tests)
}

func TestEchoReply(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	tests := []encodingtest.MU{
		{&EchoReply{Data: data}, data},
	}

	encodingtest.RunMU(t, tests)
}
