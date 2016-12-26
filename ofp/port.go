package ofp

import (
	"io"
	"net"
	"strings"

	"github.com/netrack/openflow/internal/encoding"
)

// PortFeature defines the features of port available in datapath.
type PortFeature uint32

const (
	// PortFeature10MbitHalfDuplex is set when 10Mb half-duplex rate is
	// supported.
	PortFeature10MbitHalfDuplex PortFeature = 1 << iota

	// PortFeature10MbitFullDuplex is set when 10Mbit full-duplex rate
	// is supported.
	PortFeature10MbitFullDuplex

	// PortFeature100MbitHalfDuplex is set when 100Mbit half-duplex rate
	// is supported.
	PortFeature100MbitHalfDuplex

	// PortFeature100MbitFullDuplex is set when 100Mbit full-duplex rate
	// is supported.
	PortFeature100MbitFullDuplex

	// PortFeature1GbitHalfDuplex is set when 1Gbit half-duplex rate is
	// supported.
	PortFeature1GbitHalfDuplex

	// PortFeature1GbitFullDuplex is set when 1Gib full-duplex rate is
	// supported.
	PortFeature1GbitFullDuplex

	// PortFeature10GbitFullDuplex is set when 10Gbit full-duplex rate
	// is supported.
	PortFeature10GbitFullDuplex

	// PortFeature40GbitFullDuplex is set when 40Gbit full-duplex rate
	// is supported.
	PortFeature40GbitFullDuplex

	// PortFeature100GbitFullDuplex is set when 100Gbit full-duplex rate
	// is supported.
	PortFeature100GbitFullDuplex

	// PortFeature1TbitFullDuplex is set when 1Tbit dull-duplex rate is
	// supported.
	PortFeature1TbitFullDuplex

	// PortFeatureOther is set when the rate is other than in this list.
	PortFeatureOther

	// PortFeatureCopper is set when the medium is copper.
	PortFeatureCopper

	// PortFeatureFiber is set when the medium is fiber.
	PortFeatureFiber

	// PortFeatureAutoneg is set when the port supports autonegotiation
	// of common transition parameters, such a speed, duplex mode, and
	// flow control.
	PortFeatureAutoneg

	// PortFeaturePause is set when the port supports the pause frame
	// used to halt the transmission of the sender for specific period
	// of time.
	PortFeaturePause

	// PortFeaturePauseAsym is set when the port supports flow control
	// through the asymmetric pause frames.
	PortFeaturePauseAsym
)

var portFeaturesText = []struct {
	mask PortFeature
	text string
}{
	{PortFeature10MbitHalfDuplex, "10 Mbps half-duplex"},
	{PortFeature10MbitFullDuplex, "10 Mbps full-duplex"},
	{PortFeature100MbitHalfDuplex, "100 Mbps half-duplex"},
	{PortFeature100MbitFullDuplex, "100 Mbps full-duplex"},
	{PortFeature1GbitHalfDuplex, "1 Gbps half-duplex"},
	{PortFeature1GbitFullDuplex, "1 Gbps full-duplex"},
	{PortFeature10GbitFullDuplex, "10 Gbps full-duplex"},
	{PortFeature40GbitFullDuplex, "40 Gbps full-duplex"},
	{PortFeature100GbitFullDuplex, "100 Gbps full-duplex"},
	{PortFeature1TbitFullDuplex, "1 Tbps full-duplex"},
	{PortFeatureOther, "other"},
	{PortFeatureCopper, "copper"},
	{PortFeatureFiber, "fiber"},
	{PortFeatureAutoneg, "autoneg"},
	{PortFeaturePause, "pause"},
	{PortFeaturePauseAsym, "pause asym"},
}

// String returns the human-redable representation of the port features.
func (f PortFeature) String() string {
	var features []string

	// Iterate though all of the port features and check if the
	// respective bit is set.
	for _, feature := range portFeaturesText {
		if feature.mask&f != 0 {
			features = append(features, feature.text)
		}
	}

	// Return space-joined list of the port features.
	return strings.Join(features, " ")
}

