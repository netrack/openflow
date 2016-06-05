package ofp

import (
	"io"
	"net"
	"strings"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	PF_10MB_HD  PortFeatures = 1 << iota
	PF_10MB_FD  PortFeatures = 1 << iota
	PF_100MB_HD PortFeatures = 1 << iota
	PF_100MB_FD PortFeatures = 1 << iota
	PF_1GB_HD   PortFeatures = 1 << iota
	PF_1GB_FD   PortFeatures = 1 << iota
	PF_10GB_FD  PortFeatures = 1 << iota
	PF_40GB_FD  PortFeatures = 1 << iota
	PF_100GB_FD PortFeatures = 1 << iota
	PF_1TB_FD   PortFeatures = 1 << iota
	PF_OTHER    PortFeatures = 1 << iota

	PF_COPPER     PortFeatures = 1 << iota
	PF_FIBER      PortFeatures = 1 << iota
	PF_AUTONEG    PortFeatures = 1 << iota
	PF_PAUSE      PortFeatures = 1 << iota
	PF_PAUSE_ASYM PortFeatures = 1 << iota
)

type PortFeatures uint32

func (f PortFeatures) String() string {
	// 10Mbps full-duplex copper
	var repr string

	switch {
	case f&PF_10MB_HD != 0:
		repr = "10 Mbps half-duplex"
	case f&PF_10MB_FD != 0:
		repr = "10 Mbps full-duplex"
	case f&PF_100MB_HD != 0:
		repr = "100 Mbps half-duplex"
	case f&PF_100MB_FD != 0:
		repr = "100 Mbps full-duplex"
	case f&PF_1GB_HD != 0:
		repr = "1 Gbps half-duplex"
	case f&PF_1GB_FD != 0:
		repr = "1 Gbps full-duplex"
	case f&PF_10GB_FD != 0:
		repr = "10 Gbps full-duplex"
	case f&PF_40GB_FD != 0:
		repr = "40 Gbps full-duplex"
	case f&PF_100GB_FD != 0:
		repr = "100 Gbps full-duplex"
	case f&PF_1TB_FD != 0:
		repr = "1 Tbps full-duplex"
	case f&PF_OTHER != 0:
		repr = "other"
	default:
		repr = "unknown"
	}

	switch {
	case f&PF_COPPER != 0:
		repr += " copper"
	case f&PF_FIBER != 0:
		repr += " fiber"
	case f&PF_AUTONEG != 0:
		repr += " autoneg"
	case f&PF_PAUSE != 0:
		repr += " pause"
	case f&PF_PAUSE_ASYM != 0:
		repr += " pause asym"
	default:
		repr += " unknown"
	}

	return repr
}

const (
	// Port is administratively down
	PC_PORT_DOWN PortConfig = 1 << iota

	// Drop all packets received by port
	PC_NO_RCV PortConfig = 1 << iota

	// Drop packets forwarded to port
	PC_NO_FWD PortConfig = 1 << iota

	// Do not send packet-in msgs for port
	PC_NO_PACKET_IN PortConfig = 1 << iota
)

type PortConfig uint32

func (c PortConfig) String() string {
	var repr []string

	if c&PC_PORT_DOWN != 0 {
		repr = append(repr, "DOWN")
	}

	if c&PC_NO_RCV != 0 {
		repr = append(repr, "NO_RCV")
	}

	if c&PC_NO_FWD != 0 {
		repr = append(repr, "NO_FWD")
	}

	if c&PC_NO_PACKET_IN != 0 {
		repr = append(repr, "NO_PACKET_IN")
	}

	if len(repr) == 0 {
		repr = append(repr, "UP")
	}

	return strings.Join(repr, ",")
}

const (
	// PS_LINK_DOWN bit indicates that the physical link is not present.
	PS_LINK_DOWN PortState = 1 << iota

	// PS_BLOCKED bit indicates that a switch protocol outside
	// of OpenFlow, such as 802.1D Spanning Tree, is preventing
	// the use of that port with OFPP_FLOOD.
	PS_BLOCKED PortState = 1 << iota

	// PS_LIVE indicates that port available for live for Fast Failover Group
	PS_LIVE PortState = 1 << iota
)

// Current state of the physical port. These are
// not configurable from the controller
type PortState uint32

