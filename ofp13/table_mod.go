package ofp13

const (
	TT_MAX Table = 0xfe
	TT_ALL Table = 0xff
)

type Table uint8

type TableConfigFlags uint32

type TableMod struct {
	Header  Header
	TableId Table
	Config  TableConfigFlags
}
