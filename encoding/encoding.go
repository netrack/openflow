package encoding

import (
	"bytes"
	"encoding/binary"
	"io"
)

// reader type used to calculate the count of bytes retrieved from the
// configured reader instance.
type reader struct {
	io.Reader
	read int64
}

// Read implements io.Reader interface.
func (r *reader) Read(b []byte) (int, error) {
	n, err := r.Reader.Read(b)
	r.read += int64(n)
	return n, err
}

func WriteTo(w io.Writer, v ...interface{}) (int64, error) {
	var (
		wbuf bytes.Buffer
		err  error
	)

	for _, elem := range v {
		switch elem := elem.(type) {
		case io.WriterTo:
			_, err = elem.WriteTo(&wbuf)
		default:
			err = binary.Write(&wbuf, binary.BigEndian, elem)
		}

		if err != nil {
			return 0, err
		}
	}

	return wbuf.WriteTo(w)
}

func ReadFrom(r io.Reader, v ...interface{}) (int64, error) {
	var (
		num int64
		err error
	)

	rd := &reader{r, 0}

	for _, elem := range v {
		switch elem := elem.(type) {
		case io.ReaderFrom:
			num, err = elem.ReadFrom(r)
			rd.read += num
		default:
			err = binary.Read(rd, binary.BigEndian, elem)
		}

		if err != nil {
			return rd.read, err
		}
	}

	return rd.read, nil
}
