package ofp13

const (
	MAX_TABLE_NAME_LEN = 32

	TT_MAX Table = 0xfe
	TT_ALL Table = 0xff
)

type Table uint8

const (
	TC_DEPRECATED_MASK TableConfig = 3
)

type TableConfig uint32

type TableMod struct {
	TableID Table
	_       pad3
	Config  TableConfig
}

type TableStats struct {
	TableID      Table
	_            pad3
	ActiveCount  uint32
	LookupCount  uint64
	MatchedCount uint64
}

type TableFeatures struct {
	Length  uint16
	TableID Table
	_       pad5
	name    [MAX_TABLE_NAME_LEN]byte

	MetadataMatch uint64
	MetadataWrite uint64
	Config        TableConfig

	MaxEntries uint32
	Properties []TableFeaturePropHeader
}

const (
	TFPT_INSTRUCTIONS TableFeaturePropType = iota
	TFPT_INSTRUCTIONS_MISS
	TFPT_NEXT_TABLES
	TFPT_NEXT_TABLES_MISS
	TFPT_WRITE_ACTIONS
	TFPT_WRITE_ACTIONS_MISS
	TFPT_APPLY_ACTIONS
	TFPT_APPLY_ACTIONS_MISS
	TFPT_MATCH
	TFPT_WILDCARDS
	TFPT_WRITE_SETFIELD
	TFPT_WRITE_SETFIELD_MISS
	TFPT_APPLY_SETFIELD
	TFPT_APPLY_SETFIELD_MISS
	TFPT_EXPERIMENTER      TableFeaturePropType = 0xfffe
	TFPT_EXPERIMENTER_MISS TableFeaturePropType = 0xffff
)

type TableFeaturePropType uint16

type TableFeaturePropHeader struct {
	Type   TableFeaturePropType
	Length uint16
}

type TableFeaturePropInstructions struct {
	Type   TableFeaturePropType
	Length uint16
	//TODO: InstructionID  []Instruction
}

type TableFeaturePropNextTables struct {
	Type        TableFeaturePropType
	Length      uint16
	NextTableID []Table
}

type TableFeaturePropActions struct {
	Type     TableFeaturePropType
	Length   uint16
	ActionID []interface{}
}

type TableFeaturePropOXM struct {
	Type   TableFeaturePropType
	Length uint16
	OXMID  []OXM
}

type TableFeaturePropExperimenter struct {
	Type             TableFeaturePropType
	Length           uint16
	Experimenter     uint32
	ExpType          uint32
	ExperimenterData []uint32
}
