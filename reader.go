package of

import (
	"bytes"
	"io"
	"sync"
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
	var buf bytes.Buffer
	var mulErr error

	var once sync.Once
	mulW := MultiWriterTo(w...)

	fn := func(b []byte) (int, error) {
		once.Do(func() {
			_, mulErr = mulW.WriteTo(&buf)
		})

		if mulErr != nil {
			return 0, mulErr
		}

		return buf.Read(b)
	}

	return readerFunc(fn)
}

func MultiWriterTo(w ...io.WriterTo) io.WriterTo {
	fn := func(wr io.Writer) (int64, error) {
		var n int64

		for _, writer := range w {
			if writer == nil {
				continue
			}

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
