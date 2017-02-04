package openflow

import (
	"fmt"
	"sync"
)

// A Matcher interface is used by multiplexer to find the handler for
// the received request.
type Matcher interface {
	Match(*Request) bool
}

// MatcherFunc type is an adapter to allow the use of ordinary functions
// as OpenFlow request matchers.
type MatcherFunc struct {
	Func func(*Request) bool
}

// Match implements Matcher interface and calls fn(r).
func (m *MatcherFunc) Match(r *Request) bool {
	return m.Func(r)
}

// TypeMatcher used to match requests by their types.
type TypeMatcher Type

// Match implements Matcher interface and matches the request by the type.
func (t TypeMatcher) Match(r *Request) bool {
	return r.Header.Type == Type(t)
}

// MultiMatcher creates a new Matcher instance that matches the request
// by all specified criterias.
func MultiMatcher(m ...Matcher) Matcher {
	fn := func(r *Request) bool {
		for _, matcher := range m {
			if !matcher.Match(r) {
				return false
			}
		}

		return true
	}

	return &MatcherFunc{fn}
}

type muxEntry struct {
	matcher Matcher
	handler Handler

	// once means handler will be executed only once, it will
	// be removed from the list of registered handlers.
	once bool
}

// ServeMux is an OpenFlow request multiplexer.
type ServeMux struct {
	mu       sync.RWMutex
	handlers map[Matcher]*muxEntry
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{handlers: make(map[Matcher]*muxEntry)}
}

// DefaultHandler is the Handler used for Requests that don't have a Matcher.
var DefaultHandler = DiscardHandler

// handle appends to the list of registered handlers a new one. If the
// matcher or handler of the entry is not defined, a panic function will
// be called.
func (mux *ServeMux) handle(e *muxEntry) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if e.matcher == nil {
		panic("openflow: nil matcher")
	}

	if e.handler == nil {
		panic("openflow: nil handler")
	}

	if _, dup := mux.handlers[e.matcher]; dup {
		text := "openflow: multiple registrations for %v"
		panic(fmt.Errorf(text, e.matcher))
	}

	mux.handlers[e.matcher] = e
}

// Handle registers the handler for the given pattern.
func (mux *ServeMux) Handle(m Matcher, h Handler) {
	mux.handle(&muxEntry{m, h, false})
}

// Handle registers disposable handler for the given pattern.
//
// It is not guaranteed that handler will process the first message
// exemplar of the matching message.
func (mux *ServeMux) HandleOnce(m Matcher, h Handler) {
	mux.handle(&muxEntry{m, h, true})
}

// HandleFunc registers handler function for the given pattern.
func (mux *ServeMux) HandleFunc(m Matcher, h HandlerFunc) {
	mux.Handle(m, h)
}

// HandleFunc registers handler function for the given message type.
func (mux *ServeMux) Handler(r *Request) Handler {
	var matcher Matcher
	var entry *muxEntry

	// Acquire a read lock to ensure that list would not
	// be modified during the search of registered handler.
	mux.mu.RLock()
	var matched bool

	// Try to match the processing request to any of the registered
	// filters.
	for matcher, entry = range mux.handlers {
		if matched = matcher.Match(r); matched {
			break
		}
	}

	mux.mu.RUnlock()

	// Use the DefaultHandler when there are no matching entries in the list.
	if !matched {
		return DefaultHandler
	}

	// If the retrieved entry is not disposable one, we will
	// return it as is without any processing.
	if !entry.once {
		return entry.handler
	}

	// But when the entry is disposable, we need to remove it
	// from the list.

	// We need to acquire the write lock in order to remove the
	// entry from the root. This procedure does not guarantee,
	// that handler will process the first message.
	mux.mu.Lock()
	defer mux.mu.Unlock()

	// If the concurrent message have already started the message
	// processing, it will be no longer presented in the list.
	if _, ok := mux.handlers[matcher]; !ok {
		return DiscardHandler
	}

	// Remove the entry from the list if it is marked as disposable.
	delete(mux.handlers, matcher)
	return entry.handler
}

// Serve implements Handler internface. It processing the request and
// writes back the response.
func (mux *ServeMux) Serve(rw ResponseWriter, r *Request) {
	h := mux.Handler(r)
	h.Serve(rw, r)
}

// TypeMux is an OpenFlow request multiplexer. It matches the type
// of the OpenFlow message against a list of registered handlers and calls
// the marching handler.
type TypeMux struct {
	mux *ServeMux
}

func NewTypeMux() *TypeMux {
	return &TypeMux{NewServeMux()}
}

// Handle registers the handler for the given message type.
func (mux *TypeMux) Handle(t Type, h Handler) {
	mux.mux.Handle(TypeMatcher(t), h)
}

// HandleOnce registers a disposable handler for the given message type.
func (mux *TypeMux) HandleOnce(t Type, h Handler) {
	mux.mux.HandleOnce(TypeMatcher(t), h)
}

// HandleFunc registers handler function for the given message type.
func (mux *TypeMux) HandleFunc(t Type, f HandlerFunc) {
	mux.Handle(t, f)
}

// Handle returns a Handler instance for the given OpenFlow request.
func (mux *TypeMux) Handler(r *Request) Handler {
	return mux.mux.Handler(r)
}

// Serve implements Handler internface. It processing the request and
// writes back the response.
func (mux *TypeMux) Serve(rw ResponseWriter, r *Request) {
	mux.mux.Serve(rw, r)
}

// DefaultMux is an instance of the TypeMux used as
// a default handler in the DefaultServer instance.
var DefaultMux = NewTypeMux()

// Handle registers the handler for the given message type message in the
// DefaultMux. The documentation for TypeMux
func Handle(t Type, handler Handler) {
	DefaultMux.Handle(t, handler)
}

// HandleOnce registers a disposable handler for the given message type
// in the DefaultMux.
func HandleOnce(t Type, handler Handler) {
	DefaultMux.HandleOnce(t, handler)
}

// HandleFunc registers the handler function on the given message type
// in the DefaultMux.
func HandleFunc(t Type, f func(ResponseWriter, *Request)) {
	DefaultMux.HandleFunc(t, f)
}
