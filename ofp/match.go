package ofp

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

const (
	// Switch input port.
	XMTypeInPort XMType = iota // 0

	// Switch physical input port.
	XMTypeInPhyPort

	// Metadata passed between tables.
	XMTypeMetadata

	// Ethernet destination address.
	XMTypeEthDst

	// Ethernet source address.
	XMTypeEthSrc

	// Ethernet frame type.
	XMTypeEthType // 5

	// VLAN identificator.
	XMTypeVlanID

	// VLAN priority.
	XMTypeVlanPCP

	// IP DSCP (6 bits in ToS field).
	XMTypeIPDSCP

	// IP ECN (2 bits in ToS field).
	XMTypeIPECN

	// IP protocol.
	XMTypeIPProto // 10

	// IPv4 source address.
	XMTypeIPv4Src

	// IPv4 destination address.
	XMTypeIPv4Dst

	// TCP source port.
	XMTypeTCPSrc

	// TCP destination port.
	XMTypeTCPDst

	// UDP source port.
	XMTypeUDPSrc // 15

	// UDP destination port.
	XMTypeUDPDst

	// SCTP source port.
	XMTypeSCTPSrc

	// SCTP destination port.
	XMTypeSCTPDst

	// ICMPv4 type.
	XMTypeICMPv4Type

	// ICMPv4 code.
	XMTypeICMPv4Code // 20

	// ARP opcode.
	XMTypeARPOpcode

	// ARP source IPv4 address.
	XMTypeARPSPA

	// ARP target IPv4 address.
	XMTypeARPTPA

	// ARP source hardware address.
	XMTypeARPSHA

	// ARP target hardware address.
	XMTypeARPTHA // 25

	// IPv6 source address.
	XMTypeIPv6Src

	// IPv6 destination address.
	XMTypeIPv6Dst

	// IPv6 Flow Label.
	XMTypeIPv6Flabel

	// ICMPv6 type.
	XMTypeICMPv6Type

	// ICMPv6 code.
	XMTypeICMPv6Code // 30

	// Target address for ND.
	XMTypeIPv6NDTarget

	// Source link-layer for ND.
	XMTypeIPv6NDSLL

	// Target link-layer for ND.
	XMTypeNDTLL

	// MPLS label.
	XMTypeMPLSLabel

	// MPLS TC.
	XMTypeMPLSTC // 35

	// MPLS BoS bit.
	XMTypeMPLSBOS

	// PBB I-SID.
	XMTypePBBISID

	// Logical Port Metadata.
	XMTypeTunnelID

	// IPv6 Extension Header pseudo-field.
	XMTypeIPv6EXTHDR
)

type XMType uint8

const (
	// Backward compatibility with NXM.
	XMClassNicira0 XMClass = iota

	// Backward compatibility with NXM.
	XMClassNicira1

	// The class XMC_OPENFLOW_BASIC contains the basic set of OpenFlow
	// match fields.
	XMClassOpenflowBasic XMClass = 0x8000

	// The optional class XMC_EXPERIMENTER is used for experimenter
	// matches.
	XMClassExperimenter XMClass = 0xffff
)

// XMClass represents an OXM Class ID. The high order bit differentiate
// reserved classes from member classes.
//
// Classes 0x0000 to 0x7FFF are member classes, allocated by ONF.
//
// Classes 0x8000 to 0xFFFE are reserved classes, reserved for
// standardisation.
type XMClass uint16

const (
	// Bit that indicate that a VLAN id is set.
	VlanNone VlanID = iota << 12

	// No VLAN id was set.
	VlanPresent VlanID = iota << 12
)

// VlanID represents bit definitions for VLAN id values. It allows matching
// of packets with any tag, independent of the tag's value, and to supports
// matching packets without a VLAN tag.
//
// Values of this type could be used as OXM matching values.
type VlanID uint16

