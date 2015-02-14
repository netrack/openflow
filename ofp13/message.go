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
	BufferID   uint32
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
	GenerationID uint64
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

// Packet received on port (datapath -> controller)
type PacketIn struct {
	// ID assigned by datapath
	BufferID uint32
	// Full length of frame
	//TotalLength uint16
	// Reason packet is being sent
	Reason PacketInReason
	// ID of the table that was looked up
	TableID Table
	// Cookie of the flow entry that was looked up
	Cookie uint64
	// Packet metadata. Variable size.
	Match Match
	// Followed by Link Layer packet frame
}

func (p *PacketIn) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &p.BufferID)
	if err != nil {
		return err
	}

	var length uint16
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.Reason)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.TableID)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.Cookie)
	if err != nil {
		return err
	}

	err = p.Match.Read(r)
	if err != nil {
		return err
	}

	var padding pad2
	return binary.Read(r, binary.BigEndian, &padding)
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
