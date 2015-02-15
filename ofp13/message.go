package ofp13

import (
	"bytes"
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const VERSION uint8 = 0x04

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
	// The elements field is a set of hello elements,
	// containing optional data to inform the initial handshake
	// of the connection
	Elements []HelloElemHeader
}

func (h *Hello) Bytes() []byte { return Bytes(h) }

func (h *Hello) WriteTo(w io.Writer) (int64, error) {
	return binary.Write(w, binary.BigEndian, h.Elements)
}

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

func (p *PacketIn) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint16

	n, err = binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&p.BufferID, &length, &p.Reason, &p.TableID, &p.Cookie,
	})

	if err != nil {
		return
	}

	var nn int64

	nn, err = p.Match.ReadFrom(r)
	n += nn

	if err != nil {
		return
	}

	var padding pad2
	nn, err = binary.Read(r, binary.BigEndian, &padding)
	return n + nn, err
}

// Send packet (controller -> datapath)
type PacketOut struct {
	// ID assigned by datapath (NO_BUFFER if none).
	// The BufferID is the same given in the PacketIn message.
	BufferID uint32
	// Packet's input port or OFPP_CONTROLLER
	InPort PortNo
	// Action list
	Actions Actions
}

func (p *PacketOut) Bytes() (b []byte) { return Bytes(p) }

func (p *PacketOut) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	for _, action := range p.Actions {
		n, err = action.WriteTo(&buf)
		if err != nil {
			return
		}
	}

	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		p.BufferID, p.InPort, uint16(buf.Len()), pad6{}, buf.Bytes(),
	})
}
