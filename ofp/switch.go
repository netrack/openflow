package ofp

import (
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// Capability represents a switch capabilities bitmap.
type Capability uint32

const (
	// CapabilityFlowStats is set when flow statistics is supported
	// by the datapath.
	CapabilityFlowStats Capability = 1 << iota

	// CapabilityTableStats is set when table statistics is supported
	// by the datapath.
	CapabilityTableStats

	// CapabilityPortStats is set when port statistics is supported
	// by the datapath.
	CapabilityPortStats

	// CapabilityGroupStats is set when group statistics is supported
	// by the datapath.
	CapabilityGroupStats

	// CapabilityIPReasm is set when datapath can reassemble IP
	// fragments.
	CapabilityIPReasm

	// CapabilityQueueStats is set when queue statistics is supported
	// by the datapath.
	CapabilityQueueStats

	// CapabilityPortBlocked is set when switch can block looping ports.
	//
	// This bit indicates that a switch protocol outside of OpenFlow, such
	// as 802.1D Spanning Tree, will detect topology loops and block ports
	// to prevent packet loops. If this bit is not set, in most cases the
	// controller should implement a mechanism to prevent packet loops.
	CapabilityPortBlocked Capability = 1 << 8
)

// ConfigFlag represents a switch configuration flags.
type ConfigFlag uint16

const (
	// ConfigFlagFragNormal is set, when no special handling for IP
	// fragments is configured on the switch.
	//
	// This type of fragments handling means that an attempt should be
	// made to pass the fragments through the OpenFlow tables. If any
	// field is not present (e.g., the TCP/UDP ports didn't fit), then
	// the packet should not match any entry that has that field set.
	ConfigFlagFragNormal ConfigFlag = iota

	// ConfigFlagFragDrop is set, when switch drops IP fragments.
	ConfigFlagFragDrop

	// ConfigFlagFragReasm is set, when switch reassembles IP fragments.
	ConfigFlagFragReasm

	// ConfigFlagFragMask is an IP fragment reassemble mask.
	ConfigFlagFragMask
)

// SwitchFeatures is mesage used as a response to the controller on
// feature request message.
//
// For exampe, to retireve the switch configuration, the following
// request can be sent:
//
//	req := of.NewRequest(of.TypeFeaturesRequest, nil)
type SwitchFeatures struct {
	// Datapath unique ID. The lower 48-bits are for a MAC address,
	// while the upper 16-bits are implementer-defined.
	DatapathID uint64

	// NumBuffers is a number of max packets buffered at once.
	NumBuffers uint32

	// NumberTables is a number of tables supported by datapath.
	NumTables uint8

	// AuxiliaryID identifies auxiliary connections.
	AuxiliaryID uint8

	// Capabilities is a bitmap of supported capabilities.
	Capabilities Capability

	// Reserved bytes.
	Reserved uint32
}

// WriteTo implements io.WriterTo interface. It serializes the switch
// features into the wire format.
func (s *SwitchFeatures) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		s.DatapathID, s.NumBuffers, s.NumTables, s.AuxiliaryID,
		pad2{}, s.Capabilities, s.Reserved)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// switch features from the wire format.
func (s *SwitchFeatures) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r,
		&s.DatapathID, &s.NumBuffers, &s.NumTables, &s.AuxiliaryID,
		&defaultPad2, &s.Capabilities, &s.Reserved)
}

// SwitchConfig is a message used as a response to the controller
// on switch configuration request message.
//
// To retireve the switch configuration, the following request can
// be sent:
//
//	req := of.NewRequest(of.TypeGetConfigRequest, nil)
type SwitchConfig struct {
	// Flags is bitmap of switch configuration flags.
	Flags ConfigFlag

	// MissSendLength is a max bytes of packet that datapath
	// should send to the controller.
	MissSendLength uint16
}

// WriteTo implements io.WriterTo interface. It serializes the
// switch configuration into the wire format.
func (sc *SwitchConfig) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *sc)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the switch configuration from the wire format.
func (sc *SwitchConfig) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &sc.Flags, &sc.MissSendLength)
}
