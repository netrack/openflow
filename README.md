# openflow - The OpenFlow protocol library

The openflow library is a pure Go implementation of the OpenFlow protocol. The ideas of the programming interface mostly borrowed from the Go standard HTTP library.

# Installation

```bash
$ go get github.com/netrack/openflow
```

# Usage

The usage is pretty similar to the handling HTTP request, but instead of routes we are using message types.

```go
// Define the OpenFlow handler for hello messages.
of.HandleFunc(of.TypeHello, func(rw of.ResponseWriter, r *of.Request) {
        rw.WriteHeader(&of.Header{Type: of.TypeHello})
})

// Start the TCP server on 6633 port.
of.ListenAndServe(":6633", nil)
```

```go
tm := &TypeMatcher{Type: TypePacketIn}
m := &of.RequestMatcher{tm}

d := NewRequestDispatcher()
d.HandleFunc(m, func(rw of.ResponseWriter, r *of.Request) {
        // Send back hello response.
        rw.WriteHeader(&of.Header{Type: of.TypeHello})
})
```

# License

The openflow library is distributed under MIT license, therefore you are free to do with code whatever you want. See the [LICENSE](LICENSE) file for full license text.
