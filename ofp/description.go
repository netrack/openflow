package ofp

import (
	"io"

	"github.com/netrack/openflow/internal/encoding"
)

const (
	// descLen is a length of the maximum length of the description
	// attributes (e.g. manufacturer, hardware, software).
	descLen = 256

	// serialNumLen is a length of the switch serial number.
	serialNumLen = 32
)

// Description is an information about the switch manufacturer, hardware
// revision, software revision, serial number, and a description field.
//
// This message is used as a response body of multipart request. For
// example, to retrieve detailed description about the switch, the
// following request could be sent:
//
//	req := ofp.NewMultipartRequest(
//		ofp.MultipartTypeDescription, nil)
//	...
type Description struct {
	// Manufacturer description.
	Manufacturer string

	// Hardware description.
	Hardware string

	// Software description.
	Software string

	// Serial number.
	SerialNum string

	// Human readable description of datapath.
	Datapath string
}

// copyTo copies the specified string into the slice of bytes of given
// length.
func (d *Description) copyTo(s string, length int) []byte {
	b := make([]byte, length)
	copy(b, []byte(s))
	return b
}

// WriteTo implements io.WriterTo interface. It serializes the switch
// description into the wire format.
func (d *Description) WriteTo(w io.Writer) (int64, error) {
	mfr := d.copyTo(d.Manufacturer, descLen)
	hw := d.copyTo(d.Hardware, descLen)
	sw := d.copyTo(d.Software, descLen)
	sn := d.copyTo(d.SerialNum, serialNumLen)
	dp := d.copyTo(d.Datapath, descLen)

	return encoding.WriteTo(w, mfr, hw, sw, sn, dp)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// switch description from the wire format.
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
