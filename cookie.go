package of

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

type CookieJar interface {
	SetCookies(uint64)
	Cookies() uint64
}

// CookieReader is the interface to read cookie jars.
//
// CookieReader parses the body of the handling request and returns the
// cookie jar with containing cookies or nil when error occurs.
type CookieReader interface {
	ReadCookie(io.Reader) (CookieJar, error)
}

// The CookieReaderFunc is an adapter to allow use of ordinary functions
// as OpenFlow handlers. If fn is a function with the appropriate signature,
// CookieReaderFunc(fn) is a Reader that calls fn.
type CookieReaderFunc func(io.Reader) (CookieJar, error)

// CookieJar calls the function with the specifier reader argument.
func (fn CookieReaderFunc) CookiesJar(r io.Reader) (CookieJar, error) {
	return fn(r)
}

// CookieMux provides mechanism to hook up the message handler with an
// opaque data. Filter is safe for concurrent use by multiple goroutines.
type CookieFilter struct {
	Cookies uint64

	// Reader is an OpenFlow message unmarshaler. CookieMux will use the
	// it to access the request cookie value. If the cookie matches, the
	// registered handler will be called to process the request. Otherwise
	// the request will be skipped.
	Reader CookieReader
}

// Filter compares the cookie from the message with the given one.
//
// Cookie of each incoming request will be compared to the given cookie
// jar cookie. If the request cookie matches the registered one, the given
// handler will be used to process the request.
func (f *CookieFilter) Filter(r *Request) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Parse the incoming request to access the cookies.
	jar, err := f.Reader.ReadCookie(bytes.NewBuffer(body))
	if err != nil {
		return false
	}

	r.Body = bytes.NewBuffer(body)
	return jar.Cookies() == f.Cookies
}