// PortConfig is a flag to indicate behavior of the physical port.
//
// These flags are used in Port structure to describe the current
// configuration. They are used in the PortMod message to configure
// the port's behavior.
type PortConfig uint32

const (
	// PortConfigDown is set when the port is administratively down.
	PortConfigDown PortConfig = 1 << iota

	// PortConfigNoSTP is set when 802.1D spanning tree is disabled
	// on this port.
	PortConfigNoSTP

	// PortConfigNoRcv is set when port drops all received packets.
	PortConfigNoRcv

	// PortConfigNoRcvSTP is set when it drops received 802.1D packets.
	PortConfigNoRcvSTP

	// PortConfigNoFlood is set when it is not used to flood packets.
	PortConfigNoFlood

	// PortConfigNoFwd is set when port drops packets forwarded to it.
	PortConfigNoFwd

	// PortConfigNoPacketIn is set when it is not allowed to sent
	// packet-in messages to this port.
	PortConfigNoPacketIn
)

var portConfigText = []struct {
	mask PortConfig
	text string
}{
	{PortConfigDown, "down"},
	{PortConfigNoSTP, "no STP"},
	{PortConfigNoRcv, "no recv"},
	{PortConfigNoRcvSTP, "no recv STP"},
	{PortConfigNoFlood, "no flood"},
	{PortConfigNoFwd, "no fwd"},
	{PortConfigNoPacketIn, "no packet in"},
}

// String returns a human representation of the port configuration.
func (c PortConfig) String() string {
	var configs []string

	// When the port is not DOWN, we will automatically assume
	// it is UP and append the respective text message.
	if PortConfigDown&c == 0 {
		configs = append(configs, "up")
	}

	for _, config := range portConfigText {
		if config.mask&c != 0 {
			configs = append(configs, config.text)
		}
	}

	return strings.Join(configs, " ")
}

// PortState defines the current state of the physical port. These are
// not configurable from the controller.
type PortState uint32

const (
	// PortStateLinkDown bit indicates that the physical link is not
	// present.
	PortStateLinkDown PortState = 1 << iota

	// PortStateBlocked bit indicates that a switch protocol outside of
	// OpenFlow, such as 802.1D Spanning Tree, is preventing the use of
	// that port with PortFlood.
	PortStateBlocked

	// PortStateLive indicates that port available for live for Fast
	// Failover Group.
	PortStateLive
)

var portStateText = map[PortState]string{
	PortStateLinkDown: "link down",
	PortStateBlocked:  "blocked",
	PortStateLive:     "live",
}

// PortState returns a human-readable representation of the port state.
func (s PortState) String() string {
	if text, ok := portStateText[s]; ok {
		return text
	}

	return "link up"
}

// PortNo defines a switch port number.
type PortNo uint32

const (
	// PortIn used to forward the packet out the input port. This
	// reserved port must be explicitly used in order to send back out
	// of the input port.
	PortIn PortNo = 0xfffffff8 + iota

	// PortTable used to submit the packet to the first flow table. This
	// destination port can only be used in packet-out messages.
	PortTable PortNo = 0xfffffff8 + iota

	// PortNormal used to process packets with normal L2/L3 switching.
	PortNormal PortNo = 0xfffffff8 + iota

	// PortFlood used to forward packets to all physical ports in VLAN,
	// except input port and those blocked or link down.
	PortFlood PortNo = 0xfffffff8 + iota

	// PortAll used to forward all physical ports except input port.
	PortAll PortNo = 0xfffffff8 + iota

	// PortController used to send the received packet to controller.
	PortController PortNo = 0xfffffff8 + iota

	// PortLocal is a local OpenFlow port.
	PortLocal PortNo = 0xfffffff8 + iota

	// PortAny is a wildcard port used only for flow mod (delete) and
	// flow stats requests. Selects all flows regardless of output port
	// (including flows with no output port).
	PortAny PortNo = 0xffffffff

	// PortMax is a Maximum number of physical and logical switch ports.
	PortMax PortNo = 0xffffff00
)

