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
