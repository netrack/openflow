package ofputil

import (
	"github.com/netrack/openflow/ofp"
)

// AsyncConfigMask returns the asynchronous configuration
// mask as a conjunction of master and slave bitmaps.
func AsyncConfigMask(master, slave uint32) [2]uint32 {
	return [2]uint32{master, slave}
}

// PacketInReasonBitmap returns the bitmap of packet-in reasons. The
// result could be used in the asynchronous configuration message.
func PacketInReasonBitmap(reasons ...ofp.PacketInReason) (bits uint32) {
	for _, reason := range reasons {
		bits |= (1 << uint32(reason))
	}

	return
}

// PortReasonBitmap returns the bitmap of port reasons. The result
// could be used in the asynchronous configuration message.
func PortReasonBitmap(reasons ...ofp.PortReason) (bits uint32) {
	for _, reason := range reasons {
		bits |= (1 << uint32(reason))
	}

	return
}

// FlowReasonBitmap returns the bitmap of port reasons. The result
// could be used in the asynchronous configuration message.
func FlowReasonBitmap(reasons ...ofp.FlowRemovedReason) (bits uint32) {
	for _, reason := range reasons {
		bits |= (1 << uint32(reason))
	}

	return
}
