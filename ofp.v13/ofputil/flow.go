package ofputil

import (
	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp.v13"
)

func TableFlush(table ofp.Table) *of.Request {
	r, _ := of.NewRequest(of.T_FLOW_MOD, of.NewReader(&ofp.FlowMod{
		TableID:  table,
		Command:  ofp.FC_DELETE,
		BufferID: ofp.NO_BUFFER,
		OutPort:  ofp.P_ANY,
		OutGroup: ofp.G_ANY,
		Match:    ofp.Match{ofp.MT_OXM, nil},
	}))

	return r
}

func FlowFlush(table ofp.Table, match ofp.Match) *of.Request {
	r, _ := of.NewRequest(of.T_FLOW_MOD, of.NewReader(&ofp.FlowMod{
		TableID:  table,
		Command:  ofp.FC_DELETE,
		BufferID: ofp.NO_BUFFER,
		OutPort:  ofp.P_ANY,
		OutGroup: ofp.G_ANY,
		Match:    match,
	}))

	return r
}

func FlowDrop(table ofp.Table) *of.Request {
	r, _ := of.NewRequest(of.T_FLOW_MOD, of.NewReader(&ofp.FlowMod{
		TableID:  table,
		Command:  ofp.FC_ADD,
		BufferID: ofp.NO_BUFFER,
		Match:    ofp.Match{ofp.MT_OXM, nil},
	}))

	return r
}

func EthDstAddr(hwaddr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ETH_DST, hwaddr, mask}
}

func EthSrcAddr(hwaddr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ETH_SRC, hwaddr, mask}
}

func EthType(proto uint16, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ETH_TYPE, of.Bytes(proto), mask}
}

func ARPOpType(op uint16, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ARP_OP, of.Bytes(op), mask}
}

func ARPTargetHWAddr(hwaddr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ARP_THA, hwaddr, mask}
}

func ARPTargetProtoAddr(addr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_ARP_TPA, addr, mask}
}

func IPv4DstAddr(ipaddr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_IPV4_DST, ipaddr, mask}
}

func IPv4SrcAddr(ipaddr []byte, mask []byte) ofp.OXM {
	return ofp.OXM{ofp.XMC_OPENFLOW_BASIC, ofp.XMT_OFB_IPV4_SRC, ipaddr, mask}
}
