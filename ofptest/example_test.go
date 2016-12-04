package ofptest_test

import (
	"fmt"
	"log"

	"github.com/netrack/openflow/ofp"
	"github.com/netrack/openflow/ofptest"
	of "github.com/netrack/openflow"
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

func ExampleServer() {
	ts := ofptest.NewServer(of.HandlerFunc(func(w of.ResponseWriter, r *of.Request) {
		res := &ofp.EchoReply{Data: []byte("pong")}
		w.Write(&of.Header{Type: of.TypeEchoReply}, res)
	}))

	defer ts.Close()

	echoReq := &ofp.EchoRequest{Data: []byte("ping")}
	req := of.NewRequest(of.TypeEchoRequest, echoReq)

	conn, err := of.Dial("tcp", ts.Listener.Addr().String())
	if err != nil {
		log.Fatal(err)
	}

	conn.Send(req)
	conn.Flush()
	resp, _ := conn.Receive()

	var echoResp ofp.EchoReply
	echoResp.ReadFrom(resp.Body)

	fmt.Printf("%s", echoResp.Data)
	// Output: pong
}
