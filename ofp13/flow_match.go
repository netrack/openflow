package ofp13

const (
	MT_STANDARD MatchType = iota
	MT_OXM      MatchType = iota
)

const (
	XMC_NXM_0          OXMClass = iota
	XMC_NXM_1          OXMCLass = iota
	XMC_OPENFLOW_BASIC OXMClass = 0x8000
	XMC_EXPERIMENTER   OXMClass = 0xffff
)

type MatchType uint16

type OXMClass uint16

type OXMField uint8

type Match struct {
	Type   MatchType
	Length uint16
	OXMFields
}
