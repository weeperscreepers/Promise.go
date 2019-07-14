package promise

// Callback s must be weakly typed because of Go type system
type Callback func(interface{}) interface{}

// ErrorCallback is a callback that recovers from an error
type ErrorCallback func(error) interface{}

// InOrder combines several Callbacks into one
func InOrder(callbacks []Callback) Callback {
	return func(initial interface{}) interface{} {
		val := initial
		for _, c := range callbacks {
			val = c(val)
		}
		return val
	}
}

// Packet is a thin wrapper around an interface that we can hang functions on
type Packet struct {
	value interface{}
	err   error
}

// Promise is a Packet channel
type Promise chan Packet

// Resolve a new promise
func Resolve(data interface{}) Promise {
	out := make(Promise)
	go func() {
		out <- Packet{
			value: data,
		}
	}()
	return out
}

// Reject returns a promise in an erroring state
func Reject(err error) Promise {
	out := make(Promise)
	go func() {
		out <- Packet{
			err: err,
		}
	}()
	return out
}

// Then is the classic then you know and love
func (in Promise) Then(c Callback) Promise {
	out := make(Promise)
	promise := <-in
	if promise.err != nil {
		go func() {
			out <- promise
		}()
		return out
	}

	return Resolve(c(promise.value))
}

// Catch catches an error - but you can't rethrow errors
func (in Promise) Catch(c ErrorCallback) Promise {
	out := make(Promise)
	promise := <-in
	if promise.err == nil {
		go func() {
			out <- promise
		}()
		return out
	}

	return Resolve(c(promise.err))
}
