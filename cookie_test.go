package openflow

import (
	"bytes"
	"io"
	"testing"

	"github.com/netrack/openflow/internal/encoding"
)

// fakeCookieJar implements the CookieJar interface.
type fakeCookieJar struct {
	Cookie uint64
}

func (f *fakeCookieJar) SetCookies(c uint64) {
	f.Cookie = c
}

func (f *fakeCookieJar) Cookies() uint64 {
	return f.Cookie
}

func (f *fakeCookieJar) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, &f.Cookie)
}

func (f *fakeCookieJar) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, f.Cookie)
}

func TestCookieReaderOf(t *testing.T) {
	// Write a fake cookie into the buffer.
	var rbuf bytes.Buffer
	encoding.WriteTo(&rbuf, uint64(42))

	// Create a new cookie reader from the cookie jar.
	cr := CookieReaderOf(&fakeCookieJar{})
	jar, err := cr.ReadCookie(&rbuf)

	if err != nil {
		t.Fatalf("Cookie read failed: %s", err)
	}

	if jar.Cookies() != 42 {
		t.Fatalf("Invalid cookie returned: %d", jar.Cookies())
	}
}

func TestNewCookieMatcher(t *testing.T) {
	jar := &fakeCookieJar{}
	m := NewCookieMatcher(jar)

	if m.Cookies != jar.Cookies() {
		t.Fatalf("Cookies are not set: %d != %d", m.Cookies, jar.Cookies())
	}

	// Create a request that we are going to match with the matcher.
	r := NewRequest(TypeHello, &fakeCookieJar{jar.Cookie})
	if !m.Match(r) {
		t.Fatalf("Request expected to match the cookie: %d", jar.Cookie)
	}

	// Create a new request with a message, that have distinct cookie.
	r = NewRequest(TypeHello, &fakeCookieJar{jar.Cookie + 1})
	if m.Match(r) {
		t.Fatalf("Request is not expected to match the cookie: %d", jar.Cookie)
	}
}
