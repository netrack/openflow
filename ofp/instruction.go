package ofp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	// InstructionTypeGotoTable is used to setup the next table in
	// the lookup pipeline.
	InstructionTypeGotoTable InstructionType = 1 + iota

	// InstructionTypeWriteMetadata is used to setup the metadata
	// field for use later in pipeline.
	InstructionTypeWriteMetadata InstructionType = 1 + iota

	// InstructionTypeWriteActions is used to write the action(s)
	// onto the datapath action set.
	InstructionTypeWriteActions InstructionType = 1 + iota

	// InstructionTypeApplyActions is used to apply the action(s)
	// immediately.
	InstructionTypeApplyActions InstructionType = 1 + iota

	// InstructionTypeClearActions is used to clear all actions from
	// the datapath action set.
	InstructionTypeClearActions InstructionType = 1 + iota

	// InstructionTypeMeter is used to apply meter (rate limiter).
	InstructionTypeMeter InstructionType = 1 + iota

	// InstructionTypeExperimenter is an experimenter instruction.
	InstructionTypeExperimenter InstructionType = 0xffff
)

// InstructionType represents a type of the flow modification instruction.
type InstructionType uint16

// String returns a string representation of instruction type.
func (it InstructionType) String() string {
	text, ok := instructionTypeText[it]
	if !ok {
		return fmt.Sprintf("InstructionType(%d)", it)
	}
	return text
}

var instructionTypeText = map[InstructionType]string{
	InstructionTypeGotoTable:     "InstructionGotoTable",
	InstructionTypeWriteMetadata: "InstructionWriteMetadata",
	InstructionTypeWriteActions:  "InstructionWriteActions",
	InstructionTypeApplyActions:  "InstructionApplyActions",
	InstructionTypeClearActions:  "InstructionClearActions",
	InstructionTypeMeter:         "InstructionMeter",
	InstructionTypeExperimenter:  "InstructionExperimenter",
}

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
type instruction struct {
	// Type is an instruction type.
	Type InstructionType

	// Length of this structure in bytes.
	Len uint16
}

func (i *instruction) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &i.Type, &i.Len)
}

const (
	// instructionLen is a minimum length of the instruction.
	instructionLen uint16 = 8

	// instructionHeaderLen is a length of the instruction header.
	instructionHeaderLen uint16 = 4
)

// Instruction is an interface representing an OpenFlow action.
type Instruction interface {
	encoding.ReadWriter

	// Type returns the type of the instruction.
	Type() InstructionType
}

// Instructions group the set of instructions.
type Instructions []Instruction

// WriteTo implements io.WriterTo interface. It serializes the set of
// instructions into the wire format.
func (i *Instructions) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, inst := range *i {
		_, err = inst.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	return encoding.WriteTo(w, buf.Bytes())
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the set
// of actions from the wire format.
func (i *Instructions) ReadFrom(r io.Reader) (n int64, err error) {
	var instType InstructionType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := instructionMap[instType]; ok {
			rd, err := rm.MakeReader()
			*i = append(*i, rd.(Instruction))
			return rd, err
		}

		return nil, fmt.Errorf("ofp: unknown instruction type: %s", instType)
	}

	return encoding.ScanFrom(r, &instType,
		encoding.ReaderMakerFunc(rm))
}

// InstructionGotoTable represents a packet processing pipeline
// redirection message.
type InstructionGotoTable struct {
	// Table indicates the next table in the packet processing
	// pipeline.
	Table Table
}

// Type implements Instruction interface and returns the type on
// the instruction.
func (i *InstructionGotoTable) Type() InstructionType {
	return InstructionTypeGotoTable
}

// WriteTo implements io.WriterTo interface. It serializes the
// instruction used to redirect the processing pipeline into the
// wire format.
func (i *InstructionGotoTable) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instruction{i.Type(), 8}, i.Table, pad3{})
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// instruction used to redirect the processing pipeline from the wire
// format.
func (i *InstructionGotoTable) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &instruction{}, &i.Table, &defaultPad3)
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

