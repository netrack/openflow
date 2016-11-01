package ofp

import (
	"testing"

	"github.com/netrack/openflow/encoding/encodingtest"
)

func TestTableMod(t *testing.T) {
	tests := []encodingtest.MU{
		{&TableMod{
			TableID: TableMax,
			Config:  TableConfigDeprecatedMask,
		}, []byte{
			0xfe,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
			0x00, 0x00, 0x00, 0x03, // Configuration.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTableStats(t *testing.T) {
	tests := []encodingtest.MU{
		{&TableStats{
			TableID:      TableMax,
			ActiveCount:  267,
			LookupCount:  132,
			MatchedCount: 54,
		}, []byte{
			0xfe,             // Table identifier.
			0x00, 0x00, 0x00, // 3-byte padding.
			0x00, 0x00, 0x01, 0x0b, // Active count.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x84, // Lookup count.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x36, // Matched count.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestTableFeatures(t *testing.T) {
}
