// +build integration

package openflow_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	of "github.com/netrack/openflow"
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
	// Wait for 3 successfull echo-request replies and a
	// reply on the multipart aggregated statistics.
	nreqs := 4
	ch := make(chan struct{}, nreqs)
	mux := of.NewServeMux()

	statsHandler := func(rw of.ResponseWriter, r *of.Request) {
		ch <- struct{}{}
	}

	featuresHandler := func(rw of.ResponseWriter, r *of.Request) {
		var resp ofp.MultipartReply
		var features ofp.TableFeatures

		resp.ReadFrom(r.Body)
		features.ReadFrom(r.Body)

		// Read only the first table feature.
		log.Println(features)
	}

	helloHandler := func(rw of.ResponseWriter, r *of.Request) {
		rw.Write(r.Header.Copy(), nil)

		// Send a multipart request to retrieve statistics from
		// the connected switch.
		body := &ofp.AggregateStatsRequest{
			OutPort:  ofp.PortAny,
			OutGroup: ofp.GroupAny,
			Match: ofputil.ExtendedMatch(
				ofputil.MatchInPort(1),
			),
		}

		req := ofp.NewMultipartRequest(
			ofp.MultipartTypeAggregate, body)

		header := &of.Header{Type: of.TypeMultipartRequest}
		pattern := of.TransactionMatcher(header)

		// Create a matcher based on the header transaction. It
		// will be used to handle the multipart response.
		mux.HandleOnce(pattern, of.HandlerFunc(statsHandler))
		rw.Write(header, req)

		req = ofp.NewMultipartRequest(
			ofp.MultipartTypeTableFeatures, nil)

		header = &of.Header{Type: of.TypeMultipartRequest}
		pattern = of.TransactionMatcher(header)

		mux.Handle(pattern, of.HandlerFunc(featuresHandler))
		rw.Write(header, req)
	}

	echoHandler := func(rw of.ResponseWriter, r *of.Request) {
		var req ofp.EchoRequest
		req.ReadFrom(r.Body)

		header := r.Header.Copy()
		header.Type = of.TypeEchoReply

		rw.Write(header, &ofp.EchoReply{req.Data})
		ch <- struct{}{}
	}

	mux.HandleFunc(of.TypeMatcher(of.TypeHello), helloHandler)
	mux.HandleFunc(of.TypeMatcher(of.TypeEchoRequest), echoHandler)

	ln, _ := net.Listen("tcp", "127.0.0.1:6633")
	s := ofptest.NewUnstartedServer(mux, ln)

	s.Start()
	defer s.Close()

	if err := waitAll(ch, nreqs, 20*time.Second); err != nil {
		t.Fatal(err.Error())
	}
}
