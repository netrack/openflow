package net

import (
	"encoding/binary"
	"io"
	"net"
)

type EtherType uint16

type EthernetII struct {
	Dst       [6]byte
	Src       [6]byte
	EtherType EtherType
}

func (eth *EthernetII) HWDst() (hwaddr net.HardwareAddr) {
	for _, b := range eth.Dst {
		hwaddr = append(hwaddr, b)
	}
	return
}

func (eth *EthernetII) HWSrc() (hwaddr net.HardwareAddr) {
	for _, b := range eth.Src {
		hwaddr = append(hwaddr, b)
	}
	return
}

func (eth *EthernetII) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, eth)
}

type Ethernet8021q struct {
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLAN      VLAN
	EtherType EtherType
}

func (eth *Ethernet8021q) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, eth)
}

type VLAN struct {
	TPID uint16
	TCI  uint16
}

func (v *VLAN) PCP() int {
	return int((v.TCI & 0xe000) >> 13)
}

func (v *VLAN) DEI() int {
	return int((v.TCI & 0x1000) >> 12)
}

func (v *VLAN) VID() int {
	return int((v.TCI & 0x0fff))
}
