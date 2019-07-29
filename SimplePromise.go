package promise

import (
	"log"
)

type packetChannel chan Packet

// TypeA is a Packet channel and implements a classic Promise
type TypeA struct {
	packets packetChannel
	pause   *chan bool
}

/*
	newTypeA makes a packetChannel-based promise
	and allows the optional pause channel
*/
func newTypeA(ch *chan bool) TypeA {
	if ch == nil {
		x := make(chan bool)
		ch = &x
	}
	return TypeA{
		make(packetChannel),
		ch,
	}
}

// Resolve a new promise
func Resolve(data interface{}) Promise {
	out := newTypeA(nil)
	go func() { out.packets <- NewDataPacket(data) }()
	return out
}

// Reject returns a promise in an erroring state
func Reject(err error) Promise {
	out := newTypeA(nil)
	go func() { out.packets <- NewErrorPacket(err) }()
	return out
}

// Paused returns a promise in an erroring state
func Paused(data interface{}) Promise {
	// pauseChannel := make(chan bool)
	out := newTypeA(nil) // &pauseChannel)
	//go func() { pauseChannel <- true }()
	go func() { out.packets <- NewPausedPacket(data) }()
	return out
}

// Then is the classic then you know and love
func (in TypeA) Then(c Callback) Promise {
	out := newTypeA(in.pause)

	go func() {
		inPacket := <-in.packets
		if inPacket.err != nil {
			out.packets <- inPacket
		} else {
			if inPacket.paused {
				<-(*in.pause)
			}
			out.packets <- Packet{
				value:  c(inPacket.value),
				paused: inPacket.paused,
			}
		}
	}()

	return out
}

// Catch catches an error - but you can't rethrow errors
func (in TypeA) Catch(c ErrorCallback) Promise {
	out := newTypeA(in.pause)
	go func() {
		inPacket := <-in.packets
		if inPacket.err == nil {
			out.packets <- inPacket
		} else {
			if inPacket.paused {
				<-(*in.pause)
			}
			out.packets <- Packet{
				value:  c(inPacket.err),
				paused: inPacket.paused,
			}
		}
	}()
	return out
}

// Step unpauses to execute one function
func (in TypeA) Step() PausedPromise {
	out := newTypeA(in.pause)
	go func() {
		(*in.pause) <- true
	}()
	return out
}

// Allocate calls the step function n times,
// allowing a paused promise to behave normally for n
// iterations
func (in TypeA) Allocate(n int) PausedPromise {
	go func() {
		for i := 0; i < n; i++ {
			in.Step()
		}
	}()
	return in
}

// Deallocate clears the pause channel
func (in TypeA) Deallocate() PausedPromise {
	out := newTypeA(nil)
	go func() {
		inPacket := <-in.packets
		out.packets <- inPacket
	}()
	return out
}

// Continue causes this promise to behave like a normal one again
func (in TypeA) Continue() Promise {
	completed := make(chan bool)
	in.Then(func(v interface{}) interface{} {
		completed <- true
		return nil
	})
	out := newTypeA(nil)
	go func() {
		for {
			select {
			case <-completed:
				break
			default:
				in.Step()
			}
		}
	}()
	return out
}

// Resolver is a function signature that takes anything and turns it into a promise
type Resolver func(data interface{})

// Rejecter is a function signature that takes an error and turns it into an erroring promise
type Rejecter func(err error)

// New lets you wrap a promise over some asynchronous task
func New(initial func(Resolver, Rejecter)) TypeA {
	out := newTypeA(nil)

	res := func(data interface{}) {
		go func() {
			out.packets <- Packet{
				value: data,
			}
		}()
	}
	rej := func(err error) {
		go func() {
			out.packets <- Packet{
				err: err,
			}
		}()
	}
	go initial(res, rej)
	return out
}

// Log just logs what comes to it and send it on it's way
func Log(log *log.Logger) Callback {
	return func(v interface{}) interface{} {
		log.Print(v)
		return v
	}
}
