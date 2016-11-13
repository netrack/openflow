package of

import (
	"bytes"
	"io"
)

type writerToFunc func(io.Writer) (int64, error)

func (fn writerToFunc) WriteTo(w io.Writer) (int64, error) {
	return fn(w)
}

type readerFunc func([]byte) (int, error)

func (fn readerFunc) Read(b []byte) (int, error) {
	return fn(b)
}

func newReader(w ...io.WriterTo) io.Reader {
	fn := func(b []byte) (int, error) {
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
	}

	return readerFunc(fn)
}

func MultiWriterTo(w ...io.WriterTo) io.WriterTo {
	fn := func(wr io.Writer) (int64, error) {
		var n int64

		for _, writer := range w {
			nn, err := writer.WriteTo(wr)
			n += nn

			if err != nil {
				return n, err
			}
		}

		return n, nil
	}

	return writerToFunc(fn)
}
