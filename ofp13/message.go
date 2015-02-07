package ofp13

import (
	"encoding/binary"
	"io"
)

const VERSION uint8 = 0x04

type Reader interface {
	Read(io.Reader) error
}

type Writer interface {
	Write(io.Writer) error
}

type PacketOut struct {
	BufferId   uint32
	InPort     PortNo
	ActionsLen uint16
	_          pad6
	Actions    []ActionHeader
}

const (
	CR_ROLE_NOCHANGE ControllerRole = iota
	CR_ROLE_EQUAL
	CR_ROLE_MASTER
	CR_ROLE_SLAVE
)

type ControllerRole uint32

type RoleRequest struct {
	Role         ControllerRole
	_            pad4
	GenerationId uint64
}

type AsyncConfig struct {
	PacketInMask    []uint32 //TODO
	PortStatusMask  []uint32
	FlowRemovedMask []uint32
}

const (
	R_NO_MATCH PacketInReason = iota
	R_ACTION
	R_INVALID_TTL
)

type PacketInReason uint8

type packetin struct {
	BufferId    uint32
	TotalLength uint16
	Reason      PacketInReason
	TableId     Table
	Cookie      uint64
	Match       Match
	_           pad2
}

type PacketIn struct {
	packetin
	Data []byte
}

func (p *PacketIn) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &p.packetin)
	if err != nil {
		return err
	}

	p.Data = make([]byte, p.TotalLength)
	return binary.Read(r, binary.BigEndian, p.Data)
}

const (
	HET_VERSIONBITMAP HelloElemType = 1
)

type HelloElemType uint16

type Hello struct {
	Elements []HelloElemHeader
}

func (h *Hello) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h.Elements)
}

func (h *Hello) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &h.Elements)
}

type HelloElemHeader struct {
	Type   HelloElemType
	Length uint16
}

type HelloElemVersionBitmap struct {
	Type    HelloElemType
	Length  uint16
	Bitmaps []uint32
}

type ExperimenterHeader struct {
	Experimenter uint32
	ExpType      uint32
}
