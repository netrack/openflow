package encoding

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

// reader type used to calculate the count of bytes retrieved from the
// configured reader instance.
type reader struct {
	io.Reader
	read int64
}

// ReadWriter describes typed that are capable both, to write their
// representation into the writer and read it from the reader.
type ReadWriter interface {
	io.ReaderFrom
	io.WriterTo
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

type ReaderMaker interface {
	MakeReader() io.ReaderFrom
}

type readerMakerFunc func() io.ReaderFrom

func (fn readerMakerFunc) MakeReader() io.ReaderFrom {
	return fn()
}

func ReaderMakerOf(v interface{}) ReaderMaker {
	valueType := reflect.TypeOf(v)
	return readerMakerFunc(func() io.ReaderFrom {
		value := reflect.New(valueType)
		return value.Interface().(io.ReaderFrom)
	})
}

// SkipEndOfFile returns nil of the given error caused by the
// end of file.
func SkipEndOfFile(err error) error {
	if err == io.EOF {
		return nil
	}

	return err
}