const (
	// IPv6ExtensionHeaderNoNext bit is set when "no next header"
	// encountered.
	IPv6ExtensionHeaderNoNext IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderESP bit is set when Encrypted Sec Payload header
	// present.
	IPv6ExtensionHeaderESP IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderAuth bit is set when authentication header
	// present.
	IPv6ExtensionHeaderAuth IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderDest bit is set when one or two destination
	// headers present.
	IPv6ExtensionHeaderDest IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderFrag bit is set when fragment header present.
	IPv6ExtensionHeaderFrag IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderRouter bit is set when router header present.
	IPv6ExtensionHeaderRouter IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderHop bit is set when hop-by-hop header present.
	IPv6ExtensionHeaderHop IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderUnrep bit is set when unexpected repeats
	// encountered.
	IPv6ExtensionHeaderUnrep IPv6ExtensionHeader = 1 << iota

	// IPv6ExtensionHeaderUnseq bit is set when unexpected sequencing
	// encountered.
	IPv6ExtensionHeaderUnseq IPv6ExtensionHeader = 1 << iota
)

// IPv6ExtensionHeader represents bit definitions for IPv6 Extension
// Header pseudo field. It indicates the presence of various IPv6
// extension headers in the packet header.
//
// Values of this type could be used as XM matching values.
type IPv6ExtensionHeader uint16

// The XM defines the flow match fields are described using
// the OpenFlow Extensible Match (OXM) format, which is a
// compact type-length-value (TLV) format.
type XM struct {
	// Match class that contains related match type
	Class XMClass
	// Class-specific value, identifying one of the
	// match types within the match class.
	Type  XMType
	Value XMValue
	Mask  XMValue
}

// readAllXM uses all avaiable bytes retrieved from the reader to
// unmarshal them to the list of extensible matchers. The caller
// responsible of passing limited reader to prevent from read of
// unnecessary data.
func readAllXM(r io.Reader, xms *[]XM) (int64, error) {
	// Read all available bytes from the reader, they will be
	// used to unmarshal them into the list of extensible matchers.
	buf, err := ioutil.ReadAll(r)
	n := int64(len(buf))

	if err != nil {
		return n, err
	}

	rbuf := bytes.NewBuffer(buf)

	for rbuf.Len() > 7 {
		var xm XM

		_, err = xm.ReadFrom(rbuf)
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
	var length uint8

	n, err = encoding.ReadFrom(
		r, &xm.Class, &xm.Type, &length)
	if err != nil {
		return
	}

	hasmask := (xm.Type & 1) == 1
	xm.Type >>= 1

	var m int64

	xm.Value, xm.Mask = make(XMValue, length), nil
	m, err = encoding.ReadFrom(r, &xm.Value)
	n += m

	if err != nil {
		return
	}

	if hasmask {
		length /= 2
		xm.Mask = make(XMValue, length)

		copy(xm.Mask, xm.Value[length:])
		xm.Value = xm.Value[:length]
	}

	return
}

// WriteTo implements io.WriterTo interface. It serializes the Openflow
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

type XMValue []byte

func (val XMValue) UInt32() (v uint32) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

func (val XMValue) UInt16() (v uint16) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

func (val XMValue) UInt8() (v uint8) {
	encoding.ReadFrom(bytes.NewBuffer(val), &v)
	return
}

type XMExperimenterHeader struct {
	XM           XM
	Experimenter uint32
}

const (
	// OpenFlow 1.1 match type MT_STANDARD is deprecated.
	MatchTypeStandard MatchType = iota

	// OpenFlow Extensible Match type.
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
// bytes from the given reader to the Match structure.
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

	limrd := io.LimitReader(r, int64(length))
	nn, err = readAllXM(limrd, &m.Fields)

	return n + nn, err
}

// WriteTo implements io.WriterTo interface. It serializes the Match
// structure to the given writer.
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

	if length%8 != 0 {
		_, err = buf.Write(make(pad, 8-length%8))
		if err != nil {
			return
		}
	}

	return encoding.WriteTo(
		w, m.Type, uint16(length), buf.Bytes())
}
