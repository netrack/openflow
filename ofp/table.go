package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// maxTableNameLen defines the maximum length of the table name.
const maxTableNameLen = 32

// Table defines a switch table number.
type Table uint8

// String returns a string representation of the table.
func (t Table) String() string {
	return fmt.Sprintf("Table(%d)", t)
}

const (
	// TableMax defines the last usable table number.
	TableMax Table = 0xfe

	// TableAll defines the wildcard table used for table config, flow
	// stats and flow deletes.
	TableAll Table = 0xff
)

// TableConfig defines the flags to configure the table. Reserved for
// future use.
type TableConfig uint32

const (
	// TableConfigDeprecatedMask defines the deprecated bits of the
	// table configuration.
	TableConfigDeprecatedMask TableConfig = 3
)

// TableMod is a message used to configure or modify behavior of a
// flow table.
type TableMod struct {
	// The Table chooses the table to which the configuration change should
	// be applied. If the Table is TableAll, the configuration is applied
	// to all tables in the switch.
	Table Table

	// The config field is a bitmap that is provided for backward
	// compatibility with earlier version of the specification, it is
	// reserved for future use.
	Config TableConfig
}

// WriteTo implements io.WriterTo interface. It serializes the table
// modification message into the wire format.
func (t *TableMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, t.Table, pad3{}, t.Config)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// table modification message from the wire format.
func (t *TableMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &t.Table, &pad3{}, &t.Config)
}

// TableStats defines a multipart request body used to query information
// about tables presented within a switch.
type TableStats struct {
	// Table identifies a table within a switch. Lower numbered tables
	// are consulted first.
	Table Table

	// ActiveCount is a number of active entries.
	ActiveCount uint32

	// LookupCount is a number of packets looked up in table.
	LookupCount uint64

	// MatchedCount is a number of packets that hit table.
	MatchedCount uint64
}

// WriteTo implements io.WriterTo interface. It serializes the table
// statistics into the wire format.
func (t *TableStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, t.Table, pad3{},
		t.ActiveCount, t.LookupCount, t.MatchedCount)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// table statistics from the wire format.
func (t *TableStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &t.Table, &defaultPad3,
		&t.ActiveCount, &t.LookupCount, &t.MatchedCount)
}

// tableFeaturesLen defines the length of the table features header.
const tableFeaturesLen = 64

// TableFeatures is a body of multipart request and reply. It is used
// to query for the capabilities of existing table, and to optionally
// ask the switch to reconfigure its tables to match a supplied
// configuration.
//
// In general, the table feature capabilities represents all possible
// features of a table, however some features may be mutually exclusive
// and the current capabilities structures do not allow us to represent
// such exclusions.
type TableFeatures struct {
	// Table identifies a table within a switch.
	Table Table

	// Name is human-readable name of the table.
	Name string

	// MetadataMatch specifies bits of metadata can match.
	MetadataMatch uint64

	// MetadataWrite specifies bits of metadata can write.
	MetadataWrite uint64

	// Config is a bitmap of table configurations.
	Config TableConfig

	// MaxEntries is a maximum number of entries supported.
	MaxEntries uint32

	// Properties is a list of table properties.
	Properties []TableProp
}

// WriteTo implements io.WriterTo interface. It serializes the table
// features into the wire format.
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

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// table features from the wire format.
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

		return nil, fmt.Errorf("ofp: unknown table property type: %s", tablePropType)
	}

	limrd := io.LimitReader(r, int64(length)-n)
	nn, err := encoding.ScanFrom(limrd, &tablePropType,
		encoding.ReaderMakerFunc(rm))

	return n + nn, err
}

// TablePropType defines the table property types.
//
// Low order bit cleared indicates a property for a regular Flow Entry.
// Low order bit set indicates a property for the Table-Miss Flow Entry.
type TablePropType uint16

func (t TablePropType) String() string {
	if str, have := tablePropTypeText[t]; have {
		return str
	}
	return fmt.Sprintf("TablePropType(%d)", t)
}

