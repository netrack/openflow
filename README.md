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
    rw.Header().Set(of.TypeHeaderKey, of.TypeHello)
    rw.WriteHeader()
})

// Start the TCP server on 6633 port.
of.ListenAndServe()
```

```go
match := &of.RequestMatcher{Disposable: true}
match.Filter(&of.TypeFilter{TypePacketIn})
match.Filter(&of.CookieFilter{0x123abc, &of.CookieReader{&ofp.PacketIn{}}})

dispatcher := NewRequestDispatcher()
dispatcher.HandleFunc(match, func(rw of.ResponseWriter, r *of.Request) {
	rw.WriteHeader()
})
```

# License

The openflow library is distributed under MIT license, therefore you are free to do with code whatever you want. See the [LICENSE](LICENSE) file for full license text.
