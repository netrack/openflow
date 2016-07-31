package binary

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ByteOrder binary.ByteOrder

var (
	BigEndian    ByteOrder = binary.BigEndian
	LittleEndian ByteOrder = binary.LittleEndian
)

type reader struct {
	io.Reader
	read int64
}

func (r *reader) Read(b []byte) (int, error) {
	n, err := r.Reader.Read(b)
	r.read += int64(n)
	return n, err
}

func Read(r io.Reader, order ByteOrder, data interface{}) (
	n int64, err error) {

	rd := &reader{r, 0}
	err = binary.Read(rd, order, data)
	return rd.read, err
}

func Write(w io.Writer, order ByteOrder, data interface{}) (
	n int64, err error) {

	var wbuf bytes.Buffer

	err = binary.Write(&wbuf, order, data)
	if err != nil {
		return
	}

	return wbuf.WriteTo(w)
}

func ReadSlice(r io.Reader, order ByteOrder, slice []interface{}) (
	int64, error) {

	rd := &reader{r, 0}

	for _, elem := range slice {
		err := binary.Read(rd, order, elem)
		if err != nil {
			return rd.read, err
		}
	}

	return rd.read, nil
}

func WriteSlice(w io.Writer, order ByteOrder, slice []interface{}) (
	n int64, err error) {

	var wbuf bytes.Buffer

	for _, elem := range slice {
		err = binary.Write(&wbuf, order, elem)
		if err != nil {
			return
		}
	}

	return wbuf.WriteTo(w)
}
