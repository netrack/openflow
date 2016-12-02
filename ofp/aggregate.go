package ofp

import (
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

// Aggregate information about multiple flow entries is requested.
type AggregateStatsRequest struct {
	Table Table

	OutPort  PortNo
	OutGroup Group

	Cookie     uint64
	CookieMask uint64

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

func (a *AggregateStatsRequest) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w,
		a.Table,
		pad3{},
		a.OutPort,
		a.OutGroup,
		pad4{},
		a.Cookie,
		a.CookieMask,
		&a.Match)
}

func (a *AggregateStatsRequest) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r,
		&a.Table,
		&defaultPad3,
		&a.OutPort,
		&a.OutGroup,
		&defaultPad4,
		&a.Cookie,
		&a.CookieMask,
		&a.Match)
}

type AggregateStats struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
}

func (a *AggregateStats) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(
		w, a.PacketCount, a.ByteCount, a.FlowCount, pad4{})
}

func (a *AggregateStats) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(
		r, &a.PacketCount, &a.ByteCount, &a.FlowCount, &defaultPad4)
}