// WriteTo implements io.WriterTo interface. It serializes instruction
// used to write metadata into the wire format.
func (i *InstructionWriteMetadata) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instruction{i.Type(), 24},
		pad4{}, i.Metadata, i.MetadataMask)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes instruction
// used to write metadata from the wire format.
func (i *InstructionWriteMetadata) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &instruction{},
		&defaultPad4, &i.Metadata, &i.MetadataMask)
}

// writeIntructionActions serializes the instruction with action.
// It is shared among the Appply/Clear/Write instructions.
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
	header := instruction{t, uint16(len(buf)) + instructionLen}
	return encoding.WriteTo(w, header, pad4{}, buf)
}

// readInstructionActions deserializes the instruction with actions.
// It is shared among the Apply/Clear/Write instructions.
func readInstructionActions(r io.Reader, actions *Actions) (int64, error) {
	var read int64

	// Read the header of the instruction at first to retrieve
	// the size of actions in the packet.
	var header instruction
	num, err := encoding.ReadFrom(r, &header, &defaultPad4)
	read += num

	if err != nil {
		return read, err
	}

	// Limit the reader to the size of actions, so we could know
	// where is the a border of the message.
	limrd := io.LimitReader(r, int64(header.Len-8))
	num, err = actions.ReadFrom(limrd)
	read += num

	return read, err
}

// InstructionApplyActions is an instruction used to apply the
// list of actions to the processing packet in-order.
type InstructionApplyActions struct {
	// Actions is a list of actions to apply.
	Actions Actions
}

// Type implements Instruction interface and returns the type of the
// instruction.
func (i *InstructionApplyActions) Type() InstructionType {
	return InstructionTypeApplyActions
}

// WriteTo implements io.WriterTo interface. It serializes the
// instruction used to apply actions into the wire format.
func (i *InstructionApplyActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), i.Actions)
}

// ReadFrom implements io.ReadFrom interface. It deserializes
// the instruction used to apply actions from the wire format.
func (i *InstructionApplyActions) ReadFrom(r io.Reader) (int64, error) {
	return readInstructionActions(r, &i.Actions)
}

// InstructionWriteActions represents a bundle of actions that should
// be merged into the current action set.
type InstructionWriteActions struct {
	// Actions is a list of actions to write.
	Actions Actions
}

// Type implements Instruction interface. It returns the type of
// instruction.
func (i *InstructionWriteActions) Type() InstructionType {
	return InstructionTypeWriteActions
}

// WriteTo implements io.WriterTo interface. It serializes instruction
// used to write actions into the wire format.
func (i *InstructionWriteActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), i.Actions)
}

// ReadFrom implements io.ReadFrom interface. It serializes instruction
// used to write actions from the wire format.
func (i *InstructionWriteActions) ReadFrom(r io.Reader) (int64, error) {
	return readInstructionActions(r, &i.Actions)
}

// InstructionClearActions is an instruction used to clear the set
// of actions.
type InstructionClearActions struct{}

// Type implements Instruction interface. It returns the type of the
// instruction.
func (i *InstructionClearActions) Type() InstructionType {
	return InstructionTypeClearActions
}

// WriteTo implements io.WriterTo interface. It serializes the
// instruction used to clear actions into the wire format.
func (i *InstructionClearActions) WriteTo(w io.Writer) (int64, error) {
	return writeInstructionActions(w, i.Type(), nil)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// instruction used to clear actions from the wire format.
func (i *InstructionClearActions) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &defaultPad8)
}

// InstructionMeter is an instruction used to apply meter (rate
// limiter).
type InstructionMeter struct {
	// Meter indicates which meter to apply on the packet.
	Meter Meter
}

// Type implements Instruction interface and returns type of the
// instruction.
func (i *InstructionMeter) Type() InstructionType {
	return InstructionTypeMeter
}

// WriteTo implements io.WriterTo interface. It serializes the
// instruction used to apply meter into the wire format.
func (i *InstructionMeter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, instruction{i.Type(), 8}, i.Meter)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// instruction used to apply meter from the wire format.
func (i *InstructionMeter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &instruction{}, &i.Meter)
}