const (
	// TablePropTypeInstructions indicates instructions property.
	TablePropTypeInstructions TablePropType = iota

	// TablePropTypeInstructionsMiss indicates instructions property for
	// table-miss.
	TablePropTypeInstructionsMiss

	// TablePropTypeNextTables indicates next table property.
	TablePropTypeNextTables

	// TablePropTypeNextTablesMiss indicates next table property for
	// table-miss.
	TablePropTypeNextTablesMiss

	// TablePropTypeWriteActions indicates write actions property.
	TablePropTypeWriteActions

	// TablePropTypeWriteActionsMiss indicates write actions property for
	// table-miss.
	TablePropTypeWriteActionsMiss

	// TablePropTypeApplyActions indicates apply actions property.
	TablePropTypeApplyActions

	// TablePropTypeApplyActionsMiss indicates apply actions property for
	// table-miss.
	TablePropTypeApplyActionsMiss

	// TablePropTypeMatch indicates match property.
	TablePropTypeMatch

	// TablePropTypeWildcards indicates wildcards property.
	TablePropTypeWildcards = 1 + iota

	// TablePropTypeWriteSetField indicates write set-field property.
	TablePropTypeWriteSetField = 2 + iota

	// TablePropTypeWriteSetFieldMiss indicates write set-field property
	// for table-miss.
	TablePropTypeWriteSetFieldMiss

	// TablePropTypeApplySetField indicates apply set-field property.
	TablePropTypeApplySetField

	// TablePropTypeApplySetFieldMiss indicates apply set-field property
	// for table-miss.
	TablePropTypeApplySetFieldMiss

	// TablePropTypeExperimenter indicates experimenter property.
	TablePropTypeExperimenter TablePropType = 0xfffe

	// TablePropTypeExperimenterMiss indicates experimenter property for
	// table-miss.
	TablePropTypeExperimenterMiss TablePropType = 0xffff
)

var tablePropTypeText = map[TablePropType]string{
	TablePropTypeInstructions:      "TablePropTypeInstructions",
	TablePropTypeInstructionsMiss:  "TablePropTypeInstructionsMiss",
	TablePropTypeNextTables:        "TablePropTypeNextTables",
	TablePropTypeNextTablesMiss:    "TablePropTypeNextTablesMiss",
	TablePropTypeWriteActions:      "TablePropTypeWriteActions",
	TablePropTypeWriteActionsMiss:  "TablePropTypeWriteActionsMiss",
	TablePropTypeApplyActions:      "TablePropTypeApplyActions",
	TablePropTypeApplyActionsMiss:  "TablePropTypeApplyActionsMiss",
	TablePropTypeMatch:             "TablePropTypeMatch",
	TablePropTypeWildcards:         "TablePropTypeWildcards",
	TablePropTypeWriteSetField:     "TablePropTypeWriteSetField",
	TablePropTypeWriteSetFieldMiss: "TablePropTypeWriteSetFieldMiss",
	TablePropTypeApplySetField:     "TablePropTypeApplySetField",
	TablePropTypeApplySetFieldMiss: "TablePropTypeApplySetFieldMiss",
	TablePropTypeExperimenter:      "TablePropTypeExperimenter",
	TablePropTypeExperimenterMiss:  "TablePropTypeExperimenterMiss",
}

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

// TableProp is an interface representing OpenFlow table property.
type TableProp interface {
	encoding.ReadWriter

	// Type returns the type of the table property.
	Type() TablePropType
}

// tableProp defines a common header for all table properties.
type tableProp struct {
	Type TablePropType
	Len  uint16
}

func yieldTableProp(r io.Reader, miss *bool) (
	*tableProp, io.Reader, int64, error) {

	header := new(tableProp)
	n, err := encoding.ReadFrom(r, header)
	if err != nil {
		return nil, nil, n, err
	}

	// Set the type of the table property (if it
	// is a miss property or not).
	if miss != nil {
		*miss = (header.Type & 1) == 1
	}

	limrdlen := int64(header.Len - tablePropLen)
	limrd := io.LimitReader(r, limrdlen)
	return header, limrd, n, nil
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
	headerlen := tablePropLen + xmlen*len(xms)
	header := tableProp{tp.Type(), uint16(headerlen)}

	n, err := encoding.WriteTo(w, header)
	if err != nil {
		return n, err
	}

	nn, err := encoding.WriteSliceTo(w, xms)
	if n += nn; err != nil {
		return n, err
	}

	nn, err = encoding.WriteTo(w, makePad(headerlen))
	return n + nn, err
}

func readTablePropXM(r io.Reader, xms *[]XM, miss *bool) (int64, error) {
	header, limrd, n, err := yieldTableProp(r, miss)
	if err != nil {
		return n, err
	}

	*xms = (*xms)[:0]
	nn, err := readAllXM(limrd, xms, false)
	if n += nn; err != nil {
		return n, err
	}

	nn, err = encoding.ReadFrom(r, makePad(int(header.Len)))
	return n + nn, err
}

func writeTablePropActions(w io.Writer, p TableProp, a []ActionType) (int64, error) {
	nacs := len(a)

	proplen := tablePropLen + int(actionHeaderLen)*nacs
	header := tableProp{p.Type(), uint16(proplen)}

	acs := make([]action, nacs)
	for ii, actionType := range a {
		acs[ii] = action{actionType, actionHeaderLen}
	}

	return encoding.WriteTo(w, header, acs, makePad(proplen))
}

