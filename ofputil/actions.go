package ofputil

import (
	"github.com/netrack/openflow/ofp"
)

// ActionsApply returns a list of instructions with a single element used
// to apply the set of specified actions.
func ActionsApply(actions ...ofp.Action) ofp.Instructions {
	return ofp.Instructions{&ofp.InstructionApplyActions{actions}}
}

// ActionsWrite returns a list of instructions with a single element used
// to write the set of specified actions.
func ActionsWrite(actions ...ofp.Action) ofp.Instructions {
	return ofp.Instructions{&ofp.InstructionWriteActions{actions}}
}

// ActionsClear returns a list of instructions with a single element used
// to clear actions.
func ActionsClear() ofp.Instructions {
	return ofp.Instructions{&ofp.InstructionClearActions{}}
}
