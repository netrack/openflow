package ofp

import (
	"fmt"
	"io"
	"net"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	PortFeature10MbitHalfDuplex  PortFeature = 1 << iota
	PortFeature10MbitFullDuplex  PortFeature = 1 << iota
	PortFeature100MbitHalfDuplex PortFeature = 1 << iota
	PortFeature100MbitFullDuplex PortFeature = 1 << iota
	PortFeature1GbitHalfDuplex   PortFeature = 1 << iota
	PortFeature1GbitFullDuplex   PortFeature = 1 << iota
	PortFeature10GbitFullDuplex  PortFeature = 1 << iota
	PortFeature40GbitFullDuplex  PortFeature = 1 << iota
	PortFeature100GbitFullDuplex PortFeature = 1 << iota
	PortFeature1TbitFullDuplex   PortFeature = 1 << iota
	PortFeatureOther             PortFeature = 1 << iota

	PortFeatureCopper    PortFeature = 1 << iota
	PortFeatureFiber     PortFeature = 1 << iota
	PortFeatureAutoneg   PortFeature = 1 << iota
	PortFeaturePause     PortFeature = 1 << iota
	PortFeaturePauseAsym PortFeature = 1 << iota
)

const (
	portFeatureSpeedMask  = 0x000007ff
	portFeatureMediumMask = 0x0000f800
)

type PortFeature uint32

var portFeaturesText = map[PortFeature]string{
	PortFeature10MbitHalfDuplex:  "10 Mbps half-duplex",
	PortFeature10MbitFullDuplex:  "10 Mbps full-duplex",
	PortFeature100MbitHalfDuplex: "100 Mbps half-duplex",
	PortFeature100MbitFullDuplex: "100 Mbps full-duplex",
	PortFeature1GbitHalfDuplex:   "1 Gbps half-duplex",
	PortFeature1GbitFullDuplex:   "1 Gbps full-duplex",
	PortFeature10GbitFullDuplex:  "10 Gbps full-duplex",
	PortFeature40GbitFullDuplex:  "40 Gbps full-duplex",
	PortFeature100GbitFullDuplex: "100 Gbps full-duplex",
	PortFeature1TbitFullDuplex:   "1 Tbps full-duplex",
	PortFeatureOther:             "other",
	PortFeatureCopper:            "copper",
	PortFeatureFiber:             "fiber",
	PortFeatureAutoneg:           "autoneg",
	PortFeaturePause:             "pause",
	PortFeaturePauseAsym:         "pause asym",
}

func (f PortFeature) String() string {
	speed, ok := portFeaturesText[f&portFeatureSpeedMask]
	if !ok {
		speed = "unknown"
	}

	medium, ok := portFeaturesText[f&portFeatureMediumMask]
	if !ok {
		medium = "unknown"
	}

	return fmt.Sprintf("%s %s", speed, medium)
}

const (
	// Port is administratively down
	PortConfigDown PortConfig = 1 << iota

	// Drop all packets received by port
	PortConfigNoRcv PortConfig = 1 << iota

	// Drop packets forwarded to port
	PortConfigNoFwd PortConfig = 1 << iota

	// Do not send packet-in msgs for port
	PortConfigNoPacketIn PortConfig = 1 << iota
)

type PortConfig uint32

var portConfigText = map[PortConfig]string{
	PortConfigDown:       "down",
	PortConfigNoRcv:      "no recv",
	PortConfigNoFwd:      "no fwd",
	PortConfigNoPacketIn: "no packet in",
}

func (c PortConfig) String() string {
	if text, ok := portConfigText[c]; ok {
		return text
	}

	return "up"
}

const (
	// PS_LINK_DOWN bit indicates that the physical link is not present.
	PortStateLinkDown PortState = 1 << iota

	// PS_BLOCKED bit indicates that a switch protocol outside
	// of OpenFlow, such as 802.1D Spanning Tree, is preventing
	// the use of that port with OFPP_FLOOD.
	PortStateBlocked PortState = 1 << iota

	// PS_LIVE indicates that port available for live for Fast Failover Group
	PortStateLive PortState = 1 << iota
)

// Current state of the physical port. These are
// not configurable from the controller
type PortState uint32

var portStateText = map[PortState]string{
	PortStateLinkDown: "link down",
	PortStateBlocked:  "blocked",
	PortStateLive:     "live",
}

func (s PortState) String() string {
	if text, ok := portStateText[s]; ok {
		return text
	}

	return "link up"
}

