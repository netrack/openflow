package ofp

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// OpenFlow 1.1 match type MT_STANDARD is deprecated.
	MT_STANDARD MatchType = iota

	// OpenFlow Extensible Match type.
	MT_OXM
)

// MatchType indicates the match structure (set of fields that compose the
// match) in use.
//
// The match type is placed in the type field at the beginning of all match
// structures.
type MatchType uint16

const (
	// Switch input port.
	XMT_OFB_IN_PORT OXMField = iota

	// Switch physical input port.
	XMT_OFB_IN_PHY_PORT

	// Metadata passed between tables.
	XMT_OFB_METADATA

	// Ethernet destination address.
	XMT_OFB_ETH_DST

	// Ethernet source address.
	XMT_OFB_ETH_SRC

	// Ethernet frame type.
	XMT_OFB_ETH_TYPE

	// VLAN identificator.
	XMT_OFB_VLAN_VID

	// VLAN priority.
	XMT_OFB_VLAN_PCP

	// IP DSCP (6 bits in ToS field).
	XMT_OFB_IP_DSCP

	// IP ECN (2 bits in ToS field).
	XMT_OFB_IP_ECN

	// IP protocol.
	XMT_OFB_IP_PROTO

	// IPv4 source address.
	XMT_OFB_IPV4_SRC

	// IPv4 destination address.
	XMT_OFB_IPV4_DST

	// TCP source port.
	XMT_OFB_TCP_SRC

	// TCP destination port.
	XMT_OFB_TCP_DST

	// UDP source port.
	XMT_OFB_UDP_SRC

	// UDP destination port.
	XMT_OFB_UDP_DST

	// SCTP source port.
	XMT_OFB_SCTP_SRC

	// SCTP destination port.
	XMT_OFB_SCTP_DST

	// ICMPv4 type.
	XMT_OFB_ICMPV4_TYPE

	// ICMPv4 code.
	XMT_OFB_ICMPV4_CODE

	// ARP opcode.
	XMT_OFB_ARP_OP

	// ARP source IPv4 address.
	XMT_OFB_ARP_SPA

	// ARP target IPv4 address.
	XMT_OFB_ARP_TPA

	// ARP source hardware address.
	XMT_OFB_ARP_SHA

	// ARP target hardware address.
	XMT_OFB_ARP_THA

	// IPv6 source address.
	XMT_OFB_IPV6_SRC

	// IPv6 destination address.
	XMT_OFB_IPV6_DST

	// IPv6 Flow Label.
	XMT_OFB_IPV6_FLABEL

	// ICMPv6 type.
	XMT_OFB_ICMPV6_TYPE

	// ICMPv6 code.
	XMT_OFB_ICMPV6_CODE

	// Target address for ND.
	XMT_OFB_IPV6_ND_TARGET

	// Source link-layer for ND.
	XMT_OFB_IPV6_ND_SLL

	// Target link-layer for ND.
	XMT_OFB_IPV6_ND_TLL

	// MPLS label.
	XMT_OFB_MPLS_LABEL

	// MPLS TC.
	XMT_OFB_MPLS_TC

	// MPLS BoS bit.
	XMT_OFP_MPLS_BOS

	// PBB I-SID.
	XMT_OFB_PBB_ISID

	// Logical Port Metadata.
	XMT_OFB_TUNNEL_ID

	// IPv6 Extension Header pseudo-field.
	XMT_OFB_IPV6_EXTHDR
)

type OXMField uint8

const (
	// Backward compatibility with NXM.
	XMC_NXM_0 OXMClass = iota

	// Backward compatibility with NXM.
	XMC_NXM_1 OXMClass = iota

	// The class XMC_OPENFLOW_BASIC contains the basic set of OpenFlow
	// match fields.
	XMC_OPENFLOW_BASIC OXMClass = 0x8000

	// The optional class XMC_EXPERIMENTER is used for experimenter
	// matches.
	XMC_EXPERIMENTER OXMClass = 0xffff
)

// OXMClass represents an OXM Class ID. The high order bit differentiate
// reserved classes from member classes.
//
// Classes 0x0000 to 0x7FFF are member classes, allocated by ONF.
//
// Classes 0x8000 to 0xFFFE are reserved classes, reserved for
// standardisation.
type OXMClass uint16

const (
	// Bit that indicate that a VLAN id is set.
	VID_NONE VlanID = iota << 12

	// No VLAN id was set.
	VID_PRESENT VlanID = iota << 12
)

// VlanID represents bit definitions for VLAN id values. It allows matching
// of packets with any tag, independent of the tag's value, and to supports
// matching packets without a VLAN tag.
//
// Values of this type could be used as OXM matching values.
type VlanID uint16

