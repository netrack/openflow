package encoding

import (
	"bufio"
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

// Read implements io.Reader interface.
func (r *reader) Read(b []byte) (int, error) {
	n, err := r.Reader.Read(b)
	r.read += int64(n)
	return n, err
}

// ReadWriter describes typed that are capable both, to write their
// representation into the writer and read it from the reader.
type ReadWriter interface {
	io.ReaderFrom
	io.WriterTo
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

func ReadSliceFrom(r io.Reader, rm ReaderMaker, slice interface{}) (int64, error) {
	var n int64
	sliceValue := reflect.ValueOf(slice)

	for {
		reader, err := rm.MakeReader()
		if err != nil {
			return n, err
		}

		nn, err := reader.ReadFrom(r)
		n += nn

		if err != nil {
			return n, SkipEOF(err)
		}

		elem := reflect.ValueOf(reader).Elem()
		reflect.Append(sliceValue, elem)
	}

	return n, nil
}

// ReaderMaker defines factory types, used to created new exemplars
// of the io.ReaderFrom interface.
type ReaderMaker interface {
	MakeReader() (io.ReaderFrom, error)
}

// ReaderMakerOf creates a new ReaderMaker based on the specified
// type. Pointer to the type must implement io.ReaderFrom interface.
func ReaderMakerOf(v interface{}) ReaderMaker {
	valueType := reflect.TypeOf(v)
	return ReaderMakerFunc(func() (io.ReaderFrom, error) {
		value := reflect.New(valueType)
		return value.Interface().(io.ReaderFrom), nil
	})
}

// ReaderMakerFunc is a function adapter for ReaderMaker interface.
type ReaderMakerFunc func() (io.ReaderFrom, error)

// MakeReader implements ReaderMaker interface.
func (fn ReaderMakerFunc) MakeReader() (io.ReaderFrom, error) {
	return fn()
}

func ScanFrom(r io.Reader, v interface{}, rm ReaderMaker) (int64, error) {
	// Retrieve the size of the instruction type, that preceeds
	// every instruction body.
	valType := reflect.TypeOf(v)
	typeLen := int(valType.Elem().Size())

	var n, num int64
	// To keep the implementation of the instruction unmarshaling
	// consistent with marshaling, we have to put the instruction
	// type back to the reader during unmarshaling of the list of
	// instructions.
	rdbuf := bufio.NewReader(r)

	for {
		typeBuf, err := rdbuf.Peek(typeLen)
		if err != nil {
			return n, SkipEOF(err)
		}

		// Unmarshal the instruction type from the peeked bytes.
		_, err = ReadFrom(bytes.NewReader(typeBuf), v)
		if err != nil {
			return n, err
		}

		// It is defined a corresponding io.ReaderFrom factory for
		// each value type, so we could parse the raw bytes using
		// the correct implementation of the reader.
		rdfrom, err := rm.MakeReader()
		if err != nil {
			return n, err
		}

		// Read the corresponding value from the binary representation.
		num, err = rdfrom.ReadFrom(rdbuf)
		n += num

		if err != nil {
			return n, SkipEOF(err)
		}
	}

	return n, nil
}

// SkipEOF returns nil of the given error caused by the
// end of file.
func SkipEOF(err error) error {
	if err == io.EOF {
		return nil
	}

	return err
}
