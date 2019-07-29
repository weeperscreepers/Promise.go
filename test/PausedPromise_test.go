package promise_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weeperscreepers/promise"
)

type intChannel chan int

func (ch intChannel) constant(i int) promise.Callback {
	return func(v interface{}) interface{} {
		return i
	}
}

func (ch intChannel) close() promise.Callback {
	return func(v interface{}) interface{} {
		close(ch)
		return v
	}
}

func (ch intChannel) accept() promise.Callback {
	return func(v interface{}) interface{} {
		ch <- v.(int)
		return v
	}
}

func (ch intChannel) report(i int) promise.Callback {
	return func(v interface{}) interface{} {
		ch <- i
		return v
	}
}

func (ch intChannel) delete() promise.Callback {
	return func(v interface{}) interface{} {
		ch <- v.(int)
		return v
	}
}

func TestStep(t *testing.T) {
	ch := make(intChannel)

	// echo the channel number
	sayChannel0 := promise.Paused(0).
		Then(ch.accept()).(promise.PausedPromise)
	// echo the channel number and kill the channel
	sayChannel1 := promise.Paused(2).
		Then(ch.accept()).(promise.PausedPromise)

	sayChannel1.Step()
	i := <-ch
	assert.Equal(t, 2, i)

	sayChannel0.Step()
	i = <-ch
	assert.Equal(t, 0, i)

}

/*
	Promise should not blow up if there's nothing left to run,
	and we should be able to continue adding to it
*/
func TestEndOfPromise(t *testing.T) {
	ch := make(intChannel)
	sayChannel0 := promise.
		Paused(1234).
		Then(ch.accept()).
		Then(ch.constant(2)).(promise.PausedPromise)

	sayChannel0.Step()

	sayChannel0.Step()
	sayChannel0.Step()
	sayChannel0.Step()
	sayChannel0.Step()
	sayChannel0.Step()

	i := <-ch
	assert.Equal(t, 1234, i)
	sayChannel0.Then(ch.accept())
	i = <-ch
	assert.Equal(t, 2, i)
}

/*
	Since our Step() function has a weird behavior,
	we're gonna make it easier to do
*/
func TestAllocation(t *testing.T) {
	ch := make(intChannel)
	sayChannel0 := promise.Paused(0).(promise.PausedPromise)
	sayChannel0.
		Then(ch.accept()).
		Then(ch.report(1)).
		Then(ch.report(2))
	sayChannel1 := promise.Paused(100).(promise.PausedPromise)
	sayChannel1.
		Then(ch.accept()).
		Then(ch.report(99)).
		Then(ch.report(98)).
		Then(ch.report(97))

	sayChannel0.Allocate(2)

	i := <-ch
	assert.Equal(t, 0, i)
	i = <-ch
	assert.Equal(t, 1, i)

	sayChannel1.Allocate(3)

	i = <-ch
	assert.Equal(t, 100, i)
	i = <-ch
	assert.Equal(t, 99, i)
	i = <-ch
	assert.Equal(t, 98, i)
}

func TestNonExecution(t *testing.T) {

	// this should not happen
	closeChannel := func(ch chan int) promise.Callback {
		return func(v interface{}) interface{} {
			close(ch)
			return v
		}
	}

	const (
		M1 = 1234
		M2 = 4321
	)

	ch := make(intChannel)
	sayChannel0 := promise.Paused(M1).
		Then(ch.accept()).
		Then(closeChannel(ch)).(promise.PausedPromise)

	sayChannel1 := promise.Paused(M2).
		Then(ch.accept()).(promise.PausedPromise)
	sayChannel0.Step()
	log.Print("waiting a few seconds...")
	time.Sleep(time.Second * 3)
	i := <-ch
	assert.Equal(t, M1, i)

	sayChannel1.Step()
	i = <-ch
	assert.Equal(t, M2, i)

}

func (ch intChannel) reportAndIncrement(v interface{}) interface{} {
	ch <- v.(int)
	return v.(int) + 1
}
func (ch intChannel) reportAndDecrement(v interface{}) interface{} {
	ch <- v.(int)
	return v.(int) - 1
}

func TestPausedPromiseDemo(t *testing.T) {
	ch := make(intChannel)

	sayChannel0 := promise.Paused(0).
		Then(ch.reportAndIncrement).
		Then(ch.reportAndIncrement).(promise.PausedPromise)

	sayChannel1 := promise.Paused(100).
		Then(ch.reportAndDecrement).
		Then(ch.reportAndDecrement).
		Then(ch.reportAndDecrement).(promise.PausedPromise)

	sayChannel1.Step()
	i := <-ch
	assert.Equal(t, 100, i)

	sayChannel0.Step()
	i = <-ch
	assert.Equal(t, 0, i)

	sayChannel1.Continue()
	i = <-ch
	assert.Equal(t, 99, i)
	i = <-ch
	assert.Equal(t, 98, i)
}
