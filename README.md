# openflow - The OpenFlow protocol library

[![Build Status][BuildStatus]](https://travis-ci.org/netrack/openflow)
[![Documentation][Documentation]](https://godoc.org/github.com/netrack/openflow)

The openflow library is a pure Go implementation of the OpenFlow protocol.
The ideas of the programming interface mostly borrowed from the Go standard
HTTP library.

# Installation

```bash
$ go get github.com/netrack/openflow
```

# Usage

The usage is pretty similar to the handling HTTP request, but instead of routes
we are using message types.

```go
// Define the OpenFlow handler for hello messages.
openflow.HandleFunc(openflow.TypeHello, func(rw openflow.ResponseWriter, r *openflow.Request) {
        // Send back hello response.
        rw.Write(&openflow.Header{Type: openflow.TypeHello}, nil)
})

// Start the TCP server on 6633 port.
openflow.ListenAndServe(":6633", nil)
```

```go
pattern := openflow.TypeMatcher(openflow.TypePacketIn)

mux := openflow.NewServeMux()
mux.HandleFunc(pattern, func(rw openflow.ResponseWriter, r *openflow.Request) {
        // Send back the packet-out message.
        rw.Write(&openflow.Header{Type: openflow.TypePacketOut}, nil)
})
```

# License

The openflow library is distributed under MIT license, therefore you are free
to do with code whatever you want. See the [LICENSE](LICENSE) file for full
license text.


[BuildStatus]:   https://travis-ci.org/netrack/openflow.svg?branch=master
[Documentation]: https://godoc.org/github.com/netrack/openflow?status.svg