// readTablePropActions reads the list of actions from the given reader
// and appends them into the specified list of action types. The list will
// be truncated first.
func readTablePropActions(r io.Reader, a *[]ActionType, m *bool) (int64, error) {
	header, limrd, n, err := yieldTableProp(r, m)
	if err != nil {
		return n, err
	}

	// Truncate any data specified within a list of action types,
	// in this way, list will always contain only decoded messages.
	*a = (*a)[:0]
	amaker := encoding.ReaderMakerOf(action{})

	// Read a list of action headers and aggregate the types.
	join := func(r io.ReaderFrom) { *a = append(*a, r.(*action).Type) }
	nn, err := encoding.ReadFunc(limrd, amaker, join)

	if n += nn; err != nil {
		return n, err
	}

	nn, err = encoding.ReadFrom(r, makePad(int(header.Len)))
	return n + nn, err
}

// TablePropInstructions defines the instructions property of the table.
type TablePropInstructions struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Instructions specifies a list of instructions supported by the
	// table.
	Instructions []InstructionType
}

// String returns a string representation of instructions table property.
func (t *TablePropInstructions) String() string {
	const text = "TablePropInstructions{Miss: %v, Instructions: %v}"
	return fmt.Sprintf(text, t.Miss, t.Instructions)
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropInstructions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeInstructions,
		TablePropTypeInstructionsMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the table
// instruction property into the wire format.
func (t *TablePropInstructions) WriteTo(w io.Writer) (int64, error) {
	nit := len(t.Instructions)
	proplen := tablePropLen + int(instructionHeaderLen)*nit
	header := tableProp{t.Type(), uint16(proplen)}

	its := make([]instruction, nit)
	for ii, it := range t.Instructions {
		its[ii] = instruction{it, instructionHeaderLen}
	}

	nn, err := encoding.WriteTo(w, header, its)
	if err != nil {
		return nn, err
	}

	n, err := encoding.WriteTo(w, makePad(int(header.Len)))
	return n + nn, err
}

// ReadFrom implements io.ReaderFrom interface. It serializes the table
// instruction property from the wire format.
func (t *TablePropInstructions) ReadFrom(r io.Reader) (int64, error) {
	header, limrd, n, err := yieldTableProp(r, &t.Miss)
	if err != nil {
		return n, err
	}

	// Truncate the list of instructions, but leave the underlying
	// allocated memory, thereby decoded instructions will be saved
	// without memory allocation overhead.
	t.Instructions = t.Instructions[:0]

	maker := encoding.ReaderMakerOf(instruction{})
	join := func(r io.ReaderFrom) {
		t.Instructions = append(
			t.Instructions, r.(*instruction).Type)
	}

	// Read slice of instructions and append them to the slice.
	nn, err := encoding.ReadFunc(limrd, maker, join)
	if n += nn; err != nil {
		return n, err
	}

	nn, err = encoding.ReadFrom(r, makePad(int(header.Len)))
	return n + nn, err
}

// TablePropNextTables defines the table next table property.
type TablePropNextTables struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// NextTables is the list of tables that can be directly reached from
	// the present table using InstructionGotoTable instruction.
	NextTables []Table
}

// String returns a string representation of next tables table property.
func (t *TablePropNextTables) String() string {
	const text = "TablePropNextTables{Miss: %v, NextTables: %v}"
	return fmt.Sprintf(text, t.Miss, t.NextTables)
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropNextTables) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeNextTables,
		TablePropTypeNextTablesMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the next
// tables property into the wire format.
func (t *TablePropNextTables) WriteTo(w io.Writer) (int64, error) {
	headerlen := tablePropLen + len(t.NextTables)
	header := tableProp{t.Type(), uint16(headerlen)}

	padding := makePad(headerlen)
	return encoding.WriteTo(w, header, t.NextTables, padding)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// next tables property from the wire format.
func (t *TablePropNextTables) ReadFrom(r io.Reader) (int64, error) {
	header, _, n, err := yieldTableProp(r, &t.Miss)
	if err != nil {
		return n, err
	}

	t.NextTables = make([]Table, header.Len-tablePropLen)
	padding := makePad(int(header.Len))

	nn, err := encoding.ReadFrom(r, &t.NextTables, padding)
	return n + nn, err
}

// TablePropWriteActions defines the write actions property of the
// table.
type TablePropWriteActions struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Actions is a list of actions for the feature.
	Actions []ActionType
}

