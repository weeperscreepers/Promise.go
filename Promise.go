package Promise

import (
	"errors"
	"log"
	"time"
)

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

// classic Promise.then()
func (in Promise) then(c Callback) Promise {
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

// In this implementation you cannot reject out of a .catch
func (in Promise) catch(c ErrorCallback) Promise {
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

func main() {
	go func() {
		log.Print("is async real ?")
	}()
	go Resolve(2).
		then(func(v interface{}) interface{} {
			return v.(int) + 1
		}).
		then(func(v interface{}) interface{} {
			return v.(int) * 7
		}).
		then(func(v interface{}) interface{} {
			log.Print("Answer: ", v)
			return nil
		})

	go Reject(errors.New("This is an error")).
		catch(func(e error) interface{} {
			log.Print("There was an error: ", e)
			return "We recovered from the error"
		}).
		then(func(v interface{}) interface{} {
			log.Print("Finally: ", v)
			return nil
		})

	time.Sleep(time.Second * 10)
}
