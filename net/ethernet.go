package net

import (
	"net"
)

type EtherType uint16

type Ethernet struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLAN      VLAN
	EtherType EtherType
	Payload   []byte
}

type VLAN struct {
	TPID uint16
	TCI  uint16
}

func (v *VLAN) PCP() int {
	return (v.TCI & 0xe000) >> 13
}

func (v *VLAN) DEI() int {
	return (v.TCI & 0x1000) >> 12
}

func (v *VLAN) VID() int {
	return (v.TCI & 0x0fff)
}
