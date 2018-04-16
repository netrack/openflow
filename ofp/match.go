package ofp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// XMType defines the flow match field types for OpenFlow basic class.
type XMType uint8

func (t XMType) String() string {
	text, ok := xmTypeText[t]
	if !ok {
		return fmt.Sprintf("XMType(%d)", t)
	}
	return text
}

const (
	// XMTypeInPort matches switch input port.
	XMTypeInPort XMType = iota

	// XMTypeInPhyPort matches switch physical input port.
	XMTypeInPhyPort

	// XMTypeMetadata matches metadata passed between tables.
	XMTypeMetadata

	// XMTypeEthDst matches ethernet destination address.
	XMTypeEthDst

	// XMTypeEthSrc matches ethernet source address.
	XMTypeEthSrc

	// XMTypeEthType matches ethernet frame type.
	XMTypeEthType

	// XMTypeVlanID matches VLAN identifier.
	XMTypeVlanID

	// XMTypeVlanPCP matches VLAN priority.
	XMTypeVlanPCP

	// XMTypeIPDSCP matches IP DSCP (6 bits in ToS field).
	XMTypeIPDSCP

	// XMTypeIPECN matches IP ECN (2 bits in ToS field).
	XMTypeIPECN

	// XMTypeIPProto matches IP protocol.
	XMTypeIPProto

	// XMTypeIPv4Src matches IPv4 source address.
	XMTypeIPv4Src

	// XMTypeIPv4Dst matches IPv4 destination address.
	XMTypeIPv4Dst

	// XMTypeTCPSrc matches TCP source port.
	XMTypeTCPSrc

	// XMTypeTCPDst matches TCP destination port.
	XMTypeTCPDst

	// XMTypeUDPSrc matches UDP source port.
	XMTypeUDPSrc

	// XMTypeUDPDst matches UDP destination port.
	XMTypeUDPDst

	// XMTypeSCTPSrc matches SCTP source port.
	XMTypeSCTPSrc

	// XMTypeSCTPDst matches SCTP destination port.
	XMTypeSCTPDst

	// XMTypeICMPv4Type matches ICMPv4 type.
	XMTypeICMPv4Type

	// XMTypeICMPv4Code matches ICMPv4 code.
	XMTypeICMPv4Code

	// XMTypeARPOpcode matches ARP opcode.
	XMTypeARPOpcode

	// XMTypeARPSPA matches ARP source IPv4 address.
	XMTypeARPSPA

	// XMTypeARPTPA matches ARP target IPv4 address.
	XMTypeARPTPA

	// XMTypeARPSHA matches ARP source hardware address.
	XMTypeARPSHA

	// XMTypeARPTHA matches ARP target hardware address.
	XMTypeARPTHA

	// XMTypeIPv6Src matches IPv6 source address.
	XMTypeIPv6Src

	// XMTypeIPv6Dst matches IPv6 destination address.
	XMTypeIPv6Dst

	// XMTypeIPv6FLabel matches IPv6 Flow Label.
	XMTypeIPv6FLabel

	// XMTypeICMPv6Type matches ICMPv6 type.
	XMTypeICMPv6Type

	// XMTypeICMPv6Code matches ICMPv6 code.
	XMTypeICMPv6Code

	// XMTypeIPv6NDTarget matches IPv6 target address for ND.
	XMTypeIPv6NDTarget

	// XMTypeIPv6NDSLL matches IPv6 source link-layer for ND.
	XMTypeIPv6NDSLL

	// XMTypeIPv6NDTLL matches IPv6 target link-layer for ND.
	XMTypeIPv6NDTLL

	// XMTypeMPLSLabel matches MPLS label.
	XMTypeMPLSLabel

	// XMTypeMPLSTC matches MPLS TC.
	XMTypeMPLSTC

	// XMTypeMPLSBOS matches MPLS BoS bit.
	XMTypeMPLSBOS

	// XMTypePBBISID matches PBB I-SID.
	XMTypePBBISID

	// XMTypeTunnelID matches logical Port Metadata.
	XMTypeTunnelID

	// XMTypeIPv6ExtHeader matches IPv6 extension Header pseudo-field.
	XMTypeIPv6ExtHeader
)

