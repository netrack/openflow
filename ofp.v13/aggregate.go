package ofp

import (
	"io"

	"github.com/netrack/openflow/encoding"
)

// Aggregate information about multiple flow entries is requested.
type AggregateStatsRequest struct {
	TableID Table

	OutPort  PortNo
	OutGroup Group

	Cookie     uint64
	CookieMask uint64

	Match Match
}

func (a *AggregateStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		a.TableID,
		pad3{},
		a.OutPort,
		a.OutGroup,
		pad4{},
		a.Cookie,
		a.CookieMask,
		&a.Match)
}

type AggregateStatsReply struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
}

func (a *AggregateStatsReply) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(
		w, a.PacketCount, a.ByteCount, a.FlowCount, pad4{})
}
