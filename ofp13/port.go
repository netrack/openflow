package ofp13

import (
	"net"
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

const (
	PC_PORT_DOWN    PortConfig = 1 << iota
	PC_NO_RCV       PortConfig = 1 << iota
	PC_NO_FWD       PortConfig = 1 << iota
	PC_NO_PACKET_IN PortConfig = 1 << iota
)

type PortConfig uint32

const (
	PS_LINK_DOWN PortState = 1 << iota
	PS_BLOCKED   PortState = 1 << iota
	PS_LIVE      PortState = 1 << iota
)

type PortState uint32

const (
	// Maximum number of physical and logical switch ports
	P_MAX PortNo = 0xffffff00

	// Send the packet out the input port. This reserved port must be
	// explicitly used in order to send back out of the input port.
	P_IN_PORT    PortNo = 0xfffffff8 + iota
	P_TABLE      PortNo = 0xfffffff8 + iota
	P_NORMAL     PortNo = 0xfffffff8 + iota
	P_FLOOD      PortNo = 0xfffffff8 + iota
	P_ALL        PortNo = 0xfffffff8 + iota
	P_CONTROLLER PortNo = 0xfffffff8 + iota
	P_LOCAL      PortNo = 0xfffffff8 + iota
	P_ANY        PortNo = 0xfffffff8 + iota
)

type PortNo uint32

type Port struct {
	PortNo PortNo
	HWAddr net.HardwareAddr
	Name   []byte

	Config PortConfig
	State  PortState

	Curr       PortFeatures
	Advertised PortFeatures
	Supported  PortFeatures
	Peer       PortFeatures

	CurrSpeed uint32
	MaxSpeed  uint32
}

type PortMod struct {
	Header    Header
	PortNo    PortNo
	HWAddr    net.HardwareAddr
	Config    PortConfig
	Mask      PortConfig
	Advertise PortFeatures
}