var xmTypeText = map[XMType]string{
	XMTypeInPort:        "XMTypeInPort",
	XMTypeInPhyPort:     "XMTypeInPhyPort",
	XMTypeMetadata:      "XMTypeMetadata",
	XMTypeEthDst:        "XMTypeEthDst",
	XMTypeEthSrc:        "XMTypeEthSrc",
	XMTypeEthType:       "XMTypeEthType",
	XMTypeVlanID:        "XMTypeVlanID",
	XMTypeVlanPCP:       "XMTypeVlanPCP",
	XMTypeIPDSCP:        "XMTypeIPDSCP",
	XMTypeIPECN:         "XMTypeIPECN",
	XMTypeIPProto:       "XMTypeIPProto",
	XMTypeIPv4Src:       "XMTypeIPv4Src",
	XMTypeIPv4Dst:       "XMTypeIPv4Dst",
	XMTypeTCPSrc:        "XMTypeTCPSrc",
	XMTypeTCPDst:        "XMTypeTCPDst",
	XMTypeUDPSrc:        "XMTypeUDPSrc",
	XMTypeUDPDst:        "XMTypeUDPDst",
	XMTypeSCTPSrc:       "XMTypeSCTPSrc",
	XMTypeSCTPDst:       "XMTypeSCTPDst",
	XMTypeICMPv4Type:    "XMTypeICMPv4Type",
	XMTypeICMPv4Code:    "XMTypeICMPv4Code",
	XMTypeARPOpcode:     "XMTypeARPOpcode",
	XMTypeARPSPA:        "XMTypeARPSPA",
	XMTypeARPTPA:        "XMTypeARPTPA",
	XMTypeARPSHA:        "XMTypeARPSHA",
	XMTypeARPTHA:        "XMTypeARPTHA",
	XMTypeIPv6Src:       "XMTypeIPv6Src",
	XMTypeIPv6Dst:       "XMTypeIPv6Dst",
	XMTypeIPv6FLabel:    "XMTypeIPv6FLabel",
	XMTypeICMPv6Type:    "XMTypeICMPv6Type",
	XMTypeICMPv6Code:    "XMTypeICMPv6Code",
	XMTypeIPv6NDTarget:  "XMTypeIPv6NDTarget",
	XMTypeIPv6NDSLL:     "XMTypeIPv6NDSLL",
	XMTypeIPv6NDTLL:     "XMTypeIPv6NDTLL",
	XMTypeMPLSLabel:     "XMTypeMPLSLabel",
	XMTypeMPLSTC:        "XMTypeMPLSTC",
	XMTypeMPLSBOS:       "XMTypeMPLSBOS",
	XMTypePBBISID:       "XMTypePBBISID",
	XMTypeTunnelID:      "XMTypeTunnelID",
	XMTypeIPv6ExtHeader: "XMTypeIPv6ExtHeader",
}

// XMClass represents an OXM Class ID. The high order bit differentiate
// reserved classes from member classes.
//
// Classes 0x0000 to 0x7FFF are member classes, allocated by ONF.
//
// Classes 0x8000 to 0xFFFE are reserved classes, reserved for
// standardisation.
type XMClass uint16

func (c XMClass) String() string {
	text, ok := classText[c]
	if !ok {
		return fmt.Sprintf("XMClass(%d)", c)
	}
	return text
}

const (
	// XMClassNicira0 defines a backward compatibility class with NXM.
	XMClassNicira0 XMClass = iota

	// XMClassNicira1 defines a backward compatibility class with NXM.
	XMClassNicira1

	// XMClassOpenflowBasic defines a class with the basic set of
	// OpenFlow match fields.
	XMClassOpenflowBasic XMClass = 0x8000

	// XMClassExperimenter defines a class of experimenter matches.
	XMClassExperimenter XMClass = 0xffff
)

