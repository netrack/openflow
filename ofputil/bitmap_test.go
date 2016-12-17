package ofputil

import (
	"testing"

	"github.com/netrack/openflow/ofp"
)

func TestBitmap64(t *testing.T) {
	mask := Bitmap64(3, 4)
	if mask != [2]uint32{3, 4} {
		t.Fatalf("Invalid mask returned: %v", mask)
	}
}

func TestBitmap128(t *testing.T) {
	mask := Bitmap128(3, 4, 5, 6)
	if mask != [4]uint32{3, 4, 5, 6} {
		t.Fatalf("Invalid mask returned: %v", mask)
	}
}

func TestPacketInReasonBitmap(t *testing.T) {
	bitmap := PacketInReasonBitmap(
		ofp.PacketInReasonAction,
		ofp.PacketInReasonInvalidTTL,
	)

	if bitmap != 0x6 {
		t.Fatalf("Invalid bitmap returned: %x", bitmap)
	}
}

func TestPortReasonBitmap(t *testing.T) {
	bitmap := PortReasonBitmap(
		ofp.PortReasonAdd, ofp.PortReasonDelete)

	if bitmap != 0x3 {
		t.Fatalf("Invalid bitmap returned: %x", bitmap)
	}
}

func TestFlowReasonBitmap(t *testing.T) {
	bitmap := FlowReasonBitmap(
		ofp.FlowReasonDelete,
		ofp.FlowReasonGroupDelete,
	)

	if bitmap != 0xc {
		t.Fatalf("Invalid bitmap returned: %x", bitmap)
	}
}

func TestGroupBitmap(t *testing.T) {
	bitmap := GroupBitmap(
		ofp.GroupTypeSelect,
		ofp.GroupTypeIndirect,
	)

	if bitmap != 0x6 {
		t.Fatalf("Invalid bitmap returned: %x", bitmap)
	}
}

func TestActionBitmap(t *testing.T) {
	bitmap := ActionBitmap(
		ofp.ActionTypeOutput,
		ofp.ActionTypePushVLAN,
		ofp.ActionTypePopVLAN,
	)

	if bitmap != 0x60001 {
		t.Fatalf("Invalid bitmap returned: %x", bitmap)
	}
}
