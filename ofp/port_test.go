package ofp

import (
	"net"
	"testing"

	"github.com/netrack/openflow/internal/encodingtest"
)

var hwaddr, _ = net.ParseMAC("0123.4567.89ab")

func TestPortFeatureString(t *testing.T) {
	pf := []PortFeature{
		PortFeature1GbitHalfDuplex | PortFeatureCopper,
		PortFeatureOther | PortFeaturePauseAsym,
		PortFeature10MbitHalfDuplex | PortFeatureFiber,
		PortFeature40GbitFullDuplex | PortFeaturePause,
		PortFeature1TbitFullDuplex | PortFeatureFiber,
	}

	features := map[PortFeature]string{
		pf[0]: "1 Gbps half-duplex copper",
		pf[1]: "other pause asym",
		pf[2]: "10 Mbps half-duplex fiber",
		pf[3]: "40 Gbps full-duplex pause",
		pf[4]: "1 Tbps full-duplex fiber",
	}

	for feature, text := range features {
		if feature.String() != text {
			t.Errorf("Invalid port feature, expected:\n"+
				"`%s` got:\n`%s`", text, feature.String())
		}
	}
}

func TestPortConfigString(t *testing.T) {
	pc := map[PortConfig]string{
		PortConfigDown: "down",
		PortConfig(0):  "up",
	}

	for config, text := range pc {
		if config.String() != text {
			t.Errorf("Invalid port config, expected:\n, "+
				"`%s` got:\n`%s`", text, config.String())
		}
	}
}

func TestPortStateString(t *testing.T) {
	ps := map[PortState]string{
		PortStateLinkDown: "link down",
		PortStateLive:     "live",
		PortState(0):      "link up",
	}

	for state, text := range ps {
		if state.String() != text {
			t.Errorf("Invalid port state, expected:\n"+
				"`%s` got:\n`%s`", text, state.String())
		}
	}
}

