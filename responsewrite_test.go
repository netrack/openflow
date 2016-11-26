// +build integration

package of_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp"
	"github.com/netrack/openflow/ofptest"
	"github.com/netrack/openflow/ofputil"
)

func waitAll(ch chan struct{}, n int, t time.Duration) error {
	timer := time.NewTimer(t)
	defer timer.Stop()

	for i := 0; i < n; i++ {
		select {
		case <-ch:
		case <-timer.C:
			text := "no response within the time interval"
			return fmt.Errorf(text)
		}
	}

	return nil
}

func TestResponseWrite(t *testing.T) {
	// Wait for 3 successfull echo-request replies.
	nreqs := 3
	ch := make(chan struct{}, nreqs)

	helloHandler := func(rw of.ResponseWriter, r *of.Request) {
		rw.Write(r.Header.Copy(), nil)

		// Send a few multipart requests to retrieve
		// statistics from the connected switch.
		req := &ofp.MultipartRequest{
			Type:  ofp.MultipartTypeAggregate,
			Flags: ofp.MultipartRequestMode,
		}

		stats := of.MultiWriterTo(req, &ofp.AggregateStatsRequest{
			OutPort: ofp.PortAny, OutGroup: ofp.GroupAny,
			Match: ofputil.ExtendedMatch(ofputil.MatchInPort(1)),
		})

		header := r.Header.Copy()
		header.Type = of.TypeMultipartRequest

		rw.Write(header, stats)
	}

	echoHandler := func(rw of.ResponseWriter, r *of.Request) {
		var req ofp.EchoRequest
		req.ReadFrom(r.Body)

		header := r.Header.Copy()
		header.Type = of.TypeEchoReply

		rw.Write(header, &ofp.EchoReply{req.Data})
		ch <- struct{}{}
	}

	td := of.NewTypeDispatcher()
	td.HandleFunc(of.TypeHello, helloHandler)
	td.HandleFunc(of.TypeEchoRequest, echoHandler)

	ln, _ := net.Listen("tcp", ":6633")
	s := ofptest.NewUnstartedServer(td, ln)

	s.Start()
	defer s.Close()

	if err := waitAll(ch, nreqs, 20*time.Second); err != nil {
		t.Fatal(err.Error())
	}
}