// portNameLen defines a length of the port name.
const portNameLen = 16

// Port defines the description of the switch port. This structure is
// returned within a body of the multipart request, used to get a
// description of all the ports in the system that support OpenFlow.
//
// For example, to retrieve the description of port on the system, the
// following request can be created:
//
//	body := ofp.NewMultipartRequest(ofp.MultipartTypePortDescription, nil)
//	req := of.NewRequest(of.TypeMultipartRequest, body)
type Port struct {
	// PortNo uniquely identifies a port within a switch.
	PortNo PortNo

	// HWAddr typically is the MAC address for the port.
	HWAddr net.HardwareAddr

	// Name is a human-readable name of the interface.
	Name string

	// Config describes the administrative settings.
	Config PortConfig

	// State describes the port internal state.
	State PortState

	Curr       PortFeature // Current features of the port.
	Advertised PortFeature // Features being advertised by the port.
	Supported  PortFeature // Features supported by the port.
	Peer       PortFeature // Features advertised by peer.

	// Current port bitrate in Kbps.
	CurrSpeed uint32

	// Max port bitrate in Kbps.
	MaxSpeed uint32
}

// WriteTo implements io.WriterTo interface. It serializes the port
// description into the wire format.
func (p *Port) WriteTo(w io.Writer) (int64, error) {
	name := make([]byte, portNameLen)
	copy(name, []byte(p.Name))

	return encoding.WriteTo(w, p.PortNo, pad4{}, p.HWAddr, pad2{},
		name, p.Config, p.State, p.Curr, p.Advertised, p.Supported,
		p.Peer, p.CurrSpeed, p.MaxSpeed,
	)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// port description from the wire format.
func (p *Port) ReadFrom(r io.Reader) (int64, error) {
	p.HWAddr = make(net.HardwareAddr, 6)
	var name [portNameLen]byte

	n, err := encoding.ReadFrom(r, &p.PortNo, &defaultPad4, &p.HWAddr,
		&defaultPad2, &name, &p.Config, &p.State, &p.Curr, &p.Advertised,
		&p.Supported, &p.Peer, &p.CurrSpeed, &p.MaxSpeed,
	)

	p.Name = string(name[:])
	return n, err
}

// Ports type groups the set of port descriptions.
type Ports []Port

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// list of ports from the wire format.
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

// PortMod is a message used to modify the behavior of the port.
//
// For example, to disable 802.1D spanning tree protocol on the
// second port, the following request can be created:
//
//	// Parse the hardware address of the port.
//	hwaddr, _ := net.ParseMAC("0123.4567.89ab")
//
//	pmod := ofp.PortMod{
//		PortNo: 2,
//		HWAddr: hwaddr,
//		Config: ofp.PortConfigNoSTP,
//		Mask:   ofp.PortConfigNoSTP,
//	}
//
//	req := of.NewRequest(TypePortMod, &pmod)
type PortMod struct {
	// PortNo uniquely identifies a port within a switch.
	PortNo PortNo

	// HWAddr is the hardware address. It is not configurable. This is
	// used to sanity-check the request, so it must be the same as
	// returned in an port description message.
	HWAddr net.HardwareAddr

	// Config is a bitmap of port configurations.
	Config PortConfig

	// Mask is a bitmap of flags to be changed.
	Mask PortConfig

	// Advertise is a bitmap of port features. Zero all bits to prevent
	// any action taking place.
	Advertise PortFeature
}

// WriteTo implements io.WriterTo interface. It serializes the port
// modification request into the wire format.
func (p *PortMod) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.PortNo, pad4{}, p.HWAddr, pad2{},
		p.Config, p.Mask, p.Advertise, pad4{},
	)
}

// ReadFrom implements io.ReaderFrom interfaace. It deserializes the
// port modification request from the wire format.
func (p *PortMod) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.PortNo, &defaultPad4, &p.HWAddr,
		&defaultPad2, &p.Config, &p.Mask, &p.Advertise, &defaultPad4,
	)
}

