package openflow

import (
	"sync"
)

// Runner describes types used to start a function according to the
// defined concurrency model.
type Runner interface {
	Run(func())
}

// OnDemandRoutineRunner is a runner that starts each function in a
// separate goroutine. This handler is useful for initial prototyping,
// but it is highly recommended to use runner with a fixed amount of
// workers in order to prevent over goroutining (see MultiRoutineRunner).
type OnDemandRoutineRunner struct{}

// Run starts a function in a separate go-routine. This method implements
// Runner interface.
func (_ OnDemandRoutineRunner) Run(fn func()) {
	go fn()
}

// SequentialRunner is a runner that starts each function one by one.
// New function does not start execution until the previous one is done.
//
// This runner is useful for debugging purposes.
type SequentialRunner struct{}

// Run starts a function as is. This method implements Runner interface.
func (_ SequentialRunner) Run(fn func()) {
	fn()
}

// MultiRoutineRunner is a runner that assigns each function to one
// of the workers from the pool. So there is always a constant amount
// of goroutines.
type MultiRoutineRunner struct {
	num  int
	q    chan func()
	once sync.Once
}

// NewMultiRoutineRunner creates a new instance of MultiRoutinerRunner
// with a specified amount of workers. Method panics when number is
// not positive.
func NewMultiRoutineRunner(num int) *MultiRoutineRunner {
	if num <= 0 {
		panic("number of routines must be positive")
	}
	return &MultiRoutineRunner{
		num: num,
		q:   make(chan func(), num),
	}
}

// init starts all workers.
func (mrr *MultiRoutineRunner) init() {
	for i := 0; i < mrr.num; i++ {
		go mrr.runner()
	}
}

func (mrr *MultiRoutineRunner) runner() {
	for fn := range mrr.q {
		fn()
	}
}

// Run puts a function in the waiting queue and exists. This method
// returns control to the parent caller.
func (mrr *MultiRoutineRunner) Run(fn func()) {
	mrr.once.Do(mrr.init)
	mrr.q <- fn
}
