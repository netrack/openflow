package ofp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unsafe"

	"github.com/netrack/openflow/encoding"
)

const (
	// Setup the next table in the lookup pipeline.
	InstructionTypeGotoTable InstructionType = 1 + iota

	// Setup the metadata field for use later in pipeline.
	InstructionTypeWriteMetadata InstructionType = 1 + iota

	// Write the action(s) onto the datapath action set.
	InstructionTypeWriteActions InstructionType = 1 + iota

	// Applies the action(s) immediately.
	InstructionTypeApplyActions InstructionType = 1 + iota

	// Clears all actions from the datapath action set.
	InstructionTypeClearActions InstructionType = 1 + iota

	// Apply meter (rate limiter).
	InstructionTypeMeter InstructionType = 1 + iota

	// Experimenter instruction.
	InstructionTypeExperimenter InstructionType = 0xffff
)

// InstructionType represents a type of the flow modification instruction.
type InstructionType uint16

var instructionMap = map[InstructionType]encoding.ReaderMaker{
	InstructionTypeGotoTable:     encoding.ReaderMakerOf(InstructionGotoTable{}),
	InstructionTypeWriteMetadata: encoding.ReaderMakerOf(InstructionWriteMetadata{}),
	InstructionTypeApplyActions:  encoding.ReaderMakerOf(InstructionApplyActions{}),
	InstructionTypeWriteActions:  encoding.ReaderMakerOf(InstructionWriteActions{}),
	InstructionTypeClearActions:  encoding.ReaderMakerOf(InstructionClearActions{}),
	InstructionTypeMeter:         encoding.ReaderMakerOf(InstructionMeter{}),
}

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

const instructionLen uint16 = 8

type Instruction interface {
	encoding.ReadWriter

	// Type returns the type of the instruction.
	Type() InstructionType
}

type Instructions []Instruction

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

func (i Instructions) ReadFrom(r io.Reader) (n int64, err error) {
	var instType InstructionType
	var num int64

	// Retrieve the size of the instruction type, that preceeds
	// every instruction body.
	typeLen := int(unsafe.Sizeof(instType))

	// To keep the implementation of the instruction unmarshaling
	// consistent with marshaling, we have to put the instruction
	// type back to the reader during unmarshaling of the list of
	// instructions.
	rdbuf := bufio.NewReader(r)

	for {
		typeBuf, err := rdbuf.Peek(typeLen)
		if err != nil {
			return n, encoding.SkipEndOfFile(err)
		}

		// Unmarshal the instruction type from the peeked bytes.
		num, err = encoding.ReadFrom(
			bytes.NewReader(typeBuf), &instType)

		if err != nil {
			return n, err
		}

		maker, ok := instructionMap[instType]
		if !ok {
			format := "ofp: unknown instruction type: '%x'"
			return n, fmt.Errorf(format, instType)
		}

		// It is defined a corresponding instruction factory for
		// each instruction type, so we could parse the raw bytes
		// using the correct implementation of the instruction.
		instruction := maker.MakeReader()

		// Read the corresponding instruction from the binary
		// representation.
		num, err = instruction.ReadFrom(rdbuf)
		n += num

		if err != nil && err != io.EOF {
			return n, err
		}

		i = append(i, instruction.(Instruction))
		if err == io.EOF {
			return n, nil
		}
	}

	return n, err
}

// InstructionGotoTable represents a packet processing pipeline
// redirection message.
type InstructionGotoTable struct {
	// TableID indicates the next table in the packet processing
	// pipeline.
	TableID Table
}

// Type implements Instruction interface and returns the type on
// the instruction.
func (i *InstructionGotoTable) Type() InstructionType {
	return InstructionTypeGotoTable
}

// WriteTo implements WriterTo interface.
func (i *InstructionGotoTable) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instructionhdr{i.Type(), 8}, i.TableID, pad3{})
}

func (i *InstructionGotoTable) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, instructionhdr{}, &i.TableID, &defaultPad3)
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

// Type implements Instruction interface and returns the type of the
// instruction.
func (i *InstructionWriteMetadata) Type() InstructionType {
	return InstructionTypeWriteMetadata
}

// WriteTo implements WriterTo interface.
func (i *InstructionWriteMetadata) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instructionhdr{i.Type(), 24},
		pad4{}, i.Metadata, i.MetadataMask)
}

func (i *InstructionWriteMetadata) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &instructionhdr{},
		&defaultPad4, &i.Metadata, &i.MetadataMask)
}

func writeInstructionActions(w io.Writer, t InstructionType,
	actions Actions) (int64, error) {

	// Covert the list of actions into the slice of bytes,
	// so we could include the length of the actions into
	// the instruction header.
	buf, err := actions.bytes()
	if err != nil {
		return int64(len(buf)), err
	}

	// Write the header of the instruction with the length,
	// that includes the list of instruction actions.
	header := instructionhdr{t, uint16(len(buf)) + instructionLen}
	return encoding.WriteTo(w, header, pad4{}, buf)
}

func readInstructionActions(r io.Reader, actions Actions) (int64, error) {
	var read int64

	// Read the header of the instruction at first to retrieve
	// the size of actions in the packet.
	var header instructionhdr
	num, err := encoding.ReadFrom(r, &header)
	read += num

	if err != nil {
		return read, err
	}

	// Limit the reader to the size of actions, so we could know
	// where is the a border of the message.
	limrd := io.LimitReader(r, int64(header.Len-instructionLen))
	num, err = actions.ReadFrom(limrd)
	read += num

	return read, err
}

// InstructionActions represents a bundle of action instructions.
//
// For the Apply-Actions instruction, the actions field is treated as a
// list and the actions are applied to the packet in-order.
type InstructionApplyActions struct {
	// Actions associated with IT_WRITE_ACTIONS and IT_APPLY_ACTIONS.
	Actions Actions
}

func (i *InstructionApplyActions) Type() InstructionType {
	return InstructionTypeApplyActions
}

// WriteTo implements WriterTo interface.
func (i *InstructionApplyActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), i.Actions)
}

func (i *InstructionApplyActions) ReadFrom(r io.Reader) (int64, error) {
	return readInstructionActions(r, i.Actions)
}

// For the Write-Actions instruction, the actions field is treated as a set
// and the actions are merged into the current action set.
type InstructionWriteActions struct {
	Actions Actions
}

func (i *InstructionWriteActions) Type() InstructionType {
	return InstructionTypeWriteActions
}

func (i *InstructionWriteActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), i.Actions)
}

func (i *InstructionWriteActions) ReadFrom(r io.Reader) (int64, error) {
	return readInstructionActions(r, i.Actions)
}

// For the Clear-Actions instruction, the structure does not contain any
// actions.
type InstructionClearActions struct{}

func (i *InstructionClearActions) Type() InstructionType {
	return InstructionTypeClearActions
}

func (i *InstructionClearActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), nil)
}

func (i *InstructionClearActions) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// Instruction structure for IT_METER
type InstructionMeter struct {
	// MeterID indicates which meter to apply on the packet.
	MeterID uint32
}

// Type implements Instruction interface and returns type of the
// instruction.
func (i *InstructionMeter) Type() InstructionType {
	return InstructionTypeMeter
}

// WriteTo implements WriterTo interface.
func (i *InstructionMeter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instructionhdr{i.Type(), 8}, i.MeterID)
}

func (i *InstructionMeter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &instructionhdr{}, &i.MeterID)
}
