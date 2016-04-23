package ofputil

import (
	"io"

	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp.v13"
)

func PacketInCookieReader() of.CookieReader {
	return of.CookieReaderFunc(func(r io.Reader) (of.CookieJar, error) {
		var packetIn ofp.PacketIn

		if _, err := packetIn.ReadFrom(r); err != nil {
			return nil, err
		}

		return &packetIn, nil
	})
}
