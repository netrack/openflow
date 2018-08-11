package ofp

import (
	"bytes"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestMultipartRequest(t *testing.T) {
	data := []byte{
		0x00, 0x03, // Multipart type.
		0x00, 0x01, // Multipart flags.
		0x00, 0x00, 0x00, 0x00, // 4-byte padding.
	}

	mreq := &MultipartRequest{
		Type:  MultipartTypeTable,
		Flags: MultipartRequestMode,
	}

	// Multipart request saves the rest of the body,
	// so we cannot use the unmarshal test for it.
	encodingtest.RunM(t, []encodingtest.M{{Writer: mreq, Bytes:data}})

	// Test unmarshalling manually:
	var req MultipartRequest
	n, err := req.ReadFrom(bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("failed to unmarshal request: %s", err)
	}

	if n != int64(len(data)) {
		t.Errorf("invalid length of the data retrieved")
	}

	// Validate attributes of the unmarshaled request.
	if req.Type != mreq.Type {
		t.Errorf("invalid type of the request: '%d' is "+
			"not equal to '%d'", req.Type, mreq.Type)
	}

	if req.Flags != mreq.Flags {
		t.Errorf("invalid flags of the request: '%d' are"+
			"not equal to '%d'", req.Flags, mreq.Flags)
	}
}

func TestMultipartReply(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &MultipartReply{
			Type:  MultipartTypePortStats,
			Flags: MultipartReplyMode,
		}, Bytes: []byte{
			0x00, 0x04, // Multipart type.
			0x00, 0x01, // Multipart flags.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestExperimenterMultipartHeader(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &ExperimenterMultipartHeader{
			Experimenter: 42,
			ExpType:      14,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
			0x00, 0x00, 0x00, 0x0e, // Experimenter type.
		}},
	}

	encodingtest.RunMU(t, tests)
}
