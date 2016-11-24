package ofputil

import (
	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp"
)

func TableFlush(table ofp.Table) *of.Request {
	r, _ := of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowDelete,
		BufferID: ofp.NoBuffer,
		OutPort:  ofp.PortAny,
		OutGroup: ofp.GroupAny,
		Match:    ofp.Match{ofp.MatchTypeXM, nil},
	})

	return r
}

func FlowFlush(table ofp.Table, match ofp.Match) *of.Request {
	r, _ := of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowDelete,
		BufferID: ofp.NoBuffer,
		OutPort:  ofp.PortAny,
		OutGroup: ofp.GroupAny,
		Match:    match,
	})

	return r
}

func FlowDrop(table ofp.Table) *of.Request {
	r, _ := of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowAdd,
		BufferID: ofp.NoBuffer,
		Match:    ofp.Match{ofp.MatchTypeXM, nil},
	})

	return r
}