func TestPort(t *testing.T) {
	features := PortFeature1GbitFullDuplex | PortFeatureFiber
	peer := PortFeature10GbitFullDuplex | PortFeatureCopper

	name := make([]byte, portNameLen)
	copy(name, "sw1-eth0")

	tests := []encodingtest.MU{
		{ReadWriter: &Port{
			PortNo:     PortNormal,
			HWAddr:     hwaddr,
			Name:       string(name),
			Config:     PortConfigDown,
			State:      PortStateLinkDown,
			Curr:       features,
			Advertised: features,
			Supported:  features,
			Peer:       peer,
			CurrSpeed:  42,
			MaxSpeed:   43,
		}, Bytes: []byte{
			0xff, 0xff, 0xff, 0xfa, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, // Hardware address.
			0x00, 0x00, // 2-byte padding.
			0x73, 0x77, 0x31, 0x2d, 0x65, 0x74, 0x68, 0x30,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Name.
			0x00, 0x00, 0x00, 0x01, // Port configuration.
			0x00, 0x00, 0x00, 0x01, // Port state.
			0x00, 0x00, 0x10, 0x20, // Current port features.
			0x00, 0x00, 0x10, 0x20, // Advertised port features.
			0x00, 0x00, 0x10, 0x20, // Supported port features.
			0x00, 0x00, 0x08, 0x40, // Peer port features.
			0x00, 0x00, 0x00, 0x2a, // Current speed.
			0x00, 0x00, 0x00, 0x2b, // Maximum speed.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestPortMod(t *testing.T) {
	features := PortFeature10GbitFullDuplex | PortFeatureAutoneg

	tests := []encodingtest.MU{
		{ReadWriter: &PortMod{
			PortNo:    PortFlood,
			HWAddr:    hwaddr,
			Config:    PortConfigNoFwd,
			Mask:      PortConfig(0),
			Advertise: features,
		}, Bytes: []byte{
			0xff, 0xff, 0xff, 0xfb, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, // Hardware address.
			0x00, 0x00, // 2-byte padding.
			0x00, 0x00, 0x00, 0x20, // Port configuration.
			0x00, 0x00, 0x00, 0x00, // Port configuration mask.
			0x00, 0x00, 0x20, 0x40, // Supported port features.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestPortStatus(t *testing.T) {
	features := PortFeature40GbitFullDuplex | PortFeatureFiber
	peer := PortFeature10MbitHalfDuplex | PortFeatureCopper

	name := make([]byte, portNameLen)
	copy(name, "sw1-eth1")

	tests := []encodingtest.MU{
		{ReadWriter: &PortStatus{
			Reason: PortReasonAdd,
			Port: Port{
				PortNo:     PortFlood,
				HWAddr:     hwaddr,
				Name:       string(name),
				Config:     PortConfigDown,
				State:      PortStateLinkDown,
				Curr:       features,
				Advertised: features,
				Supported:  features,
				Peer:       peer,
				CurrSpeed:  2047,
				MaxSpeed:   65535,
			},
		}, Bytes: []byte{
			0x00,                                     // Port status reason.
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 7-byte padding.

			// Port.
			0xff, 0xff, 0xff, 0xfb, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, // Hardware address.
			0x00, 0x00, // 2-byte padding.
			0x73, 0x77, 0x31, 0x2d, 0x65, 0x74, 0x68, 0x31,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Name.
			0x00, 0x00, 0x00, 0x01, // Port configuration.
			0x00, 0x00, 0x00, 0x01, // Port state.
			0x00, 0x00, 0x10, 0x80, // Current port features.
			0x00, 0x00, 0x10, 0x80, // Advertised port features.
			0x00, 0x00, 0x10, 0x80, // Supported port features.
			0x00, 0x00, 0x08, 0x01, // Peer port features.
			0x00, 0x00, 0x07, 0xff, // Current speed.
			0x00, 0x00, 0xff, 0xff, // Maximum speed.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestPortStatsRequest(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &PortStatsRequest{
			PortNo: PortNo(2),
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x02, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
		}},
	}

	encodingtest.RunMU(t, tests)
}

func TestPortStats(t *testing.T) {
	tests := []encodingtest.MU{
		{ReadWriter: &PortStats{
			PortNo:       PortNo(3),
			RxPackets:    6773009508081008653,
			TxPackets:    4449515516159517871,
			RxBytes:      5376522659808774614,
			TxBytes:      13087997567122713610,
			RxDropped:    14817532963293049516,
			TxDropped:    1483028468767136997,
			RxErrors:     8919350792100524585,
			TxErrors:     8235451563360597435,
			RxFrameErr:   12151775103618862338,
			RxOverErr:    8345656100423014979,
			RxCrcErr:     7312308705513999709,
			Collisions:   7394863223143382373,
			DurationSec:  684498643,
			DurationNSec: 2789499113,
		}, Bytes: []byte{
			0x00, 0x00, 0x00, 0x03, // Port number.
			0x00, 0x00, 0x00, 0x00, // 4-byte padding.
			0x5d, 0xfe, 0x90, 0x43, 0x3d, 0x83, 0x00, 0x0d, // Rx packets.
			0x3d, 0xbf, 0xda, 0xc9, 0x93, 0x46, 0xe4, 0xaf, // Tx packets.
			0x4a, 0x9d, 0x3e, 0xdf, 0x80, 0xbd, 0x89, 0xd6, // Rx bytes.
			0xb5, 0xa1, 0xe8, 0x71, 0xb1, 0x86, 0xac, 0x0a, // Tx bytes.
			0xcd, 0xa2, 0x73, 0xb9, 0x34, 0xb9, 0x32, 0xac, // Rx dropped.
			0x14, 0x94, 0xc6, 0x88, 0xf0, 0xa8, 0x1c, 0xe5, // Tx dropped.
			0x7b, 0xc7, 0xe6, 0x49, 0xe6, 0x40, 0x32, 0x29, // Rx errors.
			0x72, 0x4a, 0x33, 0x90, 0x47, 0x0b, 0x35, 0xbb, // Tx errors.
			0xa8, 0xa3, 0xc7, 0x12, 0xe9, 0xa3, 0x4d, 0x02, // Rx frame errors.
			0x73, 0xd1, 0xba, 0x01, 0x93, 0x49, 0x1a, 0x43, // Rx over errors.
			0x65, 0x7a, 0x8a, 0x06, 0x80, 0x28, 0x8d, 0x5d, // Rx CRC errors.
			0x66, 0x9f, 0xd4, 0xeb, 0xfa, 0x0f, 0xbd, 0x65, // Collisions.
			0x28, 0xcc, 0x9e, 0xd3, // Duration seconds.
			0xa6, 0x44, 0x60, 0xe9, // Duration nanoseconds.
		}},
	}

	encodingtest.RunMU(t, tests)
}
