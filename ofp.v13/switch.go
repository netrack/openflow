package ofp

import (
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// Flow statistics
	CapabilityFlowStats Capability = 1 << iota

	// Table statistics
	CapabilityTableStats Capability = 1 << iota

	// Port statistics
	CapabilityPortStats Capability = 1 << iota

	// Group statistics
	CapabilityGroupStats Capability = 1 << iota

	// Can reassemble IP fragments
	CapabilityIPReasm Capability = 1 << iota

	// Queue statistics
	CapabilityQueueStats Capability = 1 << iota

	// Switch will block looping ports
	CapabilityPortBlocked Capability = 1 << 8
)

type Capability uint32

const (
	// No special handling for fragments
	ConfigFlagFragNormal ConfigFlag = iota

	// Drop fragments
	ConfigFlagFragDrop

	// Reassemble (only if C_IP_REASM set)
	ConfigFlagFragReasm

	// Fragment reassemble mask
	ConfigFlagFragMask
)

// Configuration flags
type ConfigFlag uint16

// The C_FRAG_* flags indicate whether IP fragments should be
// treated normally, dropped, or reassembled. "Normal" handling
// of fragments means that an attempt should be made to pass the
// fragments through the OpenFlow tables. If any field is not
// present (e.g., the TCP/UDP ports didn’t fit), then the
// packet should not match any entry that has that field set.

// Switch features.
type SwitchFeatures struct {
	// Datapath unique ID. The lower 48-bits are for a MAC address,
	// while the upper 16-bits are implementer-defined
	DatapathID uint64

	// Max packets buffered at once
	NumBuffers uint32

	// Number of tables supported by datapath
	NumTables uint8

	// Identify auxiliary connections
	AuxiliaryID uint8

	// Bitmap of support Capabilities
	Capabilities Capability

	// Reserved
	Reserved uint32
}

func (s *SwitchFeatures) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		s.DatapathID,
		s.NumBuffers,
		s.NumTables,
		s.AuxiliaryID,
		pad2{},
		s.Capabilities,
		s.Reserved,
	)
}

func (s *SwitchFeatures) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r,
		&s.DatapathID,
		&s.NumBuffers,
		&s.NumTables,
		&s.AuxiliaryID,
		&pad2{},
		&s.Capabilities,
		&s.Reserved,
	)
}

// The controller is able to set and query configuration
// parameters in the switch with the T_SET_CONFIG and
// T_GET_CONFIG_REQUEST messages, respectively. The switch
// responds to a configuration request with an T_GET_CONFIG_REPLY message;
// it does not reply to a request to set the configuration.
type SwitchConfig struct {
	// Cofigurion flags
	Flags ConfigFlag
	// Max bytes of packet that datapath should send
	// to the controller.
	MissSendLength uint16
}

func (sc *SwitchConfig) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &sc.Flags, &sc.MissSendLength)
}

func (sc *SwitchConfig) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, *sc)
}
