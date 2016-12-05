package ofp

import (
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestError(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	tests := []encodingtest.MU{
		{&Error{
			Type: ErrTypePortModFailed,
			Code: ErrCodePortModFailedBadPort,
			Data: data,
		}, append([]byte{
			0x00, 0x07, // Error type.
			0x00, 0x00, // Error code.
		}, data...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestErrorExperimenterMsg(t *testing.T) {
	data := []byte{0x03, 0x02, 0x01, 0x00}
	tests := []encodingtest.MU{
		{&ErrorExperimenter{
			ExpType:      4,
			Experimenter: 42,
			Data:         data,
		}, append([]byte{
			0xff, 0xff, // Error type.
			0x00, 0x04, // Experimenter type.
			0x00, 0x00, 0x00, 0x2a, // Experimenter.
		}, data...)},
	}

	encodingtest.RunMU(t, tests)
}
