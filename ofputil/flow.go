package ofputil

import (
	"github.com/netrack/openflow/ofp"
	of "github.com/netrack/openflow"
)

func TableFlush(table ofp.Table) *of.Request {
	return of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowDelete,
		BufferID: ofp.NoBuffer,
		OutPort:  ofp.PortAny,
		OutGroup: ofp.GroupAny,
		Match:    ofp.Match{ofp.MatchTypeXM, nil},
	})
}

func FlowFlush(table ofp.Table, match ofp.Match) *of.Request {
	return of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowDelete,
		BufferID: ofp.NoBuffer,
		OutPort:  ofp.PortAny,
		OutGroup: ofp.GroupAny,
		Match:    match,
	})
}

func FlowDrop(table ofp.Table) *of.Request {
	return of.NewRequest(of.TypeFlowMod, &ofp.FlowMod{
		Table:    table,
		Command:  ofp.FlowAdd,
		BufferID: ofp.NoBuffer,
		Match:    ofp.Match{ofp.MatchTypeXM, nil},
	})
}
