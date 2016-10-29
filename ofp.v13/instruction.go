package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// Setup the next table in the lookup pipeline.
	IT_GOTO_TABLE InstructionType = 1 + iota

	// Setup the metadata field for use later in pipeline.
	IT_WRITE_METADATA InstructionType = 1 + iota

	// Write the action(s) onto the datapath action set.
	IT_WRITE_ACTIONS InstructionType = 1 + iota

	// Applies the action(s) immediately.
	IT_APPLY_ACTIONS InstructionType = 1 + iota

	// Clears all actions from the datapath action set.
	IT_CLEAR_ACTIONS InstructionType = 1 + iota

	// Apply meter (rate limiter).
	IT_METER InstructionType = 1 + iota

	// Experimenter instruction.
	IT_EXPERIMENTER InstructionType = 0xffff
)

// InstructionType represents a type of the flow modification instruction.
type InstructionType uint16

// Instruction header that is common to all instructions. The length
// includes the header and any padding used to make the instruction
// 64-bit aligned.
//
// NB: The length of an instruction *must* always be a multiple of eight.
type instructionhdr struct {
	// Type is an instruction type.
	Type InstructionType

	// Length of this structure in bytes.
	Len uint16
}

type instruction interface {
	io.WriterTo
}

type Instructions []instruction

// WriteTo Implements io.WriterTo interface.
func (i Instructions) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, inst := range i {
		_, err = inst.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	return encoding.WriteTo(w, buf.Bytes())
}

// InstructionGotoTable represents a packet processing pipeline
// redirection message.
type InstructionGotoTable struct {
	// TableID indicates the next table in the packet processing
	// pipeline.
	TableID Table
}

// WriteTo implements WriterTo interface.
func (i InstructionGotoTable) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instructionhdr{
		IT_GOTO_TABLE, 8}, i.TableID, pad3{})
}

// InstructionWriteMetadata setups metadata fields to use later in
// pipeline.
//
// Metadata for the next table lookup can be written using the Metadata and
// the MetadataMask in order to set specific bits on the match field.
//
// If this instruction is not specified, the metadata is passed, unchanged.
type InstructionWriteMetadata struct {
	// Metadata stores a value to write.
	Metadata uint64

	// MetadataMask specifies a metadata bit mask.
	MetadataMask uint64
}

// WriteTo implements WriterTo interface.
func (i InstructionWriteMetadata) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		instructionhdr{IT_WRITE_METADATA, 24},
		pad4{}, i.Metadata, i.MetadataMask)
}

// InstructionActions represents a bundle of action instructions.
//
// For the Apply-Actions instruction, the actions field is treated as a
// list and the actions are applied to the packet in-order.
//
// For the Write-Actions instruction, the actions field is treated as a set
// and the actions are merged into the current action set.
//
// For the Clear-Actions instruction, the structure does not contain any
// actions.
type InstructionActions struct {
	// Type specifies a type of the instruction.
	Type InstructionType

	// Actions associated with IT_WRITE_ACTIONS and IT_APPLY_ACTIONS.
	Actions Actions
}

// WriteTo implements WriterTo interface.
func (i InstructionActions) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	_, err = i.Actions.WriteTo(&buf)
	if err != nil {
		return
	}

	return encoding.WriteTo(w, instructionhdr{
		i.Type, uint16(buf.Len()) + 8}, pad4{}, buf.Bytes())
}

// Instruction structure for IT_METER
type InstructionMeter struct {
	// MeterID indicates which meter to apply on the packet.
	MeterID uint32
}

// WriteTo implements WriterTo interface.
func (i *InstructionMeter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instructionhdr{IT_METER, 8}, i.MeterID)
}
