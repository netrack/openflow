package ofp

import (
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// AggregateStatsRequest is a multipart request used to aggregate
// information about multiple flow entries.
//
// This message is used as a body of the multipart request. For example,
// to retrieve information about the flow entries matching the first
// ingress port, the following multipart request could be sent:
//
//	body := &ofp.AggregateStatsRequest{
//		Table:    ofp.TableAll,
//		OutPort:  ofp.PortAny,
//		OutGroup: ofp.GroupAny,
//		Match: ofputil.ExtendedMatch(
//			ofputil.MatchInPort(1),
//		),
//	}
//
//	req := ofp.NewMultipartRequest(
//		ofp.MutipartTypeAggregate, body)
//	...
type AggregateStatsRequest struct {
	// Table is an identifier of the table to read or TableAll
	// to read all tables.
	Table Table

	// Require matching entries to include this as an output port.
	// A value PortAny indicates no restrictions.
	OutPort PortNo

	// Require matching entries to include this as an output group.
	// A value GroupAny indicates no restrictions.
	OutGroup Group

	// Require matching entries to contain this cookie value.
	Cookie uint64

	// Mask used to restrict the cookie bits that must match. A zero
	// value indicates no restrictions.
	CookieMask uint64

	// Fields to match.
	Match Match
}

// Cookies implements of.CookieJar interface and returns the cookies of
// the aggregate statistics request
func (a *AggregateStatsRequest) Cookies() uint64 {
	return a.Cookie
}

// SetCookies implements of.CookieJar interface and sets the specified
// cookies to the aggregate statistics request.
func (a *AggregateStatsRequest) SetCookies(cookies uint64) {
	a.Cookie = cookies
}

// WriteTo implements io.WriterTo interface. It serializes the
// aggregated statistics request into the wire format with required
// padding.
func (a *AggregateStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, a.Table, pad3{},
		a.OutPort, a.OutGroup, pad4{}, a.Cookie,
		a.CookieMask, &a.Match)
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// aggregated statistics request from the wire format.
func (a *AggregateStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &a.Table, &defaultPad3,
		&a.OutPort, &a.OutGroup, &defaultPad4, &a.Cookie,
		&a.CookieMask, &a.Match)
}

// AggregateStats is a response on the aggregated statistics request.
// This message is used as a body of multipart reply.
type AggregateStats struct {
	// PacketCount is a number of packets in flow.
	PacketCount uint64

	// ByteCount is a number of bytes in flow.
	ByteCount uint64

	// FlowCount in a number of flows.
	FlowCount uint32
}

// WriteTo implements io.WriterTo interface. It serializes the
// aggregated statistics into the wire format with required
// padding.
func (a *AggregateStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(
		w, a.PacketCount, a.ByteCount, a.FlowCount, pad4{})
}

// ReadFrom implements io.ReaderFrom interface. It deserializes the
// aggregated statistics from the wire format.
func (a *AggregateStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(
		r, &a.PacketCount, &a.ByteCount, &a.FlowCount, &defaultPad4)
}
