package ofputil

import (
	"github.com/netrack/openflow/ofp"
)

// Bitmap64 returns a conjunction of two specified quartets.
func Bitmap64(q1, q2 uint32) [2]uint32 {
	return [2]uint32{q1, q2}
}

// Bitmap128 returns a conjunction of four specified quartets.
func Bitmap128(q1, q2, q3, q4 uint32) [4]uint32 {
	return [4]uint32{q1, q2, q3, q4}
}

func shl(mask, bit uint32) uint32 {
	return mask | uint32(1)<<bit
}

// PacketInReasonBitmap returns the bitmap of packet-in reasons. The
// result could be used in the asynchronous configuration message.
func PacketInReasonBitmap(reasons ...ofp.PacketInReason) (bits uint32) {
	for _, reason := range reasons {
		bits = shl(bits, uint32(reason))
	}

	return
}

// PortReasonBitmap returns the bitmap of port reasons. The result
// can be used in the asynchronous configuration message.
func PortReasonBitmap(reasons ...ofp.PortReason) (bits uint32) {
	for _, reason := range reasons {
		bits = shl(bits, uint32(reason))
	}

	return
}

// FlowReasonBitmap returns the bitmap of port reasons. The result
// can be used in the asynchronous configuration message.
func FlowReasonBitmap(reasons ...ofp.FlowRemovedReason) (bits uint32) {
	for _, reason := range reasons {
		bits = shl(bits, uint32(reason))
	}

	return
}

// GroupBitmap returns the bitmap of group types. The result can be
// used in the group features message.
func GroupBitmap(groups ...ofp.GroupType) (bits uint32) {
	for _, group := range groups {
		bits = shl(bits, uint32(group))
	}

	return
}

// ActionBitmap returns the bitmap of action types. The result can
// be used in the group features message.
func ActionBitmap(actions ...ofp.ActionType) (bits uint32) {
	for _, action := range actions {
		bits = shl(bits, uint32(action))
	}

	return
}
