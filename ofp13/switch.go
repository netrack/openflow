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
	C_FRAG_NORMAL ConfigFlags = iota
	C_FRAG_DROP   ConfigFlags = 1 << 0
	C_FRAG_REASM  ConfigFlags = 1 << 1
	C_FRAG_MASK   ConfigFlags = iota
)

type ConfigFlags uint16

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

type SwitchConfig struct {
	Flags          ConfigFlags
	MissSendLength uint16
}
