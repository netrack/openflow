package ofp13

import (
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// Flow statistics
	C_FLOW_STATS Capabilities = 1 << iota

	// Table statistics
	C_TABLE_STATS Capabilities = 1 << iota

	// Port statistics
	C_PORT_STATS Capabilities = 1 << iota

	// Group statistics
	C_GROUP_STATS Capabilities = 1 << iota

	// Can reassemble IP fragments
	C_IP_REASM Capabilities = 1 << iota

	// Queue statistics
	C_QUEUE_STATS Capabilities = 1 << iota

	// Switch will block looping ports
	C_PORT_BLOCKED Capabilities = 1 << 8
)

type Capabilities uint32

const (
	// No special handling for fragments
	C_FRAG_NORMAL ConfigFlags = iota

	// Drop fragments
	C_FRAG_DROP ConfigFlags = 1 << 0

	// Reassemble (only if C_IP_REASM set)
	C_FRAG_REASM ConfigFlags = 1 << 1
	C_FRAG_MASK  ConfigFlags = iota
)

// Configuration flags
type ConfigFlags uint16

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
	Capabilities Capabilities
	// Reserved
	Reserved uint32
}

func (s *SwitchFeatures) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		s.DatapathID,
		s.NumBuffers,
		s.NumTables,
		s.AuxiliaryID,
		pad2{},
		s.Capabilities,
		s.Reserved,
	})
}

func (s *SwitchFeatures) ReadFrom(r io.Reader) (int64, error) {
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&s.DatapathID,
		&s.NumBuffers,
		&s.NumTables,
		&s.AuxiliaryID,
		&pad2{},
		&s.Capabilities,
		&s.Reserved,
	})
}

// The controller is able to set and query configuration
// parameters in the switch with the T_SET_CONFIG and
// T_GET_CONFIG_REQUEST messages, respectively. The switch
// responds to a configuration request with an T_GET_CONFIG_REPLY message;
// it does not reply to a request to set the configuration.
type SwitchConfig struct {
	// Cofigurion flags
	Flags ConfigFlags
	// Max bytes of packet that datapath should send
	// to the controller.
	MissSendLength uint16
}

func (sc *SwitchConfig) ReadFrom(r io.Reader) (int64, error) {
	return binary.Read(r, binary.BigEndian, sc)
}

func (sc *SwitchConfig) WriteTo(w io.Writer) (int64, error) {
	return binary.Write(w, binary.BigEndian, sc)
}
