package ofp

import (
	"bytes"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestDescription(t *testing.T) {
	mk := func(s string, length int) []byte {
		b := make([]byte, length)
		copy(b, []byte(s))
		return b
	}

	mfr := mk("switch manufacturer", descLen)
	hw := mk("switch hardware", descLen)
	sw := mk("switch software", descLen)
	sn := mk("switch serial num", serialNumLen)
	dp := mk("switch datapath", descLen)

	tests := []encodingtest.MU{
		{ReadWriter: &Description{
			Manufacturer: string(mfr),
			Hardware:     string(hw),
			Software:     string(sw),
			SerialNum:    string(sn),
			Datapath:     string(dp),
		}, Bytes: bytes.Join([][]byte{mfr, hw, sw, sn, dp}, []byte(nil))},
	}

	encodingtest.RunMU(t, tests)
}
