package promise

// Callback s must be weakly typed because of Go type system
type Callback func(interface{}) interface{}

// ErrorCallback is a callback that recovers from an error
type ErrorCallback func(error) interface{}

// Packet is a thin wrapper around an interface that we can hang functions on
type Packet struct {
	value        interface{}
	err          error
	paused       bool
	pauseChannel chan interface{}
}

// Promise is an interface that supplies two basic chaining operations
type Promise interface {
	Then(c Callback) Promise
	Catch(c ErrorCallback) Promise
}

// PausedPromise is a promise that can be step()ed
type PausedPromise interface {
	Promise
	// Continue turns this back into
	// a normal promise that will execute automatically
	Continue() Promise
	// Step the promise one time
	Step() PausedPromise
	// Allocate (n) is equivalent to calling step n times
	Allocate(n int) PausedPromise
	// Deallocate resets the Step() and Allocate() functionality
	Deallocate() PausedPromise
}

// NewDataPacket constructs an instance Packet
func NewDataPacket(data interface{}) Packet {
	return Packet{
		value:  data,
		paused: false,
	}
}

// NewErrorPacket constructs an instance Packet
func NewErrorPacket(err error) Packet {
	return Packet{
		err:    err,
		paused: false,
	}
}

// NewPausedPacket constructs an instance Packet
func NewPausedPacket(data interface{}) Packet {
	return Packet{
		value:  data,
		paused: true,
	}
}
