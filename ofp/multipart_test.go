package ofp

import (
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestMultipartRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{&MultipartRequest{
			Type:  MultipartTypeFlow,
			Flags: MultipartRequestMode,
		}, []byte{
			0x00, 0x01, // Multipart type.
			0x00, 0x01, // Multipart flags.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestMultipartReply(t *testing.T) {
	tests := []encodingtest.MU{
		{&MultipartReply{
			Type:  MultipartTypePortStats,
			Flags: MultipartReplyMode,
		}, []byte{
			0x00, 0x04, // Multipart type.
			0x00, 0x01, // Multipart flags.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestExperimenterMultipartHeader(t *testing.T) {
	tests := []encodingtest.MU{
		{&ExperimenterMultipartHeader{
			Experimenter: 42,
			ExpType:      14,
		}, []byte{
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
			0x00, 0x00, 0x00, 0x0e, // Experimenter type.
		}},
	}

	encodingtest.RunMU(t, tests)
}
