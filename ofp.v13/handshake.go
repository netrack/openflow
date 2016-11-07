package ofp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/encoding"
)

const (
	// Bitmap of version supported
	HelloElemTypeVersionBitmap HelloElemType = 1
)

// Hello elements types
type HelloElemType uint16

type HelloElem interface {
	encoding.ReadWriter

	Type() HelloElemType
}

var helloElemMap = map[HelloElemType]encoding.ReaderMaker{
	HelloElemTypeVersionBitmap: encoding.ReaderMakerOf(
		HelloElemVersionBitmap{}),
}

type helloElemHdr struct {
	Type HelloElemType
	Len  uint16
}

// helloElemLen defines a length of the helloElem header.
const helloElemLen = 4

type HelloElems []HelloElem

func (h HelloElems) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	for _, elem := range h {
		_, err := elem.WriteTo(&buf)
		if err != nil {
			return 0, err
		}
	}

	return encoding.WriteTo(w, buf.Bytes())
}

func (h HelloElems) ReadFrom(r io.Reader) (int64, error) {
	var helloElemType HelloElemType

	rm := func() (io.ReaderFrom, error) {
		if rm, ok := helloElemMap[helloElemType]; ok {
			rd, err := rm.MakeReader()
			h = append(h, rd.(HelloElem))
			return rd, err
		}

		format := "ofp: unknown hello element type: '%x'"
		return nil, fmt.Errorf(format, helloElemType)
	}

	return encoding.ScanFrom(r, &helloElemType,
		encoding.ReaderMakerFunc(rm))
}

// This message includes zero or more hello elements
// having variable size. Unknown elements types must
// be ignored/skipped, to allow for future extensions.
type Hello struct {
	// Elements is a set of hello elements, containing optional
	// data to inform the initial handshake of the connection.
	Elements HelloElems
}

// Bytes serializes the message in binary format.
func (h *Hello) Bytes() []byte { return Bytes(h) }

// WriteTo implements io.WriterTo interface. It serializes
// the message in binary format.
func (h *Hello) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, h.Elements)
}

// ReadFrom implements io.ReaderFrom interface. It de-serializes
// the message from binary format.
func (h *Hello) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &h.Elements)
}

type HelloElemVersionBitmap struct {
	Bitmaps []uint32
}

func (h *HelloElemVersionBitmap) Type() HelloElemType {
	return HelloElemTypeVersionBitmap
}

func (h *HelloElemVersionBitmap) WriteTo(w io.Writer) (int64, error) {
	// Calculate the padding used to align the list of bitmaps.
	maplen := len(h.Bitmaps) * 4
	padding := make([]byte, (helloElemLen+maplen)%8)

	// Total length of the element includes the length of the
	// optional header, length of the bitmaps and the padding
	// used to align the message to 64-bits border.
	totalLen := helloElemLen + maplen + len(padding)

	// Compose the healder of the element and marshal it
	// altogether with a version bitmaps and required padding.
	header := helloElemHdr{h.Type(), uint16(totalLen)}
	return encoding.WriteTo(w, header, h.Bitmaps, padding)
}

func (h *HelloElemVersionBitmap) ReadFrom(r io.Reader) (int64, error) {
	// Read the header of the version bitmaps to retireve
	// the total length of the message.
	var header helloElemHdr
	n, err := encoding.ReadFrom(r, &header)
	if err != nil {
		return n, err
	}

	// Calculate the length of the list of version bitmaps.
	bodyLen := header.Len - helloElemLen
	limrd := io.LimitReader(r, int64(bodyLen))

	// Allocate required amount of memory used to fit all
	// element version bitmaps (assuming that uint32 takes
	// 4 bytes of the memory).
	h.Bitmaps = make([]uint32, (bodyLen)/4)
	nn, err := encoding.ReadFrom(limrd, &h.Bitmaps)
	return n + nn, err
}

type Experimenter struct {
	Experimenter uint32
	ExpType      uint32
}

func (e *Experimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e.Experimenter, e.ExpType)
}

func (e *Experimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &e.Experimenter, &e.ExpType)
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
	GenerationID uint64
}

func (rr *RoleRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, rr.Role, pad4{}, rr.GenerationID)
}

func (rr *RoleRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &rr.Role, &defaultPad4, &rr.GenerationID)
}

const AsyncConfigMaskLen = 8

type AsyncConfig struct {
	PacketInMask    [AsyncConfigMaskLen]PacketInReason
	PortStatusMask  [AsyncConfigMaskLen]PortReason
	FlowRemovedMask [AsyncConfigMaskLen]FlowRemovedReason
}

func (a *AsyncConfig) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, a.PacketInMask,
		a.PortStatusMask, a.FlowRemovedMask)
}

func (a *AsyncConfig) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &a.PacketInMask,
		&a.PortStatusMask, &a.FlowRemovedMask)
}
