package ofp

import (
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	descLen      = 256
	serialNumLen = 32
)

type Description struct {
	Manufacturer string
	Hardware     string
	Software     string
	SerialNum    string
	Datapath     string
}

func (d *Description) cpTo(s string, length int) []byte {
	b := make([]byte, length)
	copy(b, []byte(s))
	return b
}

func (d *Description) WriteTo(w io.Writer) (int64, error) {
	mfr := d.cpTo(d.Manufacturer, descLen)
	hw := d.cpTo(d.Hardware, descLen)
	sw := d.cpTo(d.Software, descLen)
	sn := d.cpTo(d.SerialNum, serialNumLen)
	dp := d.cpTo(d.Datapath, descLen)

	return encoding.WriteTo(w, mfr, hw, sw, sn, dp)
}

func (d *Description) ReadFrom(r io.Reader) (int64, error) {
	var mfr, hw, sw, dp [descLen]byte
	var sn [serialNumLen]byte

	n, err := encoding.ReadFrom(r, &mfr, &hw, &sw, &sn, &dp)
	if err != nil {
		return n, err
	}

	d.Manufacturer = string(mfr[:])
	d.Hardware = string(hw[:])
	d.Software = string(sw[:])
	d.SerialNum = string(sn[:])
	d.Datapath = string(dp[:])
	return n, err
}
