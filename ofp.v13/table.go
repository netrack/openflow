package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

const maxTableNameLen = 32

const (
	// TableMax defines the last usable table number.
	TableMax Table = 0xfe

	// TableAll defines the fake table.
	TableAll Table = 0xff
)

type Table uint8

const (
	TableConfigDeprecatedMask TableConfig = 3
)

type TableConfig uint32

// Configure/Modify behavior of a flow table
type TableMod struct {
	// The table_id chooses the table to which the configuration
	// change should be applied. If the TableID is OFPTT_ALL,
	// the configuration is applied to all tables in the switch.
	Table Table

	// The config field is a bitmap that is provided for backward
	// compatibility with earlier version of the specification,
	// it is reserved for future use.
	Config TableConfig
}

func (t *TableMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, t.Table, pad3{}, t.Config)
}

func (t *TableMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &t.Table, &pad3{}, &t.Config)
}

// Information about tables is requested with the MP_TABLE multipart
// request type. The request does not contain any data in the body.
// The body of the reply consists of an array of the TableStats
type TableStats struct {
	// Identifier of table. Lower numbered tables are consulted first
	Table Table

	// Number of active entries
	ActiveCount uint32

	// Number of packets looked up in table
	LookupCount uint64

	// Number of packets that hit table
	MatchedCount uint64
}

func (t *TableStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, t.Table, pad3{},
		t.ActiveCount, t.LookupCount, t.MatchedCount)
}

func (t *TableStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &t.Table, &defaultPad3,
		&t.ActiveCount, &t.LookupCount, &t.MatchedCount)
}

const tableFeaturesLen = 64

type TableFeatures struct {
	Table Table
	Name  string

	MetadataMatch uint64
	MetadataWrite uint64
	Config        TableConfig

	MaxEntries uint32
	Properties []TableProp
}

func (t *TableFeatures) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	for _, prop := range t.Properties {
		_, err := prop.WriteTo(&buf)
		if err != nil {
			return 0, err
		}
	}

	// Copy the table name into the fixed-length slice.
	name := make([]byte, maxTableNameLen)
	copy(name, []byte(t.Name))

	length := tableFeaturesLen + buf.Len()

	return encoding.WriteTo(w, uint16(length), t.Table, pad5{},
		name, t.MetadataMatch, t.MetadataWrite, t.Config,
		t.MaxEntries, buf.Bytes())
}

func (t *TableFeatures) String() string {
	var str string

	for _, prop := range t.Properties {
		str = fmt.Sprintf("%s <%v:%v>", str, prop.Type(), prop)
	}

	return str
}

func (t *TableFeatures) ReadFrom(r io.Reader) (int64, error) {
	var name [maxTableNameLen]byte
	var length uint16

	n, err := encoding.ReadFrom(r, &length, &t.Table,
		&defaultPad5, &name, &t.MetadataMatch, &t.MetadataWrite,
		&t.Config, &t.MaxEntries)

	if err != nil {
		return n, err
	}

	t.Name = string(name[:])
	t.Properties = nil

	var tablePropType TablePropType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := tablePropMap[tablePropType]; ok {
			rd, err := rm.MakeReader()
			t.Properties = append(t.Properties, rd.(TableProp))
			return rd, err
		}

		format := "ofp: unknown table property type: '%x'"
		return nil, fmt.Errorf(format, tablePropType)
	}

	nn, err := encoding.ScanFrom(r, &tablePropType,
		encoding.ReaderMakerFunc(rm))

	return n + nn, err
}

const (
	TablePropTypeInstructions TablePropType = iota
	TablePropTypeInstructionsMiss
	TablePropTypeNextTables
	TablePropTypeNextTablesMiss
	TablePropTypeWriteActions
	TablePropTypeWriteActionsMiss
	TablePropTypeApplyActions
	TablePropTypeApplyActionsMiss
	TablePropTypeMatch
	TablePropTypeWildcards
	TablePropTypeWriteSetField
	TablePropTypeWriteSetFieldMiss
	TablePropTypeApplySetField
	TablePropTypeApplySetFieldMiss
	TablePropTypeExperimenter     TablePropType = 0xfffe
	TablePropTypeExperimenterMiss TablePropType = 0xffff
)

type TablePropType uint16