const (
	// Send the packet out the input port. This reserved
	// port must be explicitly used in order to send back
	// out of the input port.
	PortIn PortNo = 0xfffffff8 + iota

	// Submit the packet to the first flow table. This
	// destination port can only be used in packet-out messages.
	PortTable PortNo = 0xfffffff8 + iota

	// Process with normal L2/L3 switching.
	PortNormal PortNo = 0xfffffff8 + iota

	// All physical ports in VLAN, except input port and
	// those blocked or link down.
	PortFlood PortNo = 0xfffffff8 + iota

	// All physical ports except input port.
	PortAll PortNo = 0xfffffff8 + iota

	// Send to controller.
	PortController PortNo = 0xfffffff8 + iota

	// Local openflow "port".
	PortLocal PortNo = 0xfffffff8 + iota

	// Wildcard port used only for flow mod (delete) and flow
	// stats requests. Selects all flows regardless of output port
	// (including flows with no output port).
	PortAny PortNo = 0xffffffff

	// Maximum number of physical and logical switch ports
	PortMax PortNo = 0xffffff00
)

type PortNo uint32

const portNameLen = 16

// The port description request MP_PORT_DESCRIPTION enables the
// controller to get a description of all the ports in the system
// that support OpenFlow. The request body is empty. The reply
// body consists of an array of the Port
type Port struct {
	PortNo PortNo
	HWAddr net.HardwareAddr
	Name   string

	Config PortConfig
	State  PortState

	// Current features
	Curr PortFeature
	// Features being advertised by the port
	Advertised PortFeature
	// Features supported by the port
	Supported PortFeature
	// Features advertised by peer
	Peer PortFeature

	// Current port bitrate in kbps
	CurrSpeed uint32
	// Max port bitrate in kbps
	MaxSpeed uint32
}

func (p *Port) WriteTo(w io.Writer) (int64, error) {
	name := make([]byte, portNameLen)
	copy(name, []byte(p.Name))

	return encoding.WriteTo(w,
		p.PortNo,
		pad4{},
		p.HWAddr,
		pad2{},
		name,
		p.Config,
		p.State,
		p.Curr,
		p.Advertised,
		p.Supported,
		p.Peer,
		p.CurrSpeed,
		p.MaxSpeed,
	)
}

func (p *Port) ReadFrom(r io.Reader) (int64, error) {
	p.HWAddr = make(net.HardwareAddr, 6)
	var name [portNameLen]byte

	n, err := encoding.ReadFrom(r,
		&p.PortNo,
		&defaultPad4,
		&p.HWAddr,
		&defaultPad2,
		&name,
		&p.Config,
		&p.State,
		&p.Curr,
		&p.Advertised,
		&p.Supported,
		&p.Peer,
		&p.CurrSpeed,
		&p.MaxSpeed,
	)

	p.Name = string(name[:])
	return n, err
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
	HWAddr    net.HardwareAddr
	Config    PortConfig
	Mask      PortConfig
	Advertise PortFeature
}

func (p *PortMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		p.PortNo, pad4{},
		p.HWAddr, pad2{},
		p.Config,
		p.Mask,
		p.Advertise, pad4{},
	)
}

func (p *PortMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r,
		&p.PortNo, &defaultPad4,
		&p.HWAddr, &defaultPad2,
		&p.Config,
		&p.Mask,
		&p.Advertise, &defaultPad4,
	)
}

const (
	// The port was added
	PortReasonAdd PortReason = iota

	// The port was removed
	PortReasonDelete

	// Some attribute of the port has changed
	PortReasonModify
)

type PortReason uint8

type PortStatus struct {
	Reason PortReason
	Port   Port
}

func (p *PortStatus) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.Reason, pad7{}, &p.Port)
}

func (p *PortStatus) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.Reason, &defaultPad7, &p.Port)
}

type PortStatsRequest struct {
	PortNo PortNo
}

func (p *PortStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.PortNo, pad4{})
}

func (p *PortStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.PortNo, &defaultPad4)
}

type PortStats struct {
	PortNo       PortNo
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

func (p *PortStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		p.PortNo,
		pad4{},
		p.RxPackets,
		p.TxPackets,
		p.RxBytes,
		p.TxBytes,
		p.RxDropped,
		p.TxDropped,
		p.RxErrors,
		p.TxErrors,
		p.RxFrameErr,
		p.RxOverErr,
		p.RxCrcErr,
		p.Collisions,
		p.DurationSec,
		p.DurationNSec,
	)
}

func (p *PortStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r,
		&p.PortNo,
		&defaultPad4,
		&p.RxPackets,
		&p.TxPackets,
		&p.RxBytes,
		&p.TxBytes,
		&p.RxDropped,
		&p.TxDropped,
		&p.RxErrors,
		&p.TxErrors,
		&p.RxFrameErr,
		&p.RxOverErr,
		&p.RxCrcErr,
		&p.Collisions,
		&p.DurationSec,
		&p.DurationNSec,
	)
}
