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

func MatchInPort(port ofp.PortNo) ofp.XM {
	return ofp.XM{
		Class: ofp.XMClassOpenflowBasic,
		Type:  ofp.XMTypeInPort,
		Value: bytesOf(port),
		Mask:  nil,
	}
}