var classText = map[XMClass]string{
	XMClassNicira0:       "XMClassNicira0",
	XMClassNicira1:       "XMClassNicira1",
	XMClassOpenflowBasic: "XMClassOpenflowBasic",
	XMClassExperimenter:  "XMClassExperimeter",
}

// VlanID represents bit definitions for VLAN ID values. It allows matching
// of packets with any tag, independent of the tag's value, and to supports
// matching packets without a VLAN tag.
//
// Values of this type could be used as OXM matching values.
type VlanID uint16

const (
	// VlanNone indicates that no VLAN ID was set.
	VlanNone VlanID = iota << 12

	// VlanPresent indicates that a VLAN ID is set.
	VlanPresent VlanID = iota << 12
)

// IPv6ExtensionHeader represents bit definitions for IPv6 Extension
// Header pseudo field. It indicates the presence of various IPv6
// extension headers in the packet header.
//
// Values of this type could be used as XM matching values.
type IPv6ExtensionHeader uint16

const (
	// IPv6ExtensionHeaderNoNext bit is set when "no next header"
	// encountered.
	IPv6ExtensionHeaderNoNext IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderESP bit is set when Encrypted Sec Payload header
	// present.
	IPv6ExtensionHeaderESP

	// IPv6ExtensionHeaderAuth bit is set when authentication header
	// present.
	IPv6ExtensionHeaderAuth

	// IPv6ExtensionHeaderDest bit is set when one or two destination
	// headers present.
	IPv6ExtensionHeaderDest

	// IPv6ExtensionHeaderFrag bit is set when fragment header present.
	IPv6ExtensionHeaderFrag

	// IPv6ExtensionHeaderRouter bit is set when router header present.
	IPv6ExtensionHeaderRouter

	// IPv6ExtensionHeaderHop bit is set when hop-by-hop header present.
	IPv6ExtensionHeaderHop

	// IPv6ExtensionHeaderUnrep bit is set when unexpected repeats
	// encountered.
	IPv6ExtensionHeaderUnrep

	// IPv6ExtensionHeaderUnseq bit is set when unexpected sequencing
	// encountered.
	IPv6ExtensionHeaderUnseq
)

const (
	// xmlen defines the length of the extention match header, it does
	// not include the value and mask.
	xmlen = 4
)

// The XM defines the flow match fields are described using the OpenFlow
// Extensible Match (OXM) format, which is a compact type-length-value
// (TLV) format.
//
// For example, to match the packets arrived at the first switch port,
// the following extensible match can be created:
//
//	m := ofp.XM{
//		Class: ofp.XMClassOpenflowBasic,
//		Type:  ofp.XMTypeInPort,
//		Value: ofp.XMValue{0, 0, 0, 1},
//	}
type XM struct {
	// Match class that contains related match type
	Class XMClass

	// Class-specific value, identifying one of the
	// match types within the match class.
	Type XMType

	// Type-specific value and mask.
	Value XMValue
	Mask  XMValue
}

