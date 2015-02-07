package ofp13

const (
	AT_OUTPUT       ActionType = iota
	AT_COPY_TTL_OUT ActionType = 10 + iota
	AT_COPY_TTL_IN  ActionType = 10 + iota
	AT_SET_MPLS_TTL ActionType = 12 + iota
	AT_DEC_MPLS_TTL ActionType = 12 + iota
	AT_PUSH_VLAN    ActionType = 12 + iota
	AT_POP_VLAN     ActionType = 12 + iota
	AT_PUSH_MPLS    ActionType = 12 + iota
	AT_POP_MPLS     ActionType = 12 + iota
	AT_SET_QUEUE    ActionType = 12 + iota
	AT_GROUP        ActionType = 12 + iota
	AT_SET_NW_TTL   ActionType = 12 + iota
	AT_DEC_NW_TTL   ActionType = 12 + iota
	AT_SET_FIELD    ActionType = 12 + iota
	AT_PUSH_PBB     ActionType = 12 + iota
	AT_POP_PBB      ActionType = 12 + iota
	AT_EXPERIMENTER ActionType = 0xffff
)

type ActionType uint16

const (
	// Maximum MaxLength value which can be used
	// to request a specific byte Length.
	CML_MAX = 0xffe5
	// Indicates that no buffering should be
	// applied and the whole packet is to be
	// sent to the controller.
	CML_NO_BUFFER = 0xffff
)

type ActionHeader struct {
	Type   ActionType
	Length uint16
	_      pad4
}

type ActionOutput struct {
	ActionHeader
	Port      uint32
	MaxLength uint16
	_         pad6
}

type ActionGroup struct {
	ActionHeader
	GroupId uint32
}

type ActionSetQueue struct {
	ActionHeader
	QueueId uint32
}

type ActionMPLSTTL struct {
	ActionHeader
	MPLSTTL uint8
	_       pad3
}

type ActionSetNWTTL struct {
	ActionHeader
	NWTTL uint8
	_     pad3
}

type ActionPush struct {
	ActionHeader
	Ethertype uint16
	_         pad2
}

type ActionPopMPLS struct {
	ActionHeader
	Ethertype uint16
	_         pad2
}

type ActionSetField struct {
	ActionHeader
	Fields []OXMHeader
}

type ActionExperimenterHeader struct {
	ActionHeader
	Experimenter uint32
}
