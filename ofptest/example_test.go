package ofptest

import (
	"fmt"
	"log"

	of "github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp"
)

func ExampleResponseRecorder() {
	handler := func(w of.ResponseWriter, r *of.Request) {
		w.Write(r.Header.Copy(), nil)
	}

	req := of.NewRequest(of.TypeHello, nil)
	w := NewRecorder()

	handler(w, req)
	fmt.Printf("type: %d", w.First().Header.Type)
	// Output: type: 0
}

func ExampleServer() {
	ts := NewServer(of.HandlerFunc(func(w of.ResponseWriter, r *of.Request) {
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
