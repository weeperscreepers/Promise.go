package main

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

// Promise is a thin wrapper around an interface that we can hang functions on
type Promise struct {
	value interface{}
	err   error
}

// Resolve a new promise
func Resolve(data interface{}) Promise {
	return Promise{
		value: data,
	}
}

// Reject returns a promise in an erroring state
func Reject(err error) Promise {
	return Promise{
		err: err,
	}
}

// classic Promise.then()
func (promise Promise) then(c Callback) Promise {
	if promise.err != nil {
		return promise
	}

	ch := make(chan Promise)
	go func() {
		ch <- Resolve(c(promise.value))
	}()

	return <-ch
}

// In this implementation you cannot reject out of a .catch
func (promise Promise) catch(c ErrorCallback) Promise {
	if promise.err == nil {
		return promise
	}
	ch := make(chan Promise)
	go func() {
		ch <- Resolve(c(promise.err))
	}()

	return <-ch
}

func main() {
	go func() {
		time.Sleep(time.Second * 10)
		log.Print("is async real ?")
	}()
	go Resolve(2).
		then(func(v interface{}) interface{} {
			time.Sleep(time.Second * 2)
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

}
