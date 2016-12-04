package openflow

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
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

// ReadCookie calls the function with the specifier reader argument.
func (fn CookieReaderFunc) ReadCookie(r io.Reader) (CookieJar, error) {
	return fn(r)
}

// CookieMatcher provides mechanism to hook up the message handler with an
// opaque data. Filter is safe for concurrent use by multiple goroutines.
type CookieMatcher struct {
	Cookies uint64

	// Reader is an OpenFlow message unmarshaler. CookieFilter will use
	// it to access the request cookie value. If the cookie matches, the
	// registered handler will be called to process the request. Otherwise
	// the request will be skipped.
	Reader CookieReader
}

// NewCookieMatcher creates a new cookie matcher. The cookie value will
// be randomly generated using the functions from the standard library.
func NewCookieMatcher(j CookieJar) *CookieMatcher {
	var cookies uint64

	// The default random generator does not allow to create unsigned
	// 64-bit integers, therefore, we will glue that type from the
	// two 32-bit pieces.
	cookiesLow := rand.Uint32()
	cookiesHigh := rand.Uint32()

	cookies = uint64(cookiesHigh<<32) | uint64(cookiesLow)

	// Set the generated cookies to the given cookie jar and also
	// put this value to the matcher.
	j.SetCookies(cookies)
	return &CookieMatcher{cookies, CookieReaderOf(j)}
}

// Match compares the cookie from the message with the given one.
//
// Cookie of each incoming request will be compared to the given cookie
// jar cookie. If the request cookie matches the registered one, the given
// handler will be used to process the request.
func (f *CookieMatcher) Match(r *Request) bool {
	rd, ok := r.Body.(*bytes.Buffer)

	if ok {
		// If the body is a bytes buffer, we will simply reset
		// it to reduce the amount of memory to allocate.
		defer rd.Reset()
	} else {
		// Otherwise, we will re-create a new one.
		body, err := ioutil.ReadAll(r.Body)
		defer func() { r.Body = bytes.NewBuffer(body) }()

		if err != nil {
			return false
		}

		rd = bytes.NewBuffer(body)
	}

	// Parse the incoming request to access the cookies.
	jar, err := f.Reader.ReadCookie(rd)
	if err != nil {
		return false
	}

	return jar.Cookies() == f.Cookies
}

// CookieReaderOf creates a new cookie reader instance from the cookie
// jar. It uses reflection to create a new examplar of the given type,
// so the resulting reader is safe to use in multiple go-routines.
func CookieReaderOf(j CookieJar) CookieReader {
	valueType := reflect.TypeOf(j).Elem()

	cr := func(r io.Reader) (CookieJar, error) {
		value := reflect.New(valueType)
		jar, ok := value.Interface().(io.ReaderFrom)

		// The CookieJar have to implement the io.ReaderFrom
		// interface, so it could be unmarshaled.
		if !ok {
			message := "openflow: not a valid cookie reader"
			return nil, fmt.Errorf(message)
		}

		// Unmarshal the cookie jar and return it.
		_, err := jar.ReadFrom(r)
		return jar.(CookieJar), err
	}

	return CookieReaderFunc(cr)
}
