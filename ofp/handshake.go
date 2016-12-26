package ofp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// HelloElemType defines the type of optional hello elements
// used to inform the initial handshake of the connection.
type HelloElemType uint16

const (
	// HelloElemTypeVersionBitmap is a bitmap of version supported.
	HelloElemTypeVersionBitmap HelloElemType = 1
)

// HelloElem is an optional hello element.
type HelloElem interface {
	encoding.ReadWriter

	// Type returns the type of the element.
	Type() HelloElemType
}

// Mapping of the hello elements to the implementations.
var helloElemMap = map[HelloElemType]encoding.ReaderMaker{
	HelloElemTypeVersionBitmap: encoding.ReaderMakerOf(
		HelloElemVersionBitmap{}),
}

// helloElem defines a common header for all hello elements.
type helloElem struct {
	Type HelloElemType
	Len  uint16
}

// helloElemLen defines a length of the helloElem header.
const helloElemLen = 4

// HelloElems is a list of the optional hello elements.
type HelloElems []HelloElem

// WriteTo implements io.WriterTo interface. It serializes the list
// of hello message elements into the wire format.
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

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// list of hello elements from the wire format.
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

// Hello is a message used to perform an initial handshake right after
// establishing the connection. It is used to negotiate the version of
// the protocol to communicate.
//
// For example, to create a hello message in order to response on the
// incoming message from the datapath, the following request could be
// constructed to inform that controller supports only OpenFlow 1.3:
//
//	req := of.NewRequest(of.HelloType, &ofp.Hello{
//		ofp.HelloElems{&ofp.HelloElemVersionBitmap{
//			[]uint32{1<<4},
//		}},
//	})
type Hello struct {
	// Elements is a set of hello elements, containing optional
	// data to inform the initial handshake of the connection.
	Elements HelloElems
}

// WriteTo implements io.WriterTo interface. It serializes
// the hello message into the wire format.
func (h *Hello) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, h.Elements)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the message from the wire format.
func (h *Hello) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &h.Elements)
}

// HelloElemVersionBitmap is a bitmap of the supported versions.
//
// The bitmaps field indicates the set of versions of the OpenFlow
// switch protocol a device supports, and may be used during version
// negotiation.
//
// For example, if a switch supports only version 1.0 (version=0x01)
// and version 1.3 (version=0x04), the first bitmap would be set to
// 0x00000012.
type HelloElemVersionBitmap struct {
	Bitmaps []uint32
}

// Type returns the type of hello element.
func (h *HelloElemVersionBitmap) Type() HelloElemType {
	return HelloElemTypeVersionBitmap
}

// WriteTo implements io.WriterTo interface. It serializes the version
// bitmap into the wire format.
func (h *HelloElemVersionBitmap) WriteTo(w io.Writer) (int64, error) {
	// Calculate the padding used to align the list of bitmaps.
	maplen := len(h.Bitmaps) * 4
	padding := make([]byte, (helloElemLen+maplen)%8)

	// Total length of the element includes the length of the
	// optional header, length of the bitmaps and the padding
	// used to align the message to 64-bits border.
	totalLen := helloElemLen + maplen + len(padding)

	// Compose the header of the element and marshal it
	// altogether with a version bitmaps and required padding.
	header := helloElem{h.Type(), uint16(totalLen)}
	return encoding.WriteTo(w, header, h.Bitmaps, padding)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// version bitmap from the wire format.
func (h *HelloElemVersionBitmap) ReadFrom(r io.Reader) (int64, error) {
	// Read the header of the version bitmaps to retrieve
	// the total length of the message.
	var header helloElem
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
	h.Bitmaps = make([]uint32, bodyLen/4)
	nn, err := encoding.ReadFrom(limrd, &h.Bitmaps)
	return n + nn, err
}

// Experimenter is an experimenter message header.
type Experimenter struct {
	// Experimenter identifier. If most significant bit is set to
	// zero, the low-order bytes represent experimenter's IEEE OUI.
	// Otherwise the identifier is defined by ONF.
	Experimenter uint32

	// ExpType is a experimenter message type.
	ExpType uint32
}

// WriteTo implements io.WriterTo interface. It serializes the
// experimenter header into the wire format.
func (e *Experimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e.Experimenter, e.ExpType)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the experimenter header from the wire format.
func (e *Experimenter) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &e.Experimenter, &e.ExpType)
}

// ControllerRole is a role that controller wants to assume.
type ControllerRole uint32

const (
	// ControllerRoleNoChange is used to not change the current role.
	// This enable a controller to query the current role without
	// changing it.
	ControllerRoleNoChange ControllerRole = iota

	// ControllerRoleEqual is used to change the role to a default role,
	// full access.
	ControllerRoleEqual

	// ControllerRoleMaster is used to change the role to a full access
	// role, at most one master.
	ControllerRoleMaster

	// ControllerRoleSlave is used to change the role to a read-only
	// access role.
	ControllerRoleSlave
)

// RoleRequest is a message used by controller to change its role.
//
// For example, to request the current controller role, the following
// request could be sent:
//
//	req := of.NewRequest(of.TypeRoleRequest, &ofp.RoleRequest{
//		Role: ControllerRoleNoChange,
//	})
type RoleRequest struct {
	// Role is a new role that the controller wants to assume.
	Role ControllerRole

	// GenerationID is a master election generation identifier.
	//
	// If the controller is changing the role to either slave or
	// master, the switch must validate the generation identifier
	// to check for stale messages.
	//
	// If the validation fails, the switch must discard the role
	// request and return an error message with type
	// ErrTypeRoleRequestFailed and error code
	// ErrCodeRoleRequestFailedStale.
	GenerationID uint64
}

// WriteTo implements io.WriterTo interface. It serializes the role
// request into the wire format.
func (rr *RoleRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, rr.Role, pad4{}, rr.GenerationID)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the role
// request from the wire format.
func (rr *RoleRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &rr.Role, &defaultPad4, &rr.GenerationID)
}

// AsyncConfig is a message used to configure the switch to receive
// specific types of asynchronous messages.
//
// For example, to configure the switch to send notifications on
// addition of the new port only to the controllers playing the master
// role, while the modification of the port configuration to the slave
// controllers, the following request could be constructed:
//
//	req := of.NewRequest(of.TypeAsyncConfig, &ofp.AsyncConfig{
//		PortStatusMask: ofputil.Bitmap64(
//			// Master will receive PortStats message when
//			// a new port will be added.
//			ofputil.PortReasonBitmap(ofp.PortReasonAdd),
//
//			// Slave will receive PortStats message when
//			// an existing port will be modified.
//			ofputil.PortReasonBitmap(ofp.PortReasonModify),
//		),
//	})
type AsyncConfig struct {
	// PacketInMask is a bitmap of PortReason types.
	PacketInMask [2]uint32

	// PortStatusMask is a bitmap of PortStatus types.
	PortStatusMask [2]uint32

	// FlowRemovedMask is a bitmap of FlowRemoved types.
	FlowRemovedMask [2]uint32
}

// WriteTo implements io.WriterTo interface. It serializes the
// asynchronous configuration into the wire format.
func (a *AsyncConfig) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, a.PacketInMask,
		a.PortStatusMask, a.FlowRemovedMask)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the asynchronous configuration from the wire format.
func (a *AsyncConfig) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &a.PacketInMask,
		&a.PortStatusMask, &a.FlowRemovedMask)
}