// PortReason defines types of changes being made about physical port.
type PortReason uint8

const (
	// PortReasonAdd indicates that port was added.
	PortReasonAdd PortReason = iota

	// PortReasonDelete indicates that port was removed.
	PortReasonDelete

	// PortReasonModify indicates that some attribute of the port has
	// changed.
	PortReasonModify
)

// PortStatus is the message used by the switch to inform the controller
// about the port being added, modified or removed.
type PortStatus struct {
	// Reason is the reason of this message.
	Reason PortReason

	// Port is a description of the changed port.
	Port Port
}

// WriteTo implements io.WriterTo interface. It serializes the
// port status into the wire format.
func (p *PortStatus) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.Reason, pad7{}, &p.Port)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// port status from the wire format.
func (p *PortStatus) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.Reason, &defaultPad7, &p.Port)
}

// PortStatsRequest is a multipart request message body used to retrieve
// the information about port statistics.
//
// For example, to request statistics for all ports, the following message
// can be created:
//
//	body := ofp.PortStatsRequest{ofp.PortAny}
//	req := of.NewRequest(ofp.TypeMultipartRequest,
//		ofp.NewMultipartRequest(ofp.TypePortStats, &body))
type PortStatsRequest struct {
	// PortNo uniquely identifies a port within a switch.
	//
	// Message must request statistics either for a single port or
	// of all ports using PortAny port number.
	PortNo PortNo
}

// WriteTo implements io.WriterTo interface. It serializes the port
// status request into the wire format.
func (p *PortStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.PortNo, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes
// the port stats request from the wire format.
func (p *PortStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.PortNo, &defaultPad4)
}

// PortStats defines the statistics. This is a body of the multipart
// reply returned on multipart port statistics request.
type PortStats struct {
	// PortNo uniquely identifies a port within a switch.
	PortNo PortNo

	// RxPackets is a number of received packets.
	RxPackets uint64

	// TxPackets is a number of transmitted packets.
	TxPackets uint64

	// RxBytes is a number of received bytes.
	RxBytes uint64

	// TxBytes is a number of transmitted bytes.
	TxBytes uint64

	// RxDropped is a number of packets dropped by RX.
	RxDropped uint64

	// TxDropped is a number of packets dropped by TX.
	TxDropped uint64

	// RxErrors is a number of receive errors. This is a super-set of
	// more specific receive errors and should be greater than or equal
	// to the sum of all Rx*Err values.
	RxErrors uint64

	// TxErrors is a number of transmit errors.
	TxErrors uint64

	// RxFrameErr is a number of frame alignment errors.
	RxFrameErr uint64

	// RxOverErr is a number of packets with RX overrun.
	RxOverErr uint64

	// RxCrcErr is a number of CRC errors.
	RxCrcErr uint64

	// Collisions is a number of collisions.
	Collisions uint64

	// DurationSec is a time port has been alive in seconds.
	DurationSec uint32

	// DurationNSec is a time port has been alive in nanoseconds beyond
	// DurationSec.
	DurationNSec uint32
}

// WriteTo implements io.WriterTo interface. It serializes the port stats
// into the wire format.
func (p *PortStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, p.PortNo, pad4{}, p.RxPackets,
		p.TxPackets, p.RxBytes, p.TxBytes, p.RxDropped, p.TxDropped,
		p.RxErrors, p.TxErrors, p.RxFrameErr, p.RxOverErr, p.RxCrcErr,
		p.Collisions, p.DurationSec, p.DurationNSec,
	)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the port
// stats from the wire format.
func (p *PortStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &p.PortNo, &defaultPad4, &p.RxPackets,
		&p.TxPackets, &p.RxBytes, &p.TxBytes, &p.RxDropped, &p.TxDropped,
		&p.RxErrors, &p.TxErrors, &p.RxFrameErr, &p.RxOverErr, &p.RxCrcErr,
		&p.Collisions, &p.DurationSec, &p.DurationNSec,
	)
}
