package of

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

func TestTypeDispatcher(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	rd := NewTypeDispatcher()
	rd.HandleFunc(TypeHello, func(rw ResponseWriter, r *Request) {
		rw.Write([]byte{0, 0, 0, 0})
		rw.WriteHeader()
		wg.Done()
	})

	reader := bytes.NewBuffer([]byte{4, 0, 0, 8, 0, 0, 0, 0})
	conn := &dummyConn{r: *reader}

	s := Server{Addr: "0.0.0.0:6633", Handler: rd}
	err := s.Serve(&dummyListener{conn})

	if err != io.EOF {
		t.Fatal("Serve failed:", err)
	}

	wg.Wait()

	returned := fmt.Sprintf("%x", conn.w.Bytes())
	if returned != "0400000c0000000000000000" {
		t.Fatal("Invalid data returned: ", returned)
	}
}
