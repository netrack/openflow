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

type Baker interface {
	Cookies(io.Reader) (CookieJar, error)
}

type BakerFunc func(io.Reader) (CookieJar, error)

func (fn BakerFunc) Cookies(r io.Reader) (CookieJar, error) {
	return fn(r)
}

type filterEntry struct {
	handler   Handler
	evictable bool
}

type CookieFilter struct {
	Baker Baker

	rand *rand.Rand

	handlers map[uint64]*filterEntry
	lock     sync.RWMutex
}

func NewCookieFilter() *CookieFilter {
	seed := time.Now().UTC().UnixNano()

	return &CookieFilter{
		handlers: make(map[uint64]*filterEntry),
		rand:     rand.New(rand.NewSource(seed)),
	}
}

func (m *CookieFilter) Filter(jar CookieJar, handler Handler) {
	cookies := uint64(m.rand.Int63())
	jar.SetCookies(cookies)

	m.lock.Lock()
	defer m.lock.Unlock()

	m.handlers[cookies] = &filterEntry{handler, false}
}

func (m *CookieFilter) FilterFunc(jar CookieJar, handler HandlerFunc) {
	m.Filter(jar, handler)
}

func (m *CookieFilter) Serve(rw ResponseWriter, r *Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	jar, err := m.Baker.Cookies(bytes.NewBuffer(body))
	if err != nil {
		return
	}

	m.lock.RLock()

	entry, ok := m.handlers[jar.Cookies()]
	if !ok {
		m.lock.RUnlock()
		return
	}

	if entry.evictable {
		m.lock.RUnlock()
		delete(m.handlers, jar.Cookies())
	}

	m.lock.RUnlock()

	r.Body = bytes.NewBuffer(body)
	entry.handler.Serve(rw, r)
}
