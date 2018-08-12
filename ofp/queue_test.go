package ofp

import (
	"encoding/gob"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

func TestPacketQueue(t *testing.T) {
	props := QueueProps{
		&QueuePropMinRate{42},
		&QueuePropMaxRate{43},
	}

	tests := []encodingtest.MU{
		{ReadWriter: &PacketQueue{
			Queue:      QueueAll,
			Port:       PortNormal,
			Properties: props,
		}, Bytes: []byte{
			0xff, 0xff, 0xff, 0xff, // Queue.
			0xff, 0xff, 0xff, 0xfa, // Port number.
			0x00, 0x30, // Length.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.

			// Properties
			0x00, 0x01, // Queue property min rate.
			0x00, 0x10, // Queue property length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x2a, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.

			0x00, 0x02, // Queue property max rate.
			0x00, 0x10, // Queue propertu length.
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x2b, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}},
	}

	gob.Register(QueuePropMinRate{})
	gob.Register(QueuePropMaxRate{})
	encodingtest.RunMU(t, tests)
}

func TestQueuePropExperimenter(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}

	tests := []encodingtest.MU{
		{ReadWriter: &QueuePropExperimenter{
			Experimenter: 359,
			Data:         data,
		}, Bytes: append([]byte{
			0xff, 0xff, // Queue property type.
			0x00, 0x14, // Queue property length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x00, 0x01, 0x67, // Experimenter.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}, data...)},
	}

	encodingtest.RunMU(t, tests)
}

func TestQueueStastsRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &QueueStatsRequest{
			Port:  PortNo(1),
			Queue: Queue(2),
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x01, // Port number.
			0x00, 0x00, 0x00, 0x02, // Queue number.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestQueueStats(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &QueueStats{
			Port:         PortNo(42),
			Queue:        Queue(9),
			TxBytes:      5631918746835501714,
			TxPackets:    13419917410144254707,
			TxErrors:     9021256935579401134,
			DurationSec:  3332165368,
			DurationNSec: 3773668775,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x2a, // Port number.
			0x00, 0x00, 0x00, 0x09, // Queue.
			0x4e, 0x28, 0x98, 0x42, 0xd4, 0xfe, 0x42, 0x92, // Tx bytes.
			0xba, 0x3d, 0x1f, 0xc8, 0x62, 0xba, 0x7a, 0xf3, // Tx packets.
			0x7d, 0x31, 0xf1, 0x62, 0xe0, 0xbd, 0x3b, 0xae, // Tx errors.
			0xc6, 0x9c, 0xce, 0xf8, // Duration seconds.
			0xe0, 0xed, 0x9d, 0xa7, // Duration nano seconds.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestQueueGetConfigRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &QueueGetConfigRequest{PortNo(42)}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x2a, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestQueueGetConfigReply(t *testing.T) {
	props1 := QueueProps{
		&QueuePropMinRate{1},
		&QueuePropMaxRate{2},
	}

	props2 := QueueProps{
		&QueuePropMinRate{3},
		&QueuePropMaxRate{4},
	}

	queues := []PacketQueue{
		{Queue(1), PortNo(4), props1},
		{Queue(2), PortNo(4), props2},
	}

	tests := []encodingtest.MU{
		{ReadWriter: &QueueGetConfigReply{
			Port:   PortNo(43),
			Queues: queues,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x2b, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.

			// Queue #1.
			0x00, 0x00, 0x00, 0x01, // Queue number.
			0x00, 0x00, 0x00, 0x04, // Port number.
			0x00, 0x30, // Length.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.
			// Property #1.
			0x00, 0x01, // Queue property min rate.
			0x00, 0x10, // Queue property length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x01, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.
			// Property #2.
			0x00, 0x02, // Queue property max rate.
			0x00, 0x10, // Queue propertu length.
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x02, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

			// Queue #2.
			0x00, 0x00, 0x00, 0x02, // Queue number.
			0x00, 0x00, 0x00, 0x04, // Port number.
			0x00, 0x30, // Length.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.
			// Property #1.
			0x00, 0x01, // Queue property min rate.
			0x00, 0x10, // Queue property length.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x00, 0x03, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 6-byte padding.
			// Property #2.
			0x00, 0x02, // Queue property max rate.
			0x00, 0x10, // Queue propertu length.
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x04, // Rate.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}},
	}

	encodingtest.RunMU(t, tests)
}
