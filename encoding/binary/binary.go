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

func Read(r io.Reader, order ByteOrder, data interface{}) (n int64, err error) {
	var rbuf bytes.Buffer

	n, err = rbuf.ReadFrom(r)
	if err != nil {
		return
	}

	err = binary.Read(&rbuf, order, data)
	return
}

func Write(w io.Writer, order ByteOrder, data interface{}) (n int64, err error) {
	var wbuf bytes.Buffer

	err = binary.Write(&wbuf, order, data)
	if err != nil {
		return
	}

	return wbuf.WriteTo(w)
}

func ReadSlice(r io.Reader, order ByteOrder, slice []interface{}) (n int64, err error) {
	var rbuf bytes.Buffer

	n, err = rbuf.ReadFrom(r)
	if err != nil {
		return
	}

	for _, elem := range slice {
		err = binary.Read(&rbuf, order, elem)
		if err != nil {
			return
		}
	}

	return
}

func WriteSlice(w io.Writer, order ByteOrder, slice []interface{}) (n int64, err error) {
	var wbuf bytes.Buffer

	for _, elem := range slice {
		err = binary.Write(&wbuf, order, elem)
		if err != nil {
			return
		}
	}

	return wbuf.WriteTo(w)
}
