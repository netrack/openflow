package of

import (
	"bytes"
	"io"
)

type readerFunc func([]byte) (int, error)

func (fn readerFunc) Read(b []byte) (int, error) {
	return fn(b)
}

func newReader(w ...io.WriterTo) io.Reader {
	return readerFunc(func(b []byte) (int, error) {
		var buf bytes.Buffer
		var err error

		for _, wt := range w {
			if wt == nil {
				continue
			}

			if _, err = wt.WriteTo(&buf); err != nil {
				break
			}
		}

		if err != nil {
			return 0, err
		}

		return buf.Read(b)
	})
}
