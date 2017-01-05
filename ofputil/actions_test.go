package ofputil

import (
	"reflect"
	"testing"

	"github.com/netrack/openflow/ofp"
)

func TestActionsApply(t *testing.T) {
	ac1 := ofp.ActionCopyTTLIn{}
	ac2 := ofp.ActionDecNetworkTTL{}
	ac3 := ofp.ActionOutput{Port: 2}

	its := ActionsApply(&ac1, &ac2, &ac3)
	if len(its) != 1 {
		t.Fatalf("Expected one instruction in apply actions")
	}

	// Cast the instruction interface to apply action.
	itaa, ok := its[0].(*ofp.InstructionApplyActions)
	if !ok {
		t.Fatalf("Should be an apply action instruction")
	}

	aa := ofp.Actions{&ac1, &ac2, &ac3}
	if !reflect.DeepEqual(itaa.Actions, aa) {
		t.Fatalf("Actions are not the same")
	}
}

func TestActionsWrite(t *testing.T) {
	ac1 := ofp.ActionSetNetworkTTL{128}
	ac2 := ofp.ActionGroup{3}

	its := ActionsWrite(&ac1, &ac2)
	if len(its) != 1 {
		t.Fatalf("Expected one instruction in write actions")
	}

	itwa, ok := its[0].(*ofp.InstructionWriteActions)
	if !ok {
		t.Fatalf("Should be a write action instruction")
	}

	wa := ofp.Actions{&ac1, &ac2}
	if !reflect.DeepEqual(itwa.Actions, wa) {
		t.Fatalf("Actions are not the same")
	}
}

func TestActionsClear(t *testing.T) {
	its := ActionsClear()
	if len(its) != 1 {
		t.Fatalf("Expected one instruction in clear actions")
	}

	_, ok := its[0].(*ofp.InstructionClearActions)
	if !ok {
		t.Fatalf("Should be a clear action instruction")
	}
}