// readAllXM uses all available bytes retrieved from the reader to
// unmarshal them to the list of extensible matchers. The caller
// responsible of passing limited reader to prevent from read of
// unnecessary data.
func readAllXM(r io.Reader, xms *[]XM, hasPayload bool) (int64, error) {
	// Read all available bytes from the reader, they will be
	// used to unmarshal them into the list of extensible matchers.
	buf, err := ioutil.ReadAll(r)
	n := int64(len(buf))

	if err != nil {
		return n, err
	}

	rbuf := bytes.NewBuffer(buf)
	for rbuf.Len() >= xmlen {
		var xm XM

		_, err = xm.readFrom(rbuf, hasPayload)
		if err != nil {
			return n, err
		}

		*xms = append(*xms, xm)
	}

	return n, nil
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the OpenFlow extensible match from the given reader.
func (xm *XM) ReadFrom(r io.Reader) (n int64, err error) {
	return xm.readFrom(r, true)
}

// readFrom deserializes the OpenFlow extensible match from the
// given reader.  If hasPayload is false, xm.Value and xm.Mask
// will be filled with the zero value.
func (xm *XM) readFrom(r io.Reader, hasPayload bool) (n int64, err error) {
	var length uint8

	n, err = encoding.ReadFrom(
		r, &xm.Class, &xm.Type, &length)
	if err != nil {
		return
	}

	hasmask := (xm.Type & 1) == 1
	xm.Type >>= 1

	xm.Value, xm.Mask = make(XMValue, length), nil

	if hasPayload {
		var m int64
		m, err = encoding.ReadFrom(r, &xm.Value)
		n += m
		if err != nil {
			return
		}
	}

	if hasmask {
		length /= 2
		xm.Mask = make(XMValue, length)

		copy(xm.Mask, xm.Value[length:])
		xm.Value = xm.Value[:length]
	}

	return
}

// WriteTo implements io.WriterTo interface. It serializes the OpenFlow
// extensible match into given writer.
func (xm *XM) WriteTo(w io.Writer) (int64, error) {
	var hasmask XMType
	if len(xm.Mask) > 0 {
		hasmask = 1
	}

	field := (xm.Type << 1) | hasmask

	return encoding.WriteTo(
		w, xm.Class, field,
		uint8(len(xm.Mask)+len(xm.Value)),
		xm.Value, xm.Mask)
}

// XMValue is a value of the extensible match.
type XMValue []byte

// UInt32 returns the value of extensible match as int32.
func (val XMValue) UInt32() (v uint32) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

// UInt16 returns the value of extensible match as int16.
func (val XMValue) UInt16() (v uint16) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

// UInt8 returns the value of extensible match as uint8.
func (val XMValue) UInt8() (v uint8) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

const (
	// MatchTypeStandard is an OpenFlow 1.1 match type MT_STANDARD is
	// deprecated.
	MatchTypeStandard MatchType = iota

	// MatchTypeXM is an OpenFlow Extensible Match type.
	MatchTypeXM
)

// MatchType indicates the match structure (set of fields that compose the
// match) in use.
//
// The match type is placed in the type field at the beginning of all match
// structures.
type MatchType uint16

// Match represents fields to match against flows.
type Match struct {
	// Type indicates the match structure (set of fields that compose
	// the match) in use. The match type is placed in the type field at
	// the beginning of all match structures.
	Type MatchType

	// Fields lists a set of packet match criteria.
	Fields []XM
}

// Field returns a first entry of OXM by the given type.
func (m *Match) Field(mt XMType) *XM {
	for _, xm := range m.Fields {
		if xm.Type == mt {
			return &xm
		}
	}

	return nil
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// match from the wire format.
func (m *Match) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64
	var length uint16

	// Initialize the structure attributes with default
	// values, so we could read multiple times into the
	// same variable.
	m.Type, m.Fields = 0, nil

	n, err = encoding.ReadFrom(r, &m.Type, &length)
	if err != nil {
		return
	}

	matchlen := int(length)

	// subtract the length of the already-read Type & Length fields
	rdlen := matchlen - 4

	// Limit the reader to the length of the extensible matches.
	limrd := io.LimitReader(r, int64(rdlen))
	nn, err = readAllXM(limrd, &m.Fields, true)
	if n += nn; err != nil {
		return
	}

	// Read the padding after the list of extensible matches.
	nn, err = encoding.ReadFrom(r, makePad(matchlen))
	return n + nn, err
}

// WriteTo implements io.WriterTo interface. It serializes the match
// into the wire format.
func (m *Match) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, xm := range m.Fields {
		_, err = xm.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	// Length of Match (excluding padding)
	length := buf.Len() + 4
	padding := makePad(length)

	return encoding.WriteTo(
		w, m.Type, uint16(length), buf.Bytes(), padding)
}
