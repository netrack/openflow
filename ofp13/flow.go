package ofp13

const (
	MT_STANDARD MatchType = iota
	MT_OXM
)

type MatchType uint16

const (
	XMT_OFB_IN_PORT OXMField = iota
	XMT_OFB_IN_PHY_PORT
	XMT_OFB_METADATA
	XMT_OFB_ETH_DST
	XMT_OFB_ETH_SRC
	XMT_OFB_ETH_TYPE
	XMT_OFB_VLAN_VID
	XMT_OFB_VLAN_PCP
	XMT_OFB_IP_DSCP
	XMT_OFB_IP_ECN
	XMT_OFB_IP_PROTO
	XMT_OFB_IPV4_SRC
	XMT_OFB_IPV4_DST
	XMT_OFB_TCP_SRC
	XMT_OFB_TCP_DST
	XMT_OFB_UDP_SRC
	XMT_OFB_UDP_DST
	XMT_OFB_SCTP_SRC
	XMT_OFB_SCTP_DST
	XMT_OFB_ICMPV4_TYPE
	XMT_OFB_ICMPV4_CODE
	XMT_OFB_ARP_OP
	XMT_OFB_ARP_SPA
	XMT_OFB_ARP_TPA
	XMT_OFB_ARP_SHA
	XMT_OFB_ARP_THA
	XMT_OFB_IPV6_SRC
	XMT_OFB_IPV6_DST
	XMT_OFB_IPV6_FLABEL
	XMT_OFB_ICMPV6_TYPE
	XMT_OFB_ICMPV6_CODE
	XMT_OFB_IPV6_ND_TARGET
	XMT_OFB_IPV6_ND_SLL
	XMT_OFB_IPV6_ND_TLL
	XMT_OFB_MPLS_LABEL
	XMT_OFB_MPLS_TC
	XMT_OFP_MPLS_BOS
	XMT_OFB_PBB_ISID
	XMT_OFB_TUNNEL_ID
	XMT_OFB_IPV6_EXTHDR
)

type OXMField uint8

const (
	XMC_NXM_0          OXMClass = iota
	XMC_NXM_1          OXMClass = iota
	XMC_OPENFLOW_BASIC OXMClass = 0x8000
	XMC_EXPERIMENTER   OXMClass = 0xffff
)

type OXMClass uint16

var (
	OXM_OF_IN_PORT        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IN_PORT, false, 4}
	OXM_OF_IN_PHY_PORT    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IN_PHY_PORT, false, 4}
	OXM_OF_METADATA       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_METADATA, false, 8}
	OXM_OF_ETH_DST        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ETH_DST, false, 6}
	OXM_OF_ETH_SRC        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ETH_SRC, false, 6}
	OXM_OF_ETH_TYPE       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ETH_TYPE, false, 6}
	OXM_OF_VLAN_VID       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_VLAN_VID, false, 3}
	OXM_OF_VLAN_PCP       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_VLAN_PCP, false, 1}
	OXM_OF_IP_DSCP        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IP_DSCP, false, 1}
	OXM_OF_IP_ECN         = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IP_ECN, false, 1}
	OXM_OF_IP_PROTO       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IP_PROTO, false, 1}
	OXM_OF_IPV4_SRC       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV4_SRC, false, 4}
	OXM_OF_IPV4_DST       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV4_DST, false, 4}
	OXM_OF_TCP_SRC        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_TCP_SRC, false, 2}
	OXM_OF_TCP_DST        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_TCP_DST, false, 2}
	OXM_OF_UDP_SRC        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_UDP_SRC, false, 2}
	OXM_OF_UDP_DST        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_UDP_DST, false, 2}
	OXM_OF_SCTP_SRC       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_SCTP_SRC, false, 2}
	OXM_OF_SCTP_DST       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_SCTP_DST, false, 2}
	OXM_OF_ICMPV4_TYPE    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ICMPV4_TYPE, false, 1}
	OXM_OF_ICMPV4_CODE    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ICMPV4_CODE, false, 1}
	OXM_OF_ARP_OP         = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ARP_OP, false, 2}
	OXM_OF_ARP_SPA        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ARP_SPA, false, 4}
	OXM_OF_ARP_TPA        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ARP_TPA, false, 4}
	OXM_OF_ARP_SHA        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ARP_SHA, false, 6}
	OXM_OF_ARP_THA        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ARP_THA, false, 6}
	OXM_OF_IPV6_SRC       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_SRC, false, 16}
	OXM_OF_IPV6_DST       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_DST, false, 16}
	OXM_OF_IPV6_FLABEL    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_FLABEL, false, 4}
	OXM_OF_ICMPV6_TYPE    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ICMPV6_TYPE, false, 1}
	OXM_OF_ICMPV6_CODE    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_ICMPV6_CODE, false, 1}
	OXM_OF_IPV6_ND_TARGET = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_ND_TARGET, false, 16}
	OXM_OF_IPV6_ND_SLL    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_ND_SLL, false, 6}
	OXM_OF_IPV6_ND_TLL    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_ND_TLL, false, 6}
	OXM_OF_MPLS_LABEL     = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_MPLS_LABEL, false, 3}
	OXM_OF_MPLS_TC        = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_MPLS_TC, false, 1}
	OXM_OF_MPLS_BOS       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFP_MPLS_BOS, false, 1}
	OXM_OF_PBB_ISID       = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_PBB_ISID, false, 3}
	OXM_OF_TUNNEL_ID      = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_TUNNEL_ID, false, 8}
	OXM_OF_IPV6_EXTHDR    = OXMHeader{XMC_OPENFLOW_BASIC, XMT_OFB_IPV6_EXTHDR, false, 2}
)

