package ofptest_test

import (
	"fmt"

	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofptest"
)

func ExampleResponseRecorder() {
	handler := func(w of.ResponseWriter, r *of.Request) {
		w.Write(r.Header.Copy(), nil)
	}

	req := of.NewRequest(of.TypeHello, nil)
	w := ofptest.NewRecorder()

	handler(w, req)
	fmt.Printf("type: %d", w.First().Header.Type)
	// Output: type: 0
}
