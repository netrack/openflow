package of

import (
	"sync"
)

type Matcher interface {
	Match(*Request) bool
}

type MatcherFunc func(*Request) bool

func (fn MatcherFunc) Match(r *Request) bool {
	return fn(r)
}

type TypeMatcher struct {
	Type Type
}

func (f *TypeMatcher) Match(r *Request) bool {
	return r.Header.Type == f.Type
}

type RequestMatcher []Matcher

func (d *RequestMatcher) Match(r *Request) bool {
	for _, matcher := range *d {
		if !matcher.Match(r) {
			return false
		}
	}

	return true
}

// Dispatcher is an interface used to select the matching handler to
// process the received OpenFlow message.
type Dispatcher interface {
	Dispatch(*Request) Handler
}

type DispatchHandler interface {
	Dispatcher
	Handler
}

type dispatchEntry struct {
	matcher Matcher
	handler Handler
}

type RequestDispatcher struct {
	lock     sync.RWMutex
	handlers []*dispatchEntry
}

func NewRequestDispatcher() *RequestDispatcher {
	return &RequestDispatcher{}
}

func (d *RequestDispatcher) Handle(m Matcher, h Handler) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if m == nil {
		panic("openflow: request dispatcher nil matcher")
	}

	d.handlers = append(d.handlers, &dispatchEntry{m, h})
}

func (d *RequestDispatcher) HandleFunc(m Matcher, h HandlerFunc) {
	d.Handle(m, h)
}

func (d *RequestDispatcher) Dispatch(r *Request) Handler {
	d.lock.RLock()
	defer d.lock.RUnlock()

	// Try to match the processing request to any of the
	// registered filters.
	for _, entry := range d.handlers {
		if !entry.matcher.Match(r) {
			continue
		}

		return entry.handler
	}

	// When the handler was not found, we have to return
	// the discard handler.
	return DiscardHandler
}

// TypeDispatcher is an OpenFlow request multiplexer. It matches the type
// of the OpenFlow message against a list of registered handlers and calls
// the marching handler.
type TypeDispatcher struct {
	dispatcher *RequestDispatcher
}

func NewTypeDispatcher() *TypeDispatcher {
	return &TypeDispatcher{NewRequestDispatcher()}
}

// Handle registers the handler for the given message type.
func (d *TypeDispatcher) Handle(t Type, h Handler) {
	d.dispatcher.Handle(&TypeMatcher{t}, h)
}

// HandleFunc registers handler function on the given OpenFlow
// message type.
func (d *TypeDispatcher) HandleFunc(t Type, f HandlerFunc) {
	d.Handle(t, f)
}

// Dispatch returns a Handler instance for the given OpenFlow request.
func (d *TypeDispatcher) Dispatch(r *Request) Handler {
	return d.dispatcher.Dispatch(r)
}

// Serve implements Handler internface. It processing the request and
// writes back the response.
func (d *TypeDispatcher) Serve(rw ResponseWriter, r *Request) {
	h := d.Dispatch(r)
	h.Serve(rw, r)
}

// DefaultDispatcher is an instance of the TypeDispatcher used as
// a default handler in the DefaultServer instance.
var DefaultDispatcher = NewTypeDispatcher()

// Handle registers the handler on the given type of the OpenFlow
// message in the DefaultServeMux.
func Handle(t Type, handler Handler) {
	DefaultDispatcher.Handle(t, handler)
}

// HandleFunc registers the handler function on the given type of
// the OpenFlow message in the DefaultServeMux.
func HandleFunc(t Type, f func(ResponseWriter, *Request)) {
	DefaultDispatcher.HandleFunc(t, f)
}
