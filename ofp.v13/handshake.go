package ofp

import (
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	// Bitmap of version supported
	HET_VERSIONBITMAP HelloElemType = 1
)

// Hello elements types
type HelloElemType uint16

// This message includes zero or more hello elements
// having variable size. Unknown elements types must
// be ignored/skipped, to allow for future extensions.
type Hello struct {
	// Elements is a set of hello elements, containing optional
	// data to inform the initial handshake of the connection.
	Elements []HelloElemHeader
}

// Bytes serializes the message in binary format.
func (h *Hello) Bytes() []byte { return Bytes(h) }

// WriteTo implements io.WriterTo interface. It serializes
// the message in binary format.
func (h *Hello) WriteTo(w io.Writer) (int64, error) {
	return binary.Write(w, binary.BigEndian, h.Elements)
}

// ReadFrom implements io.ReaderFrom interface. It de-serializes
// the message from binary format.
func (h *Hello) ReadFrom(r io.Reader) (int64, error) {
	return binary.Read(r, binary.BigEndian, &h.Elements)
}

// Common header for all Hello Elements
type HelloElemHeader struct {
	Type HelloElemType
	Len  uint16
}

type HelloElemVersionBitmap struct {
	Type    HelloElemType
	Len     uint16
	Bitmaps []uint32
}

type ExperimenterHeader struct {
	Experimenter uint32
	ExpType      uint32
}

const (
	// ControllerRoleNoChange defines a request to not change the
	// current role.
	ControllerRoleNoChange ControllerRole = iota

	// ControllerRoleEqual defines a default role, full access.
	ControllerRoleEqual

	// ControllerRoleMaster defines a full access role, at most one master.
	ControllerRoleMaster

	// ControllerRoleSlave defines a read-only access role.
	ControllerRoleSlave
)

// ControllerRole is a role that controller wants to assume.
type ControllerRole uint32

type RoleRequest struct {
	Role         ControllerRole
	_            pad4
	GenerationID uint64
}

type AsyncConfig struct {
	PacketInMask    []uint32 //TODO
	PortStatusMask  []uint32
	FlowRemovedMask []uint32
}