const (
	// "No next header" encountered.
	IEH_NONEXT IPv6ExtHdrFlags = 1 << iota

	// Encrypted Sec Payload header present.
	IEH_ESP IPv6ExtHdrFlags = 1 << iota

	// Authentication header present.
	IEH_AUTH IPv6ExtHdrFlags = 1 << iota

	// One or two destination headers present.
	IEH_DEST IPv6ExtHdrFlags = 1 << iota

	// Fragment header present.
	IEH_FRAG IPv6ExtHdrFlags = 1 << iota

	// Router header present.
	IEH_ROUTER IPv6ExtHdrFlags = 1 << iota

	// Hop-by-hop header present.
	IEH_HOP IPv6ExtHdrFlags = 1 << iota

	// Unexpected repeats encountered.
	IEH_UNREP IPv6ExtHdrFlags = 1 << iota

	// Unexpected sequencing encountered.
	IEH_UNSEQ IPv6ExtHdrFlags = 1 << iota
)

// IPv6ExtHdrFlags represents bit definitions for IPv6 Extension Header
// pseudo field. It indicates the presence of various IPv6 extension
// headers in the packet header.
//
// Values of this type could be used as OXM matching values.
type IPv6ExtHdrFlags uint16

// Match represents fields to match against flows.
type Match struct {
	// Type indicates the match structure (set of fields that compose
	// the match) in use. The match type is placed in the type field at
	// the beginning of all match structures.
	Type MatchType

	// OXMFields lists a set of packet match criteria.
	OXMFields []OXM
}

// Field returns a first entry of OXM by the given type.
func (m *Match) Field(field OXMField) *OXM {
	for _, oxm := range m.OXMFields {
		if oxm.Field == field {
			return &oxm
		}
	}

	return nil
}

// ReadFrom implements ReaderFrom interface.
func (m *Match) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint16
	n, err = binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&m.Type, &length,
	})

	if err != nil {
		return
	}

	var nn int64

	buf := make([]byte, length)
	nn, err = binary.Read(r, binary.BigEndian, &buf)
	n += nn

	if err != nil {
		return
	}

	rbuf := bytes.NewBuffer(buf)
	for rbuf.Len() > 7 {
		var oxm OXM

		_, err = oxm.ReadFrom(rbuf)
		if err != nil {
			return
		}

		m.OXMFields = append(m.OXMFields, oxm)
	}

	return
}

func (m *Match) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, oxm := range m.OXMFields {
		_, err = oxm.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	// Length of Match (excluding padding)
	length := buf.Len() + 4

	if length%8 != 0 {
		_, err = buf.Write(make(pad, 8-length%8))
		if err != nil {
			return
		}
	}

	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		m.Type, uint16(length), buf.Bytes(),
	})
}

// The flow match fields are described using the OpenFlow Extensible
// Match (OXM) format, which is a compact type-length-value (TLV) format.
type OXM struct {
	// Match class that contains related match type
	Class OXMClass
	// Class-specific value, identifying one of the
	// match types within the match class.
	Field OXMField
	Value OXMValue
	Mask  OXMValue
}

func (oxm *OXM) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint8

	n, err = binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&oxm.Class, &oxm.Field, &length,
	})

	if err != nil {
		return
	}

	hasmask := (oxm.Field & 1) == 1
	oxm.Field >>= 1

	var m int64

	oxm.Value = make(OXMValue, length)
	m, err = binary.Read(r, binary.BigEndian, &oxm.Value)
	n += m

	if hasmask {
		length /= 2
		oxm.Mask = make(OXMValue, length)

		m, err = binary.Read(r, binary.BigEndian, &oxm.Mask)
		n += m

		if err != nil {
			return
		}
	}

	return
}

func (oxm *OXM) WriteTo(w io.Writer) (int64, error) {
	var hasmask OXMField
	if len(oxm.Mask) > 0 {
		hasmask = 1
	}

	field := (oxm.Field << 1) | hasmask

	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		oxm.Class, field, uint8(len(oxm.Mask) + len(oxm.Value)), oxm.Value, oxm.Mask,
	})
}

type OXMValue []byte

func (val OXMValue) PortNo() (n PortNo) {
	binary.Read(bytes.NewBuffer(val), binary.BigEndian, &n)
	return
}

func (val OXMValue) UInt32() (v uint32) {
	binary.Read(bytes.NewBuffer(val), binary.BigEndian, &v)
	return
}

func (val OXMValue) UInt16() (v uint16) {
	binary.Read(bytes.NewBuffer(val), binary.BigEndian, &v)
	return
}

func (val OXMValue) UInt8() (v uint8) {
	binary.Read(bytes.NewBuffer(val), binary.BigEndian, &v)
	return
}

type OXMExperimenterHeader struct {
	OXM          OXM
	Experimenter uint32
}