func (s PortState) String() string {
	var repr []string

	if s&PS_LINK_DOWN != 0 {
		repr = append(repr, "LINK_DOWN")
	}

	if s&PS_BLOCKED != 0 {
		repr = append(repr, "BLOCKED")
	}

	if s&PS_LIVE != 0 {
		repr = append(repr, "LIVE")
	}

	if len(repr) == 0 {
		repr = append(repr, "LINK_UP")
	}

	return strings.Join(repr, ",")
}

const (
	// Send the packet out the input port. This reserved
	// port must be explicitly used in order to send back
	// out of the input port.
	P_IN_PORT PortNo = 0xfffffff8 + iota

	// Submit the packet to the first flow table. This
	// destination port can only be used in packet-out messages.
	P_TABLE PortNo = 0xfffffff8 + iota

	// Process with normal L2/L3 switching.
	P_NORMAL PortNo = 0xfffffff8 + iota

	// All physical ports in VLAN, except input port and
	// those blocked or link down.
	P_FLOOD PortNo = 0xfffffff8 + iota

	// All physical ports except input port.
	P_ALL PortNo = 0xfffffff8 + iota

	// Send to controller.
	P_CONTROLLER PortNo = 0xfffffff8 + iota

	// Local openflow "port".
	P_LOCAL PortNo = 0xfffffff8 + iota

	// Wildcard port used only for flow mod (delete) and flow
	// stats requests. Selects all flows regardless of output port
	// (including flows with no output port).
	P_ANY PortNo = 0xffffffff

	// Maximum number of physical and logical switch ports
	P_MAX PortNo = 0xffffff00
)

type PortNo uint32

const MAX_PORT_NAME_LEN = 16

// The port description request MP_PORT_DESCRIPTION enables the
// controller to get a description of all the ports in the system
// that support OpenFlow. The request body is empty. The reply
// body consists of an array of the Port
type Port struct {
	PortNo PortNo
	HWAddr net.HardwareAddr
	Name   []byte

	Config PortConfig
	State  PortState

	// Current features
	Curr PortFeatures
	// Features being advertised by the port
	Advertised PortFeatures
	// Features supported by the port
	Supported PortFeatures
	// Features advertised by peer
	Peer PortFeatures

	// Current port bitrate in kbps
	CurrSpeed uint32
	// Max port bitrate in kbps
	MaxSpeed uint32
}

func (p *Port) ReadFrom(r io.Reader) (int64, error) {
	p.HWAddr = make(net.HardwareAddr, 6)
	p.Name = make([]byte, MAX_PORT_NAME_LEN)

	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&p.PortNo,
		&pad4{},
		&p.HWAddr,
		&pad2{},
		&p.Name,
		&p.Config,
		&p.State,
		&p.Curr,
		&p.Advertised,
		&p.Supported,
		&p.Peer,
		&p.CurrSpeed,
		&p.MaxSpeed,
	})
}

type Ports []Port

func (p *Ports) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	for err == nil {
		var port Port
		nn, err = port.ReadFrom(r)
		n += nn

		if err == io.EOF {
			err = nil
			return
		}

		*p = append(*p, port)
	}

	return
}

type PortMod struct {
	PortNo    PortNo
	_         pad4
	HWAddr    net.HardwareAddr
	_         pad2
	Config    PortConfig
	Mask      PortConfig
	Advertise PortFeatures
	_         pad4
}

const (
	// The port was added
	PR_ADD PortReason = iota

	// The port was removed
	PR_DELETE

	// Some attribute of the port has changed
	PR_MODIFY
)

type PortReason uint8

type PortStatus struct {
	Reason PortReason
	_      pad7
	Desc   Port
}

type PortStatsRequest struct {
	PortNo PortNo
	_      pad4
}

type PortStats struct {
	PortNo       PortNo
	_            pad4
	RxPackets    uint64
	TxPackets    uint64
	RxBytes      uint64
	TxBytes      uint64
	RxDropped    uint64
	TxDropped    uint64
	RxErrors     uint64
	TxErrors     uint64
	RxFrameErr   uint64
	RxOverErr    uint64
	RxCrcErr     uint64
	Collisions   uint64
	DurationSec  uint32
	DurationNSec uint32
}
