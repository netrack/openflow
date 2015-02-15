package ofp13

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const VERSION uint8 = 0x04

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
	var length uint16

	err := binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&p.BufferID, &length, &p.Reason, &p.TableID, &p.Cookie,
	})

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

// Send packet (controller -> datapath)
type PacketOut struct {
	// ID assigned by datapath (NO_BUFFER if none).
	// The BufferID is the same given in the PacketIn message.
	BufferID uint32
	// Packet's input port or OFPP_CONTROLLER
	InPort PortNo
	// Action list

	// TODO: modify interface type
	Actions []interface {
		Write(io.Writer) error
	}
}

func (p *PacketOut) Write(w io.Writer) error {
	var buf bytes.Buffer

	for _, action := range p.Actions {
		err := action.Write(&buf)
		fmt.Println("here!!!", err)
		if err != nil {
			return err
		}
	}

	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		p.BufferID, p.InPort, uint16(buf.Len()), pad6{}, buf.Bytes(),
	})
}
