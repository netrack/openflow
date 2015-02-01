package ofp13

import (
	"encoding/binary"
	"io"
)

const VERSION uint8 = 0x04

const (
	// Immutable messages
	T_HELLO Type = iota
	T_ERROR
	T_ECHO_REQUEST
	T_ECHO_REPLY
	T_EXPERIMENTER

	// Switch configuration messages
	T_FEATURES_REQUEST
	T_FEATURES_REPLY
	T_GET_CONFIG_REQUEST
	T_GET_CONFIG_REPLY
	T_SET_CONFIG

	// Asynchronous messages
	T_PACKET_IN
	T_FLOW_REMOVED
	T_PORT_STATUS

	// Controller command messages
	T_PACKET_OUT
	T_FLOW_MOD
	T_GROUP_MOD
	T_PORT_MOD
	T_TABLE_MOD

	// Multipart messages
	T_MULTIPART_REQUEST
	T_MULTIPART_REPLY

	// Queue configuration messages
	T_QUEUE_GET_CONFIG_REQUEST
	T_QUEUE_GET_CONFIG_REPLY

	// Controller role change request messages
	T_ROLE_REQUEST
	T_ROLE_REPLY

	// Asynchronous message configuration
	T_ASYNC_REQUEST
	T_ASYNC_REPLY
	T_SET_ASYNC

	// Meters and rate limiters configuration messages
	T_METER_MOD
)

type Type uint8

const (
	L_HELLO uint16 = 8
	L_ERROR
	L_ECHO_REQUEST
	L_ECHO_REPLY
	L_EXPERIMENTER

	L_FEATURES_REQUEST
	L_FEATURES_REPLY
	L_GEL_CONFIG_REQUEST
	L_GEL_CONFIG_REPLY
	L_SEL_CONFIG

	L_PACKEL_IN
	L_FLOW_REMOVED
	L_PORL_STATUS

	L_PACKEL_OUT
	L_FLOW_MOD
	L_GROUP_MOD
	L_PORL_MOD
	L_TABLE_MOD

	L_MULTIPARL_REQUEST
	L_MULTIPARL_REPLY

	L_QUEUE_GEL_CONFIG_REQUEST
	L_QUEUE_GEL_CONFIG_REPLY

	L_ROLE_REQUEST
	L_ROLE_REPLY

	L_ASYNC_REQUEST
	L_ASYNC_REPLY
	L_SEL_ASYNC

	L_METER_MOD
)

const (
	CR_ROLE_NOCHANGE ControllerRole = iota
	CR_ROLE_EQUAL
	CR_ROLE_MASTER
	CR_ROLE_SLAVE
)

type ControllerRole uint32

type Header struct {
	Version uint8
	Type    Type
	Length  uint16
	Xid     uint32
}

func (h *Header) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h)
}

func (h *Header) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, h)
}

type PacketOut struct {
	Header     Header
	BufferId   uint32
	InPort     PortNo
	ActionsLen uint16
	_          pad6
	Actions    []ActionHeader
}

type RoleRequest struct {
	Header       Header
	Role         ControllerRole
	_            pad4
	GenerationId uint64
}

type AsyncConfig struct {
	Header          Header
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

type PacketIn struct {
	Header      Header
	BufferId    uint32
	TotalLength uint16
	Reason      PacketInReason
	TableId     Table
	Cookie      uint64
	Match       Match
}

const (
	RR_IDLE_TIMEOUT FlowRemovedReason = iota
	RR_HARD_TIMEOUT
	RR_DELETE
	RR_GROUP_DELETE
)

type FlowRemovedReason uint8

type FlowRemoved struct {
	Header       Header
	Cookie       uint64
	Priority     uint16
	Reason       FlowRemovedReason
	TableId      Table
	DurationSec  uint32
	DurationNSec uint32
	IdleTimeout  uint16
	HardTimeout  uint16
	PacketCount  uint64
	ByteCount    uint64
	Match        Match
}

const (
	PR_ADD PortReason = iota
	PR_DELETE
	PR_MODIFY
)

type PortReason uint8

type PortStatus struct {
	Header Header
	Reason PortReason
	_      pad7
	Desc   Port
}

const (
	HET_VERSIONBITMAP HelloElemType = 1
)

type HelloElemType uint16

type Hello struct {
	Header   Header
	Elements []HelloElemHeader
}

func NewHello() *Hello {
	return &Hello{Header{VERSION, T_HELLO, L_HELLO, 0}, []HelloElemHeader{}}
}

func (h *Hello) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, h.Header)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, h.Elements)
}

func (h *Hello) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &h.Header)
	if err != nil {
		return err
	}

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
	Header       Header
	Experimenter uint32
	ExpType      uint32
}

type pad1 [1]uint8
type pad2 [2]uint8
type pad3 [3]uint8
type pad4 [4]uint8
type pad5 [5]uint8
type pad6 [6]uint8
type pad7 [7]uint8
