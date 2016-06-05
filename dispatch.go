package of

import (
	"bytes"
)

type Filter interface {
	Filter(*Request) bool
}

type FilterFunc func(*Request) bool

func (fn FilterFunc) Filter(r *Request) bool {
	return fn(r)
}

// TypeDispatcher represents a filter
type TypeFilter struct {
	Type Type
}

func (f *TypeFilter) Filter(r *Request) bool {
	return r.Header.Type == f.Type
}

type Matcher interface {
	Match(r *Request) bool
}

type RequestMatcher struct {
	Filters []Filter

	Disposable bool
}

func (d *RequestMatcher) Filter(f Filter) {
	d.Filters = append(d.Filters, f)
}

func (d *RequestMatcher) Match(r *Request) (b bool) {
	for _, filter := range d.Filters {
		b = filter.Filter(r)
		if !b {
			return
		}
	}
}

// Dispatcher is an interface used to select the matching handler to
// process the received OpenFlow message.
type Dispatcher interface {
	Dispatch(*Request) (Handler, Matcher)
}

type RequestDispatcher struct {
	lock     sync.RWMutex
	handlers map[Matcher]Handler
}

func NewRequestDispatcher() *RequestDispatcher {
	return &RequestDispatcher{handlers: make(map[Matcher]Handler)}
}

func (d *RequestDispatcher) Handle(m Matcher, h Handler) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if m == nil {
		panic("openflow: request dispatcher nil matcher")
	}

	d.handlers[m] = h
}

func (d *RequestDispatcher) HandleFunc(m Matcher, f HandlerFunc) {
	d.Handle(m, f)
}

func (d *RequestDispatcher) Dispatch(r *Request) (Handler, Matcher) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entries, ok := d.m[r.Header.Type]
	if !ok {
		entries = append(entries, &serveMuxEntry{h: DiscardHandler})
	}

	return HandlerFunc(func(rw ResponseWriter, r *Request) {
		if len(entities) == 0 {
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		for _, entry := range entries {
			r.Body = bytes.NewBuffer(body)
			go entry.h.Serve(rw, r)
		}
	})
}

// TypeDispatcher is an OpenFlow request multiplexer. It matches the type
// of the OpenFlow message against a list of registered handlers and calls
// the marching handler.
type TypeDispatcher struct {
	dispatcher RequestDispatcher
}

// NewTypeDispatcher allocates a new instance of the TypeDispatcher.
func NewTypeDispatcher() *TypeDispatcher {
	return &TypeDispatcher{dispatcher: NewRequestDispatcher()}
}

// Handle registers the handler for the given message type.
func (d *TypeDispatcher) Handle(t Type, h Handler) {
	dispatcher := &RequestMatcher{Disposable: false}
	dispatcher.Filter(&TypeFilter{t})

	d.dispatcher.Handle(dispatcher, h)
}

// Handler registers handler function on the given OpenFlow
// message type.
func (d *TypeDispatcher) HandleFunc(t Type, f HandlerFunc) {
	d.Handle(t, f)
}

// Dispatch returns a Handler instance for the given OpenFlow request.
func (d *TypeDispatcher) Dispatch(r *Request) (Handler, Dispatcher) {
	return d.dispatcher.Dispatch(r)
}

// Serve dispatches OpenFlow requests to the registered handlers.
func (d *TypeDispatcher) Serve(rw ResponseWriter, r *Request) {
	h, _ := d.Handler(r)
	h.Serve(rw, r)
}

/*

match := &of.RequestMatcher{Disposable: true}
match.Filter(&of.TypeFilter{TypePacketIn})
match.Filter(&of.CookieFilter{0x123abc, &of.CookieReader{&ofp.PacketIn{}}})

dispatcher := NewRequestDispatcher()
dispatcher.HandleFunc(match, func(rw of.ResponseWriter, r *of.Request) {
	rw.WriteHeader()
})

*/
