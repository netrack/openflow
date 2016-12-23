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
package main

import (
    of "github.com/netrack/openflow"
)

func main() {
    // Define the OpenFlow handler for hello messages.
    of.HandleFunc(of.TypeHello, func(rw of.ResponseWriter, r *of.Request) {
        // Send back hello response.
        rw.Write(&of.Header{Type: of.TypeHello}, nil)
    })

    // Start the TCP server on 6633 port.
    of.ListenAndServe(":6633", nil)
}
```

```go
package main

import (
    "github.com/netrack/openflow/ofp"
    of "github.com/netrack/openflow"
)

func main() {
    pattern := of.TypeMatcher(of.TypePacketIn)

    mux := of.NewServeMux()
    mux.HandleFunc(pattern, func(rw of.ResponseWriter, r *of.Request) {
        var packet ofp.PacketIn
        packet.ReadFrom(r.Body)

        apply := &ofp.InstructionApplyActions{
            ofp.Actions{&ofp.ActionOutput{ofp.PortFlood, 0}},
        }

        // For each incoming packet-in request, create a
        // respective flow modification command.
        fmod := ofp.NewFlowMod(ofp.FlowAdd, packet)
        fmod.Instructions = ofp.Instructions{apply}

        rw.Write(&of.Header{Type: of.TypeFlowMod}, fmod)
    })

    of.ListenAndServe(":6633", mux)
}
```

# License

The openflow library is distributed under MIT license, therefore you are free
to do with code whatever you want. See the [LICENSE](LICENSE) file for full
license text.


[BuildStatus]:   https://travis-ci.org/netrack/openflow.svg?branch=master
[Documentation]: https://godoc.org/github.com/netrack/openflow?status.svg