// String returns a string representation of write actions table property.
func (t *TablePropWriteActions) String() string {
	const text = "TablePropWriteActions{Miss: %v, Actions: [%v]}"
	return fmt.Sprintf(text, t.Miss, t.Actions)
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropWriteActions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeWriteActions,
		TablePropTypeWriteActionsMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the write
// actions property into the wire format.
func (t *TablePropWriteActions) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropActions(w, t, t.Actions)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// write actions property from the wire format.
func (t *TablePropWriteActions) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropActions(r, &t.Actions, &t.Miss)
}

// TablePropApplyActions defines the apply actions property of the
// table.
type TablePropApplyActions struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Actions is a list of actions for the feature.
	Actions []ActionType
}

// String returns a string representation of apply actions table property.
func (t *TablePropApplyActions) String() string {
	const text = "TablePropApplyActions{Miss: %v, Actions: [%v]}"
	return fmt.Sprintf(text, t.Miss, t.Actions)
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropApplyActions) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeApplyActions,
		TablePropTypeApplyActionsMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the
// apply actions property into the wire format.
func (t *TablePropApplyActions) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropActions(w, t, t.Actions)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// apply actions property from the wire format.
func (t *TablePropApplyActions) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropActions(r, &t.Actions, &t.Miss)
}

// TablePropMatch  defines the match property of the table.
type TablePropMatch struct {
	// Fields is an array of extensible matches.
	Fields []XM
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropMatch) Type() TablePropType {
	return TablePropTypeMatch
}

// WriteTo implements io.WriterTo interface. It serializes the
// match property into the wire format.
func (t *TablePropMatch) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// match property from the wire format.
func (t *TablePropMatch) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropXM(r, &t.Fields, nil)
}

// TablePropWildcards defines the wildcard property of the table.
type TablePropWildcards struct {
	// Fields is an array of extensible matches.
	Fields []XM
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropWildcards) Type() TablePropType {
	return TablePropTypeWildcards
}

// WriteTo implements io.WriterTo interface. It serializes the wildcard
// property into the wire format.
func (t *TablePropWildcards) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// wildcard property from the wire format.
func (t *TablePropWildcards) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropXM(r, &t.Fields, nil)
}

// TablePropWriteSetField defines the write set-field property of the
// table.
type TablePropWriteSetField struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Fields is an array of extensible matches.
	Fields []XM
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropWriteSetField) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeWriteSetField,
		TablePropTypeWriteSetFieldMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the write
// set-field into the wire format.
func (t *TablePropWriteSetField) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the write
// set-field from the wire format.
func (t *TablePropWriteSetField) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropXM(r, &t.Fields, &t.Miss)
}

// TablePropApplySetField defines the apply set-field property of the
// table.
type TablePropApplySetField struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Fields is an array of extensible matches.
	Fields []XM
}

// String returns a string representation of apply set field table property.
func (t *TablePropApplySetField) String() string {
	const text = "TablePropApplySetField{Miss: %v, Fields: %v}"
	return fmt.Sprintf(text, t.Miss, t.Fields)
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropApplySetField) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeApplySetField,
		TablePropTypeApplySetFieldMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the apply
// set-field property into the wire format.
func (t *TablePropApplySetField) WriteTo(w io.Writer) (int64, error) {
	return writeTablePropXM(w, t, t.Fields)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// apply set-field property from the wire format.
func (t *TablePropApplySetField) ReadFrom(r io.Reader) (int64, error) {
	return readTablePropXM(r, &t.Fields, &t.Miss)
}

// TablePropExperimenter defines the experimenter property of the
// table.
type TablePropExperimenter struct {
	// Miss is set to true when it is a property for table-miss.
	Miss bool

	// Experimenter identifier.
	Experimenter uint32

	// Experimenter defined.
	ExpType uint32

	// Experimenter data.
	Data []byte
}

// Type implements TableProp interface. It returns the type of the
// table property.
func (t *TablePropExperimenter) Type() TablePropType {
	return tablePropType(t.Miss,
		TablePropTypeExperimenter,
		TablePropTypeExperimenterMiss)
}

// WriteTo implements io.WriterTo interface. It serializes the
// experimenter property into the wire format.
func (t *TablePropExperimenter) WriteTo(w io.Writer) (int64, error) {
	header := tableProp{t.Type(), uint16(tablePropLen + len(t.Data) + 8)}
	padding := makePad(int(header.Len))

	return encoding.WriteTo(w, header, t.Experimenter,
		t.ExpType, t.Data, padding)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// experimenter property from the wire format.
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
