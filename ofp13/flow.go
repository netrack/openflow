package ofp13

import (
	"bytes"
	"encoding/binary"
	"io"
)

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

const (
	VID_NONE    VlanID = iota << 12
	VID_PRESENT VlanID = iota << 12
)

type VlanID uint16

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

// Fields to match against flows
type Match struct {
	// Type indicates the match structure (set of fields that compose the match) in use.
	// The match type is placed in the type field at the beginning of all match structures.
	Type      MatchType
	OXMFields []OXM
}

func (m *Match) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &m.Type)
	if err != nil {
		return err
	}

	var length uint16
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	buf := make([]byte, length)
	err = binary.Read(r, binary.BigEndian, &buf)
	if err != nil {
		return err
	}

	// TODO: invalid condition
	rbuf := bytes.NewBuffer(buf)
	for rbuf.Len() > 4 {
		var oxm OXM

		err = oxm.Read(rbuf)
		if err != nil {
			return err
		}

		m.OXMFields = append(m.OXMFields, oxm)
	}

	return nil
}

type OXM struct {
	Class OXMClass
	Field OXMField
	Mask  []byte
	Value []byte
}

func (oxm *OXM) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &oxm.Class)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &oxm.Field)
	if err != nil {
		return err
	}

	hasmask := oxm.Field&1 == 1
	oxm.Field >>= 1

	var length uint8
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	if hasmask {
		length /= 2
		oxm.Mask = make([]byte, length)

		err = binary.Read(r, binary.BigEndian, &oxm.Mask)
		if err != nil {
			return err
		}
	}

	oxm.Value = make([]byte, length)
	return binary.Read(r, binary.BigEndian, &oxm.Value)
}

type OXMExperimenterHeader struct {
	OXM          OXM
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
	TableID uint8
	_       pad3
}

type InstructionWriteMetadata struct {
	Type         InstructionType
	Length       uint16
	_            pad4
	Metadata     uint64
	MetadataMask uint64
}

type InstructionActions struct {
	Type    InstructionType
	Length  uint16
	_       pad4
	Actions []ActionHeader
}

type InstructionMeter struct {
	Type    InstructionType
	Length  uint16
	MeterID uint8
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
	Cookie      uint64
	CookieMask  uint64
	TableID     Table
	Command     FlowModCommand
	IdleTimeout uint16
	HardTimeout uint16
	Priority    uint16
	BufferID    uint16
	OutPort     PortNo
	OutGroup    Group

	Flags FlowModFlags
	_     pad2
	Match Match
}

func (f *FlowMod) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, f)
}

const (
	RR_IDLE_TIMEOUT FlowRemovedReason = iota
	RR_HARD_TIMEOUT
	RR_DELETE
	RR_GROUP_DELETE
)

type FlowRemovedReason uint8

type FlowRemoved struct {
	Cookie       uint64
	Priority     uint16
	Reason       FlowRemovedReason
	TableID      Table
	DurationSec  uint32
	DurationNSec uint32
	IdleTimeout  uint16
	HardTimeout  uint16
	PacketCount  uint64
	ByteCount    uint64
	Match        Match
}

type FlowStatsRequest struct {
	TableID    Table
	_          pad3
	OutPort    PortNo
	OutGroup   Group
	_          pad4
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

type FlowStats struct {
	Length       uint16
	TableID      Table
	_            pad1
	DurationSec  uint32
	DurationNSec uint32

	Priority    uint16
	IdleTimeout uint16
	HardTimeout uint16
	Flags       FlowModFlags
	_           pad4
	Cookie      uint64
	PacketCount uint64
	ByteCount   uint64
	Match       Match
}
