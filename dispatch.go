package of

import (
	"sync"
)

// Dispatcher is an interface used to select the matching handler to
// process the received OpenFlow message.
type Dispatcher interface {
	// Dispatch returns the best matching handler for the given request.
	Dispatch(*Request) bool
}

type TypeDispatcher struct {
	mu sync.RWMutex
	m  map[Type][]Handler
}

// NewTypeDispacher creates a new instance of the message-type based
// dispather.
func NewTypeDispatcher() *TypeDispatcher {
	return &TypeDispatcher{m: make(map[Type][]Handler)}
}

func (d *TypeDispatcher) Handle(t Type, h Handler) {
}

func (d *TypeDispatcher) HandleFunc(t Type, f HandlerFunc) {
	d.Handle(t, f)
}

type CookieDispatcher struct {
}

func (d *CookieDispatcher) Handle(c *Cookie, h Handler) {
}

func (d *CookieDispatcher) HandleFunc(c *Cookie, f HandlerFunc) {
	d.Handle(t, f)
}