const (
	VID_NONE    VlanId = iota << 12
	VID_PRESENT VlanId = iota << 12
)

type VlanId uint16

const (
	IEH_NONEXT IPv6ExtHdrFlags = 1 << iota
	IEH_ESP    IPv6ExtHdrFlags = 1 << iota
	IEH_AUTH   IPv6ExtHdrFlags = 1 << iota
	IEH_DEST   IPv6ExtHdrFlags = 1 << iota
	IEH_FRAG   IPv6ExtHdrFlags = 1 << iota
	IEH_ROUTER IPv6ExtHdrFlags = 1 << iota
	IEH_HOP    IPv6ExtHdrFlags = 1 << iota
	IEH_UNREP  IPv6ExtHdrFlags = 1 << iota
	IEH_UNSEQ  IPv6ExtHdrFlags = 1 << iota
)

type IPv6ExtHdrFlags uint16

type Match struct {
	Type      MatchType
	Length    uint16
	OXMFields OXMHeader
}

type OXMHeader struct {
	Class   OXMClass
	Field   OXMField
	HasMask bool
	Length  uint8
}

type OXMExperimenterHeader struct {
	OXMHeader    OXMHeader
	Experimenter uint32
}

const (
	IT_GOTO_TABLE     InstructionType = 1 + iota
	IT_WRITE_METADATA InstructionType = 1 + iota
	IT_WRITE_ACTIONS  InstructionType = 1 + iota
	IT_APPLY_ACTIONS  InstructionType = 1 + iota
	IT_CLEAR_ACTIONS  InstructionType = 1 + iota
	IT_METER          InstructionType = 1 + iota
	IT_EXPERIMENTER   InstructionType = 0xffff
)

type InstructionType uint16

type InstrutionGotoTable struct {
	Type    InstructionType
	Length  uint16
	TableId uint8
}

type InstructionWriteMetadata struct {
	Type         InstructionType
	Length       uint16
	Metadata     uint64
	MetadataMask uint64
}

type InstructionActions struct {
	Type    InstructionType
	Length  uint16
	Actions []ActionHeader
}

type InstructionMeter struct {
	Type    InstructionType
	Length  uint16
	MeterId uint8
}

const (
	FC_ADD FlowModCommand = iota
	FC_MODIFY
	FC_MODIFY_STRICT
	FC_DELETE
	FC_DELETE_STRICT
)

type FlowModCommand uint8

const (
	// Send flow removed message when flow expires or is deleted
	FF_SEND_FLOW_REM FlowModFlags = 1 << iota
	// Check for overlapping entries first
	FF_CHECK_OVERLAP FlowModFlags = 1 << iota
	// Reset flow packet and byte counts
	FF_RESET_COUNTS FlowModFlags = 1 << iota
	// Don't keep track of packet count
	FF_NO_PKT_COUNTS FlowModFlags = 1 << iota
	// Don't keep track of byte count
	FF_NO_BYT_COUNTS FlowModFlags = 1 << iota
)

type FlowModFlags uint16

type FlowMod struct {
	Header      Header
	Cookie      uint64
	CookieMask  uint64
	TableId     Table
	Command     FlowModCommand
	IdleTimeout uint16
	HardTimeout uint16
	Priority    uint16
	BufferId    uint16
	OutPort     PortNo
	OutGroup    Group

	Flags FlowModFlags
	Match Match
}

type FlowStatsRequest struct {
	TableId    Table
	OutPort    PortNo
	OutGroup   Group
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

type FlowStats struct {
	Length       uint16
	TableId      Table
	DurationSec  uint32
	DurationNSec uint32

	Priority    uint16
	IdleTimeout uint16
	HardTimeout uint16
	Flags       FlowModFlags
	Cookie      uint64
	PacketCount uint64
	ByteCount   uint64
	Match       Match
}
