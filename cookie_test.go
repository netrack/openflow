package of

import (
	"bytes"
	"io"
	"testing"

	"github.com/netrack/openflow/encoding"
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
		t.Fatalf("Invalid cookie returned: %d", jar.Cookies)
	}
}
