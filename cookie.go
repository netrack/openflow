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

type filterEntry struct {
	handler   Handler
	evictable bool
}

// CookieHandler provides mechanism to hook up the message handler with an
// opaque randomly created data. Handler is safe for concurrent use by
// multiple goroutines.
type CookieHandler struct {
	Reader CookieReader

	rand *rand.Rand

	handlers map[uint64]*filterEntry
	lock     sync.RWMutex
}

func NewCookieHandler() *CookieHandler {
	seed := time.Now().UTC().UnixNano()

	return &CookieHandler{
		handlers: make(map[uint64]*filterEntry),
		rand:     rand.New(rand.NewSource(seed)),
	}
}

func (h *CookieHandler) Handle(jar CookieJar, handler Handler) {
	cookies := uint64(h.rand.Int63())
	jar.SetCookies(cookies)

	h.lock.Lock()
	defer h.lock.Unlock()

	h.handlers[cookies] = &filterEntry{handler, false}
}

func (h *CookieHandler) HandleFunc(jar CookieJar, handler HandlerFunc) {
	h.Handle(jar, handler)
}

func (h *CookieHandler) Unhandle(jar CookieJar) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.handlers, jar.Cookies())
}

// Serve implements Handler interface. Serve dispatches the request to the
// handler whose cookie matches.
func (h *CookieHandler) Serve(rw ResponseWriter, r *Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	jar, err := h.Reader.ReadCookie(bytes.NewBuffer(body))
	if err != nil {
		return
	}

	h.lock.RLock()
	defer h.lock.RUnlock()

	entry, ok := h.handlers[jar.Cookies()]
	if !ok {
		return
	}

	if entry.evictable {
		delete(h.handlers, jar.Cookies())
	}

	r.Body = bytes.NewBuffer(body)
	entry.handler.Serve(rw, r)
}
