package ofputil

import (
	"bytes"
	"fmt"

	"github.com/netrack/openflow/internal/encoding"
	"github.com/netrack/openflow/ofp"
)

func bytesOf(v interface{}) []byte {
	var buf bytes.Buffer

	_, err := encoding.WriteTo(&buf, v)
	if err != nil {
		text := "ofputil: unable to marshal %v"
		panic(fmt.Errorf(text, err))
	}

	return buf.Bytes()
}

func ExtendedMatch(xms ...ofp.XM) ofp.Match {
	return ofp.Match{ofp.MatchTypeXM, xms}
}

// basic creates an Openflow basic extensible match of the given type.
func basic(t ofp.XMType, val ofp.XMValue, mask ofp.XMValue) ofp.XM {
	return ofp.XM{
		Class: ofp.XMClassOpenflowBasic,
		Type:  t, Value: val, Mask: mask,
	}
}

// MatchEthType creates an Openflow basic extensible match of Ethernet
// payload type.
func MatchEthType(eth uint16) ofp.XM {
	return basic(ofp.XMTypeEthType, bytesOf(eth), nil)
}

// MatchInPort creates an Openflow basic extensible match of in port.
func MatchInPort(port ofp.PortNo) ofp.XM {
	return basic(ofp.XMTypeInPort, bytesOf(port), nil)
}

// MatchIPProto creates an Openflow basic extensible match of IP protocol
// payload type.
func MatchIPProto(ipp uint8) ofp.XM {
	return basic(ofp.XMTypeIPProto, bytesOf(ipp), nil)
}

// MatchICMPv6Type creates an Openflow basic extensible match of ICMPv6
// message type.
func MatchICMPv6Type(icmpt uint8) ofp.XM {
	return basic(ofp.XMTypeICMPv6Type, bytesOf(icmpt), nil)
}

// MatchIPv6ExtHeader creates an Openflow basic extensible match of IPv6
// extension header.
func MatchIPv6ExtHeader(header uint16) ofp.XM {
	return basic(ofp.XMTypeIPv6ExtHeader, bytesOf(header), nil)
}
