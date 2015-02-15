package binary

import (
	"encoding/binary"
	"io"
)

type ByteOrder binary.ByteOrder

var (
	BigEndian    ByteOrder = binary.BigEndian
	LittleEndian ByteOrder = binary.LittleEndian
)

func Read(r io.Reader, order ByteOrder, data interface{}) error {
	return binary.Read(r, order, data)
}

func Write(w io.Writer, order ByteOrder, data interface{}) error {
	return binary.Write(w, order, data)
}

func ReadSlice(r io.Reader, order ByteOrder, slice []interface{}) error {
	for _, elem := range slice {
		err := binary.Read(r, order, elem)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteSlice(w io.Writer, order ByteOrder, slice []interface{}) error {
	for _, elem := range slice {
		err := binary.Write(w, order, elem)
		if err != nil {
			return err
		}
	}

	return nil
}