var tablePropMap = map[TablePropType]encoding.ReaderMaker{
	TablePropTypeInstructions:      encoding.ReaderMakerOf(TablePropInstructions{}),
	TablePropTypeInstructionsMiss:  encoding.ReaderMakerOf(TablePropInstructions{}),
	TablePropTypeNextTables:        encoding.ReaderMakerOf(TablePropNextTables{}),
	TablePropTypeNextTablesMiss:    encoding.ReaderMakerOf(TablePropNextTables{}),
	TablePropTypeWriteActions:      encoding.ReaderMakerOf(TablePropWriteActions{}),
	TablePropTypeWriteActionsMiss:  encoding.ReaderMakerOf(TablePropWriteActions{}),
	TablePropTypeApplyActions:      encoding.ReaderMakerOf(TablePropApplyActions{}),
	TablePropTypeApplyActionsMiss:  encoding.ReaderMakerOf(TablePropApplyActions{}),
	TablePropTypeMatch:             encoding.ReaderMakerOf(TablePropMatch{}),
	TablePropTypeWildcards:         encoding.ReaderMakerOf(TablePropWildcards{}),
	TablePropTypeWriteSetField:     encoding.ReaderMakerOf(TablePropWriteSetField{}),
	TablePropTypeWriteSetFieldMiss: encoding.ReaderMakerOf(TablePropWriteSetField{}),
	TablePropTypeApplySetField:     encoding.ReaderMakerOf(TablePropApplySetField{}),
	TablePropTypeApplySetFieldMiss: encoding.ReaderMakerOf(TablePropApplySetField{}),
	TablePropTypeExperimenter:      encoding.ReaderMakerOf(TablePropExperimenter{}),
	TablePropTypeExperimenterMiss:  encoding.ReaderMakerOf(TablePropExperimenter{}),
}

type TableProp interface {
	encoding.ReadWriter

	Type() TablePropType
}

type tableProp struct {
	Type TablePropType
	Len  uint16
}

// tablePropLen defines the length of the table feature property
// header, it includes only the type and length of the message.
const tablePropLen = 4

// tablePropType returns the type of the table feature property, based
// on the miss flag. If the specified flag is set to true, then miss
// table feature property type will be returned.
func tablePropType(miss bool, rt, mt TablePropType) TablePropType {
	if miss {
		return mt
	}

	return rt
}

func writeTablePropXM(w io.Writer, tp TableProp, xms []XM) (int64, error) {
	var buf bytes.Buffer
	_, err := encoding.WriteSliceTo(&buf, xms)
	if err != nil {
		return 0, err
	}

	header := tableProp{tp.Type(), uint16(tablePropLen + buf.Len())}
	return encoding.WriteTo(w, header, buf.Bytes())
}

func readTablePropXM(r io.Reader, xms *[]XM) (tableProp, int64, error) {
	var header tableProp
	*xms = nil

	n, err := encoding.ReadFrom(r, &header)
	if err != nil {
		return header, n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-tablePropLen))
	nn, err := readAllXM(limrd, xms)

	return header, n + nn, err
}

func writeTablePropActions(w io.Writer, tp TableProp, a Actions) (int64, error) {
	// Write the list of actions into the temporary buffer
	// to be able cacluate the total message length.
	buf, err := a.bytes()
	if err != nil {
		return 0, err
	}

	header := tableProp{tp.Type(), uint16(tablePropLen + len(buf))}
	padding := make([]byte, (header.Len+7)/8*8-header.Len)

	return encoding.WriteTo(w, header, buf, padding)
}

func readTablePropActions(r io.Reader, a *Actions) (tableProp, int64, error) {
	var header tableProp

	n, err := encoding.ReadFrom(r, &header)
	if err != nil {
		return header, n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-tablePropLen))
	*a = nil

	nn, err := a.ReadFrom(limrd)
	n += nn

	if err != nil {
		return header, n, err
	}

	padding := make([]byte, (header.Len+7)/8*8-header.Len)
	nn, err = encoding.ReadFrom(r, padding)

	return header, n + nn, err
}

type TablePropInstructions struct {
	Miss         bool
	Instructions Instructions
}

func (t *TablePropInstructions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeInstructions,
		TablePropTypeInstructionsMiss)
}

func (t *TablePropInstructions) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of instructions into the temporary buffer
	// so we could provide the total message length in the header.
	_, err := t.Instructions.WriteTo(&buf)
	if err != nil {
		return 0, err
	}

	header := tableProp{t.Type(), uint16(tablePropLen + buf.Len())}
	padding := make([]byte, (header.Len+7)/8*8-header.Len)

	return encoding.WriteTo(w, header, buf.Bytes(), padding)
}

func (t *TablePropInstructions) ReadFrom(r io.Reader) (int64, error) {
	var header tableProp
	n, err := encoding.ReadFrom(r, &header)
	if err != nil {
		return n, err
	}

	t.Instructions = nil
	limrd := io.LimitReader(r, int64(header.Len-tablePropLen))
	nn, err := t.Instructions.ReadFrom(limrd)
	n += nn

	if err != nil {
		return n, err
	}

	// If the unmarshaled property describes miss flow-entry
	// we will assign the respective flag in the structure.
	t.Miss = header.Type == TablePropTypeInstructionsMiss
	if err != nil {
		return n, err
	}

	padding := make([]byte, (header.Len+7)/8*8-header.Len)
	nn, err = encoding.ReadFrom(r, padding)

	return n + nn, err
}

type TablePropNextTables struct {
	Miss       bool
	NextTables []Table
}

func (t *TablePropNextTables) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeNextTables,
		TablePropTypeNextTablesMiss)
}

