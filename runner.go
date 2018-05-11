package openflow

// Runner describes types used to start a function according to the
// defined concurrency model.
type Runner interface {
	Run(func())
}

// OnDemandRoutineRunner is a runner that starts each function in a
// separate goroutine. This handler is useful for initial prototyping,
// but it is highly recommended to use runner with a fixed amount of
// workers in order to prevent over goroutining.
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
