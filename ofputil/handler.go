package ofputil

import (
	"log"

	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp"
)

// EchoHandler returns a request handler that replies on
// each request with a echo message with the same data
// as it was retrieved in the original message.
//
// The method accepts optional handler, that will executed
// in case of successful message submission.
func EchoHandler(h of.Handler) of.Handler {
	fn := func(rw of.ResponseWriter, r *of.Request) {
		var req ofp.EchoRequest

		// Try to parse the retrieved request to copy
		// the data of echo message into the request.
		_, err := req.ReadFrom(r.Body)
		if err != nil {
			text := "ofputil: failed to read the message: %v"
			log.Printf(text, err)
			return
		}

		header := r.Header.Copy()
		header.Type = of.TypeEchoReply

		// Send a reply with the same data in body.
		rw.Write(header, &ofp.EchoReply{req.Data})

		// Execute optional handler.
		if h != nil {
			h.Serve(rw, r)
		}
	}

	return of.HandlerFunc(fn)
}

// HelloHandler returns a simple request handler that replies
// to each request with hello message of the specified version.
//
// The method accepts optional handler, that will executed
// in case of successful message submission.
func HelloHandler(version uint8, h of.Handler) of.Handler {
	fn := func(rw of.ResponseWriter, r *of.Request) {
		// Copy the header of retrieved message,
		// including the trasnaction identifier.
		header := r.Header.Copy()
		header.Version = version

		// Send a response to the called with a version
		// supported version in the header.
		rw.Write(header, nil)

		if h != nil {
			h.Serve(rw, r)
		}
	}

	return of.HandlerFunc(fn)
}