func (t *TablePropNextTables) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Write the list of table identifier to the temorary buffer
	// to calculate the totale length of the message.
	_, err := encoding.WriteTo(&buf, t.NextTables)
	if err != nil {
		return 0, err
	}

	header := tableProp{t.Type(), uint16(tablePropLen + buf.Len())}
	padding := make([]byte, (header.Len+7)/8*8-header.Len)

	return encoding.WriteTo(w, header, buf.Bytes(), padding)
}

func (t *TablePropNextTables) ReadFrom(r io.Reader) (int64, error) {
	var header tableProp

	// Read the header, so we could create the list of the
	// table identifiers of the required size.
	n, err := encoding.ReadFrom(r, &header)
	t.Miss = header.Type == TablePropTypeNextTablesMiss

	if err != nil {
		return 0, err
	}

	padding := make([]byte, (header.Len+7)/8*8-header.Len)
	t.NextTables = make([]Table, header.Len-tablePropLen)

	nn, err := encoding.ReadFrom(r, &t.NextTables, padding)
	return n + nn, err
}

type TablePropWriteActions struct {
	Miss    bool
	Actions Actions
}

func (t *TablePropWriteActions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeWriteActions,
		TablePropTypeWriteActionsMiss)
}

func (t *TablePropWriteActions) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropActions(w, t, t.Actions)
}

func (t *TablePropWriteActions) ReadFrom(r io.Reader) (int64, error) {
	header, n, err := readTablePropActions(r, &t.Actions)
	t.Miss = header.Type == TablePropTypeWriteActionsMiss
	return n, err
}

type TablePropApplyActions struct {
	Miss    bool
	Actions Actions
}

func (t *TablePropApplyActions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeApplyActions,
		TablePropTypeApplyActionsMiss)
}

func (t *TablePropApplyActions) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropActions(w, t, t.Actions)
}

func (t *TablePropApplyActions) ReadFrom(r io.Reader) (int64, error) {
	header, n, err := readTablePropActions(r, &t.Actions)
	t.Miss = header.Type == TablePropTypeApplyActionsMiss
	return n, err
}

type TablePropMatch struct {
	Fields []XM
}

func (t *TablePropMatch) Type() TablePropType {
	return TablePropTypeMatch
}

func (t *TablePropMatch) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

func (t *TablePropMatch) ReadFrom(r io.Reader) (int64, error) {
	_, n, err := readTablePropXM(r, &t.Fields)
	return n, err
}

type TablePropWildcards struct {
	Fields []XM
}

func (t *TablePropWildcards) Type() TablePropType {
	return TablePropTypeWildcards
}

func (t *TablePropWildcards) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

func (t *TablePropWildcards) ReadFrom(r io.Reader) (int64, error) {
	_, n, err := readTablePropXM(r, &t.Fields)
	return n, err
}

type TablePropWriteSetField struct {
	Miss   bool
	Fields []XM
}

func (t *TablePropWriteSetField) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeWriteSetField,
		TablePropTypeWriteSetFieldMiss)
}

func (t *TablePropWriteSetField) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

func (t *TablePropWriteSetField) ReadFrom(r io.Reader) (int64, error) {
	header, n, err := readTablePropXM(r, &t.Fields)
	t.Miss = header.Type == TablePropTypeWriteSetFieldMiss
	return n, err
}

type TablePropApplySetField struct {
	Miss   bool
	Fields []XM
}

func (t *TablePropApplySetField) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeApplySetField,
		TablePropTypeApplySetFieldMiss)
}

func (t *TablePropApplySetField) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

func (t *TablePropApplySetField) ReadFrom(r io.Reader) (int64, error) {
	header, n, err := readTablePropXM(r, &t.Fields)
	t.Miss = header.Type == TablePropTypeApplySetFieldMiss
	return n, err
}

type TablePropExperimenter struct {
	Miss         bool
	Experimenter uint32
	ExpType      uint32
	Data         []byte
}

func (t *TablePropExperimenter) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeExperimenter,
		TablePropTypeExperimenterMiss)
}

func (t *TablePropExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := tableProp{t.Type(), uint16(tablePropLen + len(t.Data) + 8)}
	padding := make([]byte, (header.Len+7)/8*8-header.Len)

	return encoding.WriteTo(w, header, t.Experimenter,
		t.ExpType, t.Data, padding)
}

func (t *TablePropExperimenter) ReadFrom(r io.Reader) (int64, error) {
	var header tableProp

	n, err := encoding.ReadFrom(r, &header, &t.Experimenter, &t.ExpType)
	if err != nil {
		return n, err
	}

	limrd := io.LimitReader(r, int64(header.Len-tablePropLen-8))
	t.Data, err = ioutil.ReadAll(limrd)
	n += int64(len(t.Data))

	if err != nil {
		return n, err
	}

	padding := make([]byte, (header.Len+7)/8*8-header.Len)
	nn, err := encoding.ReadFrom(r, padding)

	return n + nn, err
}
